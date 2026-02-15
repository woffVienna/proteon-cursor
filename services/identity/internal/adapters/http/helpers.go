package http

import "os"

func issuerFromEnv() string {
	if v := os.Getenv("JWT_ISSUER"); v != "" {
		return v
	}
	return "proteon.identity"
}

func audienceFromEnv() string {
	if v := os.Getenv("JWT_AUDIENCE"); v != "" {
		return v
	}
	return "proteon-api"
}
