package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/woffVienna/proteon-cursor/libs/platform/httpcommon"
	"github.com/woffVienna/proteon-cursor/libs/platform/security/jwtverifier"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/adapters/http/generated/server"
	authapp "github.com/woffVienna/proteon-cursor/services/identity/internal/application/auth"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/application/interfaces"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/domain"
)

// Handler implements server.StrictServerInterface.
type Handler struct {
	authSvc     *authapp.Service
	issuer      interfaces.TokenIssuer
	verifier    *jwtverifier.Verifier
	serviceName string
	version     string
}

// NewHandler creates an HTTP handler.
func NewHandler(
	authSvc *authapp.Service,
	issuer interfaces.TokenIssuer,
	verifier *jwtverifier.Verifier,
	serviceName string,
	version string,
) *Handler {
	return &Handler{
		authSvc:     authSvc,
		issuer:      issuer,
		verifier:    verifier,
		serviceName: serviceName,
		version:     version,
	}
}

var _ server.StrictServerInterface = (*Handler)(nil)

func (h *Handler) GetV1Health(ctx context.Context, _ server.GetV1HealthRequestObject) (server.GetV1HealthResponseObject, error) {
	svc := h.serviceName
	version := h.version
	return server.GetV1Health200JSONResponse(server.HealthResponse{
		Status:  server.Ok,
		Service: &svc,
		Version: &version,
	}), nil
}

func (h *Handler) GetV1WellKnownJwks(ctx context.Context, _ server.GetV1WellKnownJwksRequestObject) (server.GetV1WellKnownJwksResponseObject, error) {
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

func (h *Handler) PostV1AuthLogin(ctx context.Context, req server.PostV1AuthLoginRequestObject) (server.PostV1AuthLoginResponseObject, error) {
	if req.Body == nil {
		return server.PostV1AuthLogin400JSONResponse{
			BadRequestJSONResponse: server.BadRequestJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "BAD_REQUEST", Message: "missing request body"},
			}),
		}, nil
	}

	tenant := "proteon.dev"
	if req.Body.Tenant != nil && *req.Body.Tenant != "" {
		tenant = *req.Body.Tenant
	}

	pair, err := h.authSvc.Login(ctx, req.Body.Login, req.Body.Password, tenant)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			return server.PostV1AuthLogin401JSONResponse{
				UnauthorizedJSONResponse: server.UnauthorizedJSONResponse(server.ErrorResponse{
					Error: server.ErrorBody{Code: "INVALID_CREDENTIALS", Message: "invalid login or password"},
				}),
			}, nil
		}
		return server.PostV1AuthLogin500JSONResponse{
			InternalErrorJSONResponse: server.InternalErrorJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()},
			}),
		}, nil
	}

	return server.PostV1AuthLogin200JSONResponse(server.TokenPairResponse{
		AccessToken:  pair.AccessToken,
		TokenType:    server.Bearer,
		ExpiresIn:    pair.ExpiresIn,
		RefreshToken: pair.RefreshToken,
	}), nil
}

func (h *Handler) PostV1AuthRefresh(ctx context.Context, req server.PostV1AuthRefreshRequestObject) (server.PostV1AuthRefreshResponseObject, error) {
	if req.Body == nil {
		return server.PostV1AuthRefresh400JSONResponse{
			BadRequestJSONResponse: server.BadRequestJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "BAD_REQUEST", Message: "missing request body"},
			}),
		}, nil
	}
	if req.Body.RefreshToken == "" {
		return server.PostV1AuthRefresh400JSONResponse{
			BadRequestJSONResponse: server.BadRequestJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "BAD_REQUEST", Message: "refresh_token is required"},
			}),
		}, nil
	}

	pair, err := h.authSvc.Refresh(ctx, req.Body.RefreshToken)
	if err != nil {
		if err == domain.ErrInvalidRefreshToken || err == domain.ErrRefreshTokenExpired {
			return server.PostV1AuthRefresh401JSONResponse{
				UnauthorizedJSONResponse: server.UnauthorizedJSONResponse(server.ErrorResponse{
					Error: server.ErrorBody{
						Code:    "INVALID_REFRESH_TOKEN",
						Message: "invalid or expired refresh token",
					},
				}),
			}, nil
		}
		return server.PostV1AuthRefresh500JSONResponse{
			InternalErrorJSONResponse: server.InternalErrorJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()},
			}),
		}, nil
	}

	return server.PostV1AuthRefresh200JSONResponse(server.TokenPairResponse{
		AccessToken:  pair.AccessToken,
		TokenType:    server.Bearer,
		ExpiresIn:    pair.ExpiresIn,
		RefreshToken: pair.RefreshToken,
	}), nil
}

func (h *Handler) PostV1AuthLogout(ctx context.Context, req server.PostV1AuthLogoutRequestObject) (server.PostV1AuthLogoutResponseObject, error) {
	if req.Body == nil {
		return server.PostV1AuthLogout400JSONResponse{
			BadRequestJSONResponse: server.BadRequestJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "BAD_REQUEST", Message: "missing request body"},
			}),
		}, nil
	}
	if req.Body.RefreshToken == "" {
		return server.PostV1AuthLogout400JSONResponse{
			BadRequestJSONResponse: server.BadRequestJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "BAD_REQUEST", Message: "refresh_token is required"},
			}),
		}, nil
	}

	_ = h.authSvc.Logout(ctx, req.Body.RefreshToken)
	return server.PostV1AuthLogout204Response{}, nil
}

func (h *Handler) GetV1Me(ctx context.Context, _ server.GetV1MeRequestObject) (server.GetV1MeResponseObject, error) {
	req := httpcommon.HTTPRequestFromContext(ctx)
	if req == nil {
		return server.GetV1Me500JSONResponse{
			InternalErrorJSONResponse: server.InternalErrorJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "INTERNAL_ERROR", Message: "request not available in context"},
			}),
		}, nil
	}

	rawToken, err := httpcommon.ExtractBearer(req.Header.Get("Authorization"))
	if err != nil {
		return server.GetV1Me401JSONResponse{
			UnauthorizedJSONResponse: server.UnauthorizedJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "UNAUTHORIZED", Message: err.Error()},
			}),
		}, nil
	}

	claims, err := h.verifier.Verify(rawToken)
	if err != nil {
		return server.GetV1Me401JSONResponse{
			UnauthorizedJSONResponse: server.UnauthorizedJSONResponse(server.ErrorResponse{
				Error: server.ErrorBody{Code: "UNAUTHORIZED", Message: "invalid token"},
			}),
		}, nil
	}

	return server.GetV1Me200JSONResponse(server.MeResponse{
		Sub:       claims.Subject,
		Tenant:    claims.Tenant,
		Scopes:    optionalScopesPtr(claims.Scopes),
		SessionId: optionalStringPtr(claims.SessionID),
	}), nil
}

func writeJSONError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(server.ErrorResponse{
		Error: server.ErrorBody{Code: code, Message: msg},
	})
}

func optionalStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func optionalScopesPtr(scopes []string) *[]string {
	if len(scopes) == 0 {
		return nil
	}
	return &scopes
}
