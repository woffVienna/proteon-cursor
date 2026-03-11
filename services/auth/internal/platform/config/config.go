package config

import (
	platformconfig "github.com/woffVienna/proteon-cursor/libs/platform/config"
)

type Config = platformconfig.Config[ServiceConfig]

type ServiceConfig struct {
	IdentityURL string
}

func Load() (Config, error) {
	loader := platformconfig.NewLoader[ServiceConfig](platformconfig.LoaderOptions{
		WorkingDir:         ".",
		DefaultServiceName: "auth-service",
		DefaultEnvironment: "dev",
		DefaultMarket:      "AT",
		DefaultVersion:     "dev",
		DefaultPort:        "8083",
	})

	return loader.Load(func(env platformconfig.Env) (ServiceConfig, error) {
		return ServiceConfig{
			IdentityURL: env.String("IDENTITY_URL", "http://localhost:8081"),
		}, nil
	})
}
