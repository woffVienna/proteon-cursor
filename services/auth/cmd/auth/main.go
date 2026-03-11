package main

import (
	"log"

	"github.com/woffVienna/proteon-cursor/services/auth/internal/adapters/credentials"
	httpadapter "github.com/woffVienna/proteon-cursor/services/auth/internal/adapters/http"
	identityadapter "github.com/woffVienna/proteon-cursor/services/auth/internal/adapters/identity"
	"github.com/woffVienna/proteon-cursor/services/auth/internal/application/login"
	"github.com/woffVienna/proteon-cursor/services/auth/internal/platform/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	credStore, err := credentials.NewMemoryStore()
	if err != nil {
		log.Fatalf("failed to create credential store: %v", err)
	}

	identityClient := identityadapter.NewClient(cfg.Service.IdentityURL)

	loginSvc := login.NewService(credStore, identityClient)
	handler := httpadapter.NewHandler(loginSvc)

	httpCfg := httpadapter.Config{
		Port:              cfg.HTTP.Port,
		OpenAPIBundlePath: ".build/generated/openapi.bundle.yml",
		ServiceName:       cfg.ServiceName,
		Version:           cfg.Version,
	}
	srv := httpadapter.NewServer(httpCfg, handler)

	addr := ":" + cfg.HTTP.Port
	log.Printf("Auth service listening on %s", addr)
	log.Printf("Swagger UI:  %s/swagger", cfg.HTTP.PublicBaseURL)
	log.Printf("OpenAPI spec: %s/openapi.yaml", cfg.HTTP.PublicBaseURL)
	log.Printf("starting %s", srv)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
