package main

import (
	"log"

	"github.com/woffVienna/proteon-cursor/services/identity/internal/adapters/auth"
	httpadapter "github.com/woffVienna/proteon-cursor/services/identity/internal/adapters/http"
	authapp "github.com/woffVienna/proteon-cursor/services/identity/internal/application/auth"
)

func main() {
	cfg := httpadapter.DefaultConfig()

	// Domain adapters
	validator := auth.NewDemoValidator()
	refreshStore := auth.NewMemoryStore()
	issuer, err := auth.NewJWTIssuer()
	if err != nil {
		log.Fatalf("failed to create JWT issuer: %v", err)
	}

	// Application layer
	authSvc := authapp.NewService(validator, refreshStore, issuer)

	// HTTP adapter
	srv := httpadapter.NewServer(cfg, authSvc, issuer)

	addr := ":" + cfg.Port
	log.Printf("Identity service listening on %s", addr)
	log.Printf("Swagger UI:  http://localhost%s/swagger", addr)
	log.Printf("OpenAPI spec: http://localhost%s/openapi.yaml", addr)
	log.Printf("starting %s", srv)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
