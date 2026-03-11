package main

import (
	"crypto/ed25519"
	"log"
	"time"

	"github.com/woffVienna/proteon-cursor/libs/platform/security/jwtverifier"
	boauth "github.com/woffVienna/proteon-cursor/services/backoffice-gateway/internal/adapters/auth"
	httpadapter "github.com/woffVienna/proteon-cursor/services/backoffice-gateway/internal/adapters/http"
	bomw "github.com/woffVienna/proteon-cursor/services/backoffice-gateway/internal/adapters/http/middleware"
	"github.com/woffVienna/proteon-cursor/services/backoffice-gateway/internal/adapters/http/proxy"
	"github.com/woffVienna/proteon-cursor/services/backoffice-gateway/internal/platform/config"
)

const (
	jwksRetryInterval = 2 * time.Second
	jwksRetryTimeout  = 60 * time.Second
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	identityURL := cfg.Service.Upstream.IdentityURL
	log.Printf("fetching JWKS from %s (retrying up to %v)", identityURL, jwksRetryTimeout)
	var keys map[string]ed25519.PublicKey
	deadline := time.Now().Add(jwksRetryTimeout)
	for {
		keys, err = boauth.FetchJWKS(identityURL)
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			log.Fatalf("failed to fetch JWKS after %v: %v", jwksRetryTimeout, err)
		}
		log.Printf("JWKS fetch failed (will retry): %v", err)
		time.Sleep(jwksRetryInterval)
	}
	log.Printf("loaded %d signing key(s) from identity", len(keys))

	verifier := jwtverifier.New(jwtverifier.Config{
		Issuer:   cfg.Service.JWT.Issuer,
		Audience: cfg.Service.JWT.Audience,
		Keys:     keys,
		Leeway:   30 * time.Second,
	})

	appKeyMW := bomw.AppKeyMiddleware(cfg.Service.AppKey)
	authMW := bomw.Auth(verifier)

	authProxy, err := proxy.New(cfg.Service.Upstream.AuthURL)
	if err != nil {
		log.Fatalf("failed to create auth proxy: %v", err)
	}

	httpCfg := httpadapter.Config{
		Port:              cfg.HTTP.Port,
		OpenAPIBundlePath: ".build/generated/openapi.bundle.yml",
		ServiceName:       cfg.ServiceName,
		Version:           cfg.Version,
		BasePath:          cfg.Service.BasePath,
	}
	srv := httpadapter.NewServer(httpCfg, authProxy, appKeyMW, authMW)

	addr := ":" + cfg.HTTP.Port
	log.Printf("Backoffice gateway listening on %s", addr)
	log.Printf("Swagger UI:  %s/swagger", cfg.HTTP.PublicBaseURL)
	log.Printf("OpenAPI spec: %s/openapi.yaml", cfg.HTTP.PublicBaseURL)
	log.Printf("starting %s", srv)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
