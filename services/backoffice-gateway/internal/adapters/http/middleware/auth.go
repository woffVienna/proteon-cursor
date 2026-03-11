package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/woffVienna/proteon-cursor/libs/platform/httpcommon"
	"github.com/woffVienna/proteon-cursor/libs/platform/security/jwtverifier"
)

const (
	HeaderPlatformUserID      = "X-Platform-User-Id"
	HeaderPlatformTenant      = "X-Platform-Tenant"
	HeaderPlatformSubjectType = "X-Platform-Subject-Type"
)

// Auth returns a chi middleware that validates JWTs and injects verified
// identity context into downstream request headers.
func Auth(verifier *jwtverifier.Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawToken, err := httpcommon.ExtractBearer(r.Header.Get("Authorization"))
			if err != nil {
				writeAuthError(w, "UNAUTHORIZED", "missing or invalid authorization header")
				return
			}

			claims, err := verifier.Verify(rawToken)
			if err != nil {
				writeAuthError(w, "UNAUTHORIZED", "invalid token")
				return
			}

			r.Header.Set(HeaderPlatformUserID, claims.Subject)
			r.Header.Set(HeaderPlatformTenant, claims.Tenant)
			if claims.SubjectType != "" {
				r.Header.Set(HeaderPlatformSubjectType, claims.SubjectType)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func writeAuthError(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
