package http

import (
	"encoding/json"
	"net/http"

	"github.com/woffVienna/proteon-cursor/services/auth/internal/application/login"
	"github.com/woffVienna/proteon-cursor/services/auth/internal/domain"
)

// Handler handles auth HTTP requests.
type Handler struct {
	loginSvc *login.Service
}

// NewHandler creates a new Handler.
func NewHandler(loginSvc *login.Service) *Handler {
	return &Handler{loginSvc: loginSvc}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int32  `json:"expires_in"`
}

// Login handles POST /v1/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "username and password are required")
		return
	}

	result, err := h.loginSvc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid username or password")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	resp := loginResponse{
		AccessToken: result.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   result.ExpiresIn,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		// best-effort; nothing else to do
		_ = err
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]map[string]string{
		"error": {
			"code":    code,
			"message": message,
		},
	})
}
