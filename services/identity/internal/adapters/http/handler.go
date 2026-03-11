package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/adapters/http/generated/server"
	authapp "github.com/woffVienna/proteon-cursor/services/identity/internal/application/auth"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/application/interfaces"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/domain"
)

// Handler implements server.StrictServerInterface.
type Handler struct {
	authSvc     *authapp.Service
	issuer      interfaces.TokenIssuer
	serviceName string
	version     string
}

// NewHandler creates an HTTP handler.
func NewHandler(
	authSvc *authapp.Service,
	issuer interfaces.TokenIssuer,
	serviceName string,
	version string,
) *Handler {
	return &Handler{
		authSvc:     authSvc,
		issuer:      issuer,
		serviceName: serviceName,
		version:     version,
	}
}

var _ server.StrictServerInterface = (*Handler)(nil)

func (h *Handler) GetV1Health(_ context.Context, _ server.GetV1HealthRequestObject) (server.GetV1HealthResponseObject, error) {
	svc := h.serviceName
	version := h.version
	return server.GetV1Health200JSONResponse(server.HealthResponse{
		Status:  server.Ok,
		Service: &svc,
		Version: &version,
	}), nil
}

func (h *Handler) GetV1WellKnownJwks(_ context.Context, _ server.GetV1WellKnownJwksRequestObject) (server.GetV1WellKnownJwksResponseObject, error) {
	alg := "EdDSA"
	use := "sig"
	jwk := server.Jwk{
		Kty: "OKP",
		Kid: h.issuer.Kid(),
		Alg: &alg,
		Use: &use,
	}
	jwk.Set("crv", "Ed25519")
	jwk.Set("x", base64.RawURLEncoding.EncodeToString(h.issuer.PublicKey()))
	return server.GetV1WellKnownJwks200JSONResponse{Keys: []server.Jwk{jwk}}, nil
}

func (h *Handler) PostV1AuthExchange(ctx context.Context, req server.PostV1AuthExchangeRequestObject) (server.PostV1AuthExchangeResponseObject, error) {
	if req.Body == nil {
		return server.PostV1AuthExchange400JSONResponse{
			BadRequestJSONResponse: server.BadRequestJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "BAD_REQUEST", Message: "missing request body"},
			}),
		}, nil
	}

	tenant := ""
	if req.Body.Tenant != nil {
		tenant = *req.Body.Tenant
	}

	result, err := h.authSvc.Exchange(ctx, req.Body.Provider, req.Body.ExternalUserId, tenant)
	if err != nil {
		if err == domain.ErrInvalidAssertion {
			return server.PostV1AuthExchange400JSONResponse{
				BadRequestJSONResponse: server.BadRequestJSONResponse(server.ErrorResponse{
					Error: server.ErrorBody{Code: "INVALID_ASSERTION", Message: "invalid external identity assertion"},
				}),
			}, nil
		}
		return server.PostV1AuthExchange500JSONResponse{
			InternalErrorJSONResponse: server.InternalErrorJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "INTERNAL_ERROR", Message: "internal error"},
			}),
		}, nil
	}

	platformUserUUID, err := uuid.Parse(result.PlatformUserID)
	if err != nil {
		return server.PostV1AuthExchange500JSONResponse{
			InternalErrorJSONResponse: server.InternalErrorJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "INTERNAL_ERROR", Message: "internal error"},
			}),
		}, nil
	}

	return server.PostV1AuthExchange200JSONResponse(server.AuthExchangeResponse{
		AccessToken:    result.AccessToken,
		TokenType:      server.AuthExchangeResponseTokenTypeBearer,
		ExpiresIn:      result.ExpiresIn,
		PlatformUserId: platformUserUUID,
	}), nil
}

func (h *Handler) GetV1UsersUserId(ctx context.Context, req server.GetV1UsersUserIdRequestObject) (server.GetV1UsersUserIdResponseObject, error) {
	identity, err := h.authSvc.GetIdentity(ctx, req.UserId.String())
	if err != nil {
		if err == domain.ErrIdentityNotFound {
			return server.GetV1UsersUserId404JSONResponse{
				NotFoundJSONResponse: server.NotFoundJSONResponse(server.ErrorResponse{
					Error: server.ErrorBody{Code: "NOT_FOUND", Message: "platform identity not found"},
				}),
			}, nil
		}
		return server.GetV1UsersUserId500JSONResponse{
			InternalErrorJSONResponse: server.InternalErrorJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "INTERNAL_ERROR", Message: "internal error"},
			}),
		}, nil
	}

	platformUserUUID, err := uuid.Parse(identity.PlatformUserID)
	if err != nil {
		return server.GetV1UsersUserId500JSONResponse{
			InternalErrorJSONResponse: server.InternalErrorJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "INTERNAL_ERROR", Message: "internal error"},
			}),
		}, nil
	}

	var tenant *string
	if identity.Tenant != "" {
		tenant = &identity.Tenant
	}

	return server.GetV1UsersUserId200JSONResponse(server.PlatformIdentityResponse{
		PlatformUserId: platformUserUUID,
		Provider:       identity.Provider,
		ExternalUserId: identity.ExternalUserID,
		Tenant:         tenant,
		CreatedAt:      identity.CreatedAt,
	}), nil
}

func (h *Handler) PostInternalV1BackofficeTokens(ctx context.Context, req server.PostInternalV1BackofficeTokensRequestObject) (server.PostInternalV1BackofficeTokensResponseObject, error) {
	if req.Body == nil {
		return server.PostInternalV1BackofficeTokens400JSONResponse{
			BadRequestJSONResponse: server.BadRequestJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "BAD_REQUEST", Message: "missing request body"},
			}),
		}, nil
	}

	body := req.Body

	var tenant string
	if body.TenantId != nil {
		tenant = *body.TenantId
	}

	aud := "backoffice"
	if body.Audience != nil && *body.Audience != "" {
		aud = *body.Audience
	}

	result, err := h.authSvc.IssueBackofficeToken(ctx, body.UserId.String(), string(body.SubjectType), tenant, aud)
	if err != nil {
		return server.PostInternalV1BackofficeTokens500JSONResponse{
			InternalErrorJSONResponse: server.InternalErrorJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "INTERNAL_ERROR", Message: "internal error"},
			}),
		}, nil
	}

	return server.PostInternalV1BackofficeTokens200JSONResponse(server.BackofficeTokenResponse{
		AccessToken: result.AccessToken,
		TokenType:   server.BackofficeTokenResponseTokenTypeBearer,
		ExpiresIn:   result.ExpiresIn,
	}), nil
}

func writeJSONError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(server.ErrorResponse{
		Error: server.ErrorBody{Code: code, Message: msg},
	})
}
