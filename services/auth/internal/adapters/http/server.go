package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/woffVienna/proteon-cursor/libs/platform/httpcommon"
)

// Config holds HTTP server configuration.
type Config struct {
	Port              string
	OpenAPIBundlePath string
	ServiceName       string
	Version           string
}

// Server is the HTTP adapter for the auth service.
type Server struct {
	cfg     Config
	handler *Handler
}

// NewServer creates an HTTP server with the given dependencies.
func NewServer(cfg Config, handler *Handler) *Server {
	return &Server{
		cfg:     cfg,
		handler: handler,
	}
}

// Router returns the HTTP handler with all routes configured.
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(chimw.Logger)

	httpcommon.MountDocsAndHealth(r, httpcommon.DocsOptions{
		Title:        "Proteon Auth Service - Swagger UI",
		SpecFilePath: s.cfg.OpenAPIBundlePath,
		HealthRoute:  "/v1/health",
	})

	r.Post("/v1/login", s.handler.Login)

	return r
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":"+s.cfg.Port, s.Router())
}

// String implements fmt.Stringer.
func (s *Server) String() string {
	return fmt.Sprintf("auth-service(port=%s)", s.cfg.Port)
}
