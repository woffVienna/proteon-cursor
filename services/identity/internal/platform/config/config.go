package config

import (
	platformconfig "github.com/woffVienna/proteon-cursor/libs/platform/config"
)

type Config = platformconfig.Config[ServiceConfig]

type ServiceConfig struct {
	JWT JWTConfig
}

type JWTConfig struct {
	Issuer   string
	Audience string
}

func Load() (Config, error) {
	loader := platformconfig.NewLoader[ServiceConfig](platformconfig.LoaderOptions{
		WorkingDir:         ".",
		DefaultServiceName: "identity-service",
		DefaultEnvironment: "dev",
		DefaultMarket:      "AT",
		DefaultVersion:     "dev",
		DefaultPort:        "8081",
	})

	return loader.Load(func(env platformconfig.Env) (ServiceConfig, error) {
		return ServiceConfig{
			JWT: JWTConfig{
				Issuer:   env.String("JWT_ISSUER", "proteon.identity"),
				Audience: env.String("JWT_AUDIENCE", "proteon-api"),
			},
		}, nil
	})
}
