package http

import (
	"fmt"
	"net/http"
	"net/http/httputil"
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

// Server is the HTTP adapter for the API gateway.
type Server struct {
	cfg            Config
	identityProxy  *httputil.ReverseProxy
	authMiddleware func(http.Handler) http.Handler
}

// NewServer creates a gateway HTTP server.
func NewServer(
	cfg Config,
	identityProxy *httputil.ReverseProxy,
	authMiddleware func(http.Handler) http.Handler,
) *Server {
	return &Server{
		cfg:            cfg,
		identityProxy:  identityProxy,
		authMiddleware: authMiddleware,
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
		Title:        "Proteon API Gateway - Swagger UI",
		SpecFilePath: s.cfg.OpenAPIBundlePath,
		HealthRoute:  "/v1/health",
	})

	r.Group(func(r chi.Router) {
		r.Post("/v1/auth/exchange", s.identityProxy.ServeHTTP)
		r.Get("/v1/.well-known/jwks.json", s.identityProxy.ServeHTTP)
	})

	r.Group(func(r chi.Router) {
		r.Use(s.authMiddleware)
		r.Get("/v1/users/{userId}", s.identityProxy.ServeHTTP)
	})

	return r
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":"+s.cfg.Port, s.Router())
}

// String implements fmt.Stringer.
func (s *Server) String() string {
	return fmt.Sprintf("api-gateway(port=%s)", s.cfg.Port)
}
