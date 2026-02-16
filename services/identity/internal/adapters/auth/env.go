package auth

import "os"

// IssuerFromEnv returns JWT issuer from JWT_ISSUER env, or default.
func IssuerFromEnv() string {
	if v := os.Getenv("JWT_ISSUER"); v != "" {
		return v
	}
	return "proteon.identity"
}

// AudienceFromEnv returns JWT audience from JWT_AUDIENCE env, or default.
func AudienceFromEnv() string {
	if v := os.Getenv("JWT_AUDIENCE"); v != "" {
		return v
	}
	return "proteon-api"
}
