package main

import (
	"crypto/rand"
	"encoding/hex"
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

	identityStore := auth.NewMemoryIdentityStore(generateUUID)
	issuer, err := auth.NewJWTIssuer(cfg.Service.JWT.Issuer, cfg.Service.JWT.Audience)
	if err != nil {
		log.Fatalf("failed to create JWT issuer: %v", err)
	}

	authSvc := authapp.NewService(identityStore, identityStore, issuer)

	httpCfg := httpadapter.Config{
		Port:              cfg.HTTP.Port,
		OpenAPIBundlePath: ".build/generated/openapi.bundle.yml",
		ServiceName:       cfg.ServiceName,
		Version:           cfg.Version,
	}
	srv := httpadapter.NewServer(httpCfg, authSvc, issuer)

	addr := ":" + cfg.HTTP.Port
	log.Printf("Identity service listening on %s", addr)
	log.Printf("Swagger UI:  %s/swagger", cfg.HTTP.PublicBaseURL)
	log.Printf("OpenAPI spec: %s/openapi.yaml", cfg.HTTP.PublicBaseURL)
	log.Printf("starting %s", srv)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return hex.EncodeToString(b[:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:])
}
