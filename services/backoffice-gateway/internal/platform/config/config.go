package config

import (
	platformconfig "github.com/woffVienna/proteon-cursor/libs/platform/config"
)

type Config = platformconfig.Config[ServiceConfig]

type ServiceConfig struct {
	JWT      JWTConfig
	Upstream UpstreamConfig
	AppKey   string
	// BasePath is the URL path prefix when behind an ingress (e.g. /backoffice). Empty for direct access.
	BasePath string
}

type JWTConfig struct {
	Issuer   string
	Audience string
}

type UpstreamConfig struct {
	IdentityURL string
	AuthURL     string
}

func Load() (Config, error) {
	loader := platformconfig.NewLoader[ServiceConfig](platformconfig.LoaderOptions{
		WorkingDir:         ".",
		DefaultServiceName: "backoffice-gateway",
		DefaultEnvironment: "dev",
		DefaultMarket:      "AT",
		DefaultVersion:     "dev",
		DefaultPort:        "8082",
	})

	return loader.Load(func(env platformconfig.Env) (ServiceConfig, error) {
		return ServiceConfig{
			JWT: JWTConfig{
				Issuer:   env.String("JWT_ISSUER", "proteon.identity"),
				Audience: env.String("JWT_AUDIENCE", "backoffice"),
			},
			Upstream: UpstreamConfig{
				IdentityURL: env.String("IDENTITY_URL", "http://localhost:8081"),
				AuthURL:     env.String("AUTH_URL", "http://localhost:8083"),
			},
			AppKey:    env.String("APP_KEY", "dev-backoffice-key-001"),
			BasePath:  env.String("BASE_PATH", ""),
		}, nil
	})
}
