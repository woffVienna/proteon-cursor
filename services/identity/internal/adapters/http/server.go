package http

import (
	"crypto/ed25519"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/woffVienna/proteon-cursor/libs/platform/httpcommon"
	"github.com/woffVienna/proteon-cursor/libs/platform/security/jwtverifier"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/adapters/auth"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/adapters/http/generated/server"
	authapp "github.com/woffVienna/proteon-cursor/services/identity/internal/application/auth"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/application/interfaces"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Config holds HTTP server configuration.
type Config struct {
	Port              string
	OpenAPIBundlePath string
}

// DefaultConfig returns config from env with defaults.
func DefaultConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	return Config{
		Port:              port,
		OpenAPIBundlePath: ".build/generated/openapi.bundle.yml",
	}
}

// Server is the HTTP adapter.
type Server struct {
	cfg     Config
	handler *Handler
}

// NewServer creates an HTTP server with the given dependencies.
func NewServer(cfg Config, authSvc *authapp.Service, issuer interfaces.TokenIssuer) *Server {
	verifier := jwtverifier.New(jwtverifier.Config{
		Issuer:   auth.IssuerFromEnv(),
		Audience: auth.AudienceFromEnv(),
		Keys: map[string]ed25519.PublicKey{
			issuer.Kid(): issuer.PublicKey(),
		},
		Leeway: 30 * time.Second,
	})
	return &Server{
		cfg:     cfg,
		handler: NewHandler(authSvc, issuer, verifier),
	}
}

// Router returns the HTTP handler.
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Logger)
	r.Use(httpcommon.WithHTTPRequest)

	httpcommon.MountDocsAndHealth(r, httpcommon.DocsOptions{
		Title:        "Proteon Identity Service - Swagger UI",
		SpecFilePath: s.cfg.OpenAPIBundlePath,
	})

	strictOpts := server.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			writeJSONError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		},
	}

	handler := server.NewStrictHandlerWithOptions(s.handler, nil, strictOpts)
	r.Mount("/", server.HandlerFromMux(handler, chi.NewRouter()))

	return r
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":"+s.cfg.Port, s.Router())
}

// String implements fmt.Stringer.
func (s *Server) String() string {
	return fmt.Sprintf("identity-service(port=%s)", s.cfg.Port)
}
