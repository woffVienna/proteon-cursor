package httpcommon

import (
	"context"
	"net/http"
)

type ctxKey int

const ctxKeyHTTPRequest ctxKey = 1

// WithHTTPRequest stores *http.Request in the context.
// Useful for strict oapi-codegen handlers that only receive context.Context.
func WithHTTPRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxKeyHTTPRequest, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HTTPRequestFromContext extracts *http.Request from context if present.
func HTTPRequestFromContext(ctx context.Context) *http.Request {
	if v := ctx.Value(ctxKeyHTTPRequest); v != nil {
		if req, ok := v.(*http.Request); ok {
			return req
		}
	}
	return nil
}
