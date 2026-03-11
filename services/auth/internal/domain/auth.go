package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// LoginResult is the result of a successful login.
type LoginResult struct {
	AccessToken string
	ExpiresIn   int32
}
