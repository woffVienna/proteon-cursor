package config

import (
	platformconfig "github.com/woffVienna/proteon-cursor/libs/platform/config"
)

type Config = platformconfig.Config[ServiceConfig]

type ServiceConfig struct {
	JWT      JWTConfig
	Upstream UpstreamConfig
}

type JWTConfig struct {
	Issuer   string
	Audience string
}

type UpstreamConfig struct {
	IdentityURL string
}

func Load() (Config, error) {
	loader := platformconfig.NewLoader[ServiceConfig](platformconfig.LoaderOptions{
		WorkingDir:         ".",
		DefaultServiceName: "api-gateway",
		DefaultEnvironment: "dev",
		DefaultMarket:      "AT",
		DefaultVersion:     "dev",
		DefaultPort:        "8080",
	})

	return loader.Load(func(env platformconfig.Env) (ServiceConfig, error) {
		return ServiceConfig{
			JWT: JWTConfig{
				Issuer:   env.String("JWT_ISSUER", "proteon.identity"),
				Audience: env.String("JWT_AUDIENCE", "proteon-api"),
			},
			Upstream: UpstreamConfig{
				IdentityURL: env.String("IDENTITY_URL", "http://localhost:8081"),
			},
		}, nil
	})
}
