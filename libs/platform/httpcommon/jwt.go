package httpcommon

import (
	"errors"
	"strings"
)

var ErrMissingAuthorization = errors.New("missing Authorization header")
var ErrInvalidAuthorization = errors.New("invalid Authorization header")

// ExtractBearer extracts the raw JWT from an Authorization header.
func ExtractBearer(h string) (string, error) {
	if h == "" {
		return "", ErrMissingAuthorization
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(h, prefix) {
		return "", ErrInvalidAuthorization
	}
	token := strings.TrimSpace(h[len(prefix):])
	if token == "" {
		return "", ErrInvalidAuthorization
	}
	return token, nil
}
