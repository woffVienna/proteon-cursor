package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config[T any] struct {
	ServiceName string
	Environment string
	Market      string
	Version     string
	HTTP        HTTPConfig
	Service     T
}

type HTTPConfig struct {
	Port          string
	PublicBaseURL string
}

type LoaderOptions struct {
	WorkingDir string

	DefaultServiceName string
	DefaultEnvironment string
	DefaultMarket      string
	DefaultVersion     string

	DefaultPort          string
	DefaultPublicBaseURL string
}

type Loader[T any] struct {
	opts LoaderOptions
}

type Env struct{}

type ServiceParser[T any] func(Env) (T, error)

func NewLoader[T any](opts LoaderOptions) Loader[T] {
	return Loader[T]{opts: withDefaultOptions(opts)}
}

func (l Loader[T]) Load(parseService ServiceParser[T]) (Config[T], error) {
	if err := LoadLocalEnvFile(l.opts.WorkingDir); err != nil {
		return Config[T]{}, err
	}

	port := envString("PORT", l.opts.DefaultPort)
	if err := validatePort(port); err != nil {
		return Config[T]{}, err
	}

	service, err := parseService(Env{})
	if err != nil {
		return Config[T]{}, err
	}

	publicBaseURL := envString("PUBLIC_BASE_URL", l.opts.DefaultPublicBaseURL)
	if publicBaseURL == "" {
		publicBaseURL = "http://localhost:" + port
	}

	return Config[T]{
		ServiceName: envString("SERVICE_NAME", l.opts.DefaultServiceName),
		Environment: envString("ENV", l.opts.DefaultEnvironment),
		Market:      envString("MARKET", l.opts.DefaultMarket),
		Version:     envString("VERSION", l.opts.DefaultVersion),
		HTTP: HTTPConfig{
			Port:          port,
			PublicBaseURL: publicBaseURL,
		},
		Service: service,
	}, nil
}

func (Env) String(key, fallback string) string {
	return envString(key, fallback)
}

func withDefaultOptions(opts LoaderOptions) LoaderOptions {
	if opts.WorkingDir == "" {
		opts.WorkingDir = "."
	}
	if opts.DefaultEnvironment == "" {
		opts.DefaultEnvironment = "dev"
	}
	if opts.DefaultMarket == "" {
		opts.DefaultMarket = "AT"
	}
	if opts.DefaultVersion == "" {
		opts.DefaultVersion = "dev"
	}
	if opts.DefaultPort == "" {
		opts.DefaultPort = "8081"
	}
	return opts
}

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func validatePort(port string) error {
	v, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid PORT %q: %w", port, err)
	}
	if v <= 0 || v > 65535 {
		return fmt.Errorf("invalid PORT %q: out of range", port)
	}
	return nil
}
