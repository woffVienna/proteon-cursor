package middleware

import (
	"encoding/json"
	"net/http"
)

// AppKeyMiddleware validates the X-App-Key header against a configured key.
func AppKeyMiddleware(expectedKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-App-Key")
			if key == "" || key != expectedKey {
				writeAppKeyError(w, "UNAUTHORIZED", "missing or invalid app key")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func writeAppKeyError(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]map[string]string{
		"error": {
			"code":    code,
			"message": message,
		},
	})
}
