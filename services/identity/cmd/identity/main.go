package main

import (
	"log"

	"github.com/woffVienna/proteon-cursor/services/identity/internal/adapters/auth"
	httpadapter "github.com/woffVienna/proteon-cursor/services/identity/internal/adapters/http"
	authapp "github.com/woffVienna/proteon-cursor/services/identity/internal/application/auth"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/platform/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Domain adapters
	validator := auth.NewDemoValidator()
	refreshStore := auth.NewMemoryStore()
	issuer, err := auth.NewJWTIssuer(cfg.Service.JWT.Issuer, cfg.Service.JWT.Audience)
	if err != nil {
		log.Fatalf("failed to create JWT issuer: %v", err)
	}

	// Application layer
	authSvc := authapp.NewService(validator, refreshStore, issuer)

	// HTTP adapter
	httpCfg := httpadapter.Config{
		Port:              cfg.HTTP.Port,
		OpenAPIBundlePath: ".build/generated/openapi.bundle.yml",
		ServiceName:       cfg.ServiceName,
		Version:           cfg.Version,
		JWTIssuer:         cfg.Service.JWT.Issuer,
		JWTAudience:       cfg.Service.JWT.Audience,
	}
	srv := httpadapter.NewServer(httpCfg, authSvc, issuer)

	addr := ":" + cfg.HTTP.Port
	log.Printf("Identity service listening on %s (%s mode)", addr, cfg.RuntimeMode)
	log.Printf("Swagger UI:  %s/swagger", cfg.HTTP.PublicBaseURL)
	log.Printf("OpenAPI spec: %s/openapi.yaml", cfg.HTTP.PublicBaseURL)
	log.Printf("starting %s", srv)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
