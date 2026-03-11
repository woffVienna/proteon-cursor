package main

import (
	"log"
	"time"

	"github.com/woffVienna/proteon-cursor/libs/platform/security/jwtverifier"
	"github.com/woffVienna/proteon-cursor/services/api-gateway/internal/adapters/auth"
	httpadapter "github.com/woffVienna/proteon-cursor/services/api-gateway/internal/adapters/http"
	"github.com/woffVienna/proteon-cursor/services/api-gateway/internal/adapters/http/middleware"
	"github.com/woffVienna/proteon-cursor/services/api-gateway/internal/adapters/http/proxy"
	"github.com/woffVienna/proteon-cursor/services/api-gateway/internal/platform/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("fetching JWKS from %s", cfg.Service.Upstream.IdentityURL)
	keys, err := auth.FetchJWKS(cfg.Service.Upstream.IdentityURL)
	if err != nil {
		log.Fatalf("failed to fetch JWKS: %v", err)
	}
	log.Printf("loaded %d signing key(s) from identity", len(keys))

	verifier := jwtverifier.New(jwtverifier.Config{
		Issuer:   cfg.Service.JWT.Issuer,
		Audience: cfg.Service.JWT.Audience,
		Keys:     keys,
		Leeway:   30 * time.Second,
	})

	authMW := middleware.Auth(verifier)

	identityProxy, err := proxy.New(cfg.Service.Upstream.IdentityURL)
	if err != nil {
		log.Fatalf("failed to create identity proxy: %v", err)
	}

	httpCfg := httpadapter.Config{
		Port:              cfg.HTTP.Port,
		OpenAPIBundlePath: ".build/generated/openapi.bundle.yml",
		ServiceName:       cfg.ServiceName,
		Version:           cfg.Version,
	}
	srv := httpadapter.NewServer(httpCfg, identityProxy, authMW)

	addr := ":" + cfg.HTTP.Port
	log.Printf("API gateway listening on %s", addr)
	log.Printf("Swagger UI:  %s/swagger", cfg.HTTP.PublicBaseURL)
	log.Printf("OpenAPI spec: %s/openapi.yaml", cfg.HTTP.PublicBaseURL)
	log.Printf("starting %s", srv)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
