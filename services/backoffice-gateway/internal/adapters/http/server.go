package http

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/woffVienna/proteon-cursor/libs/platform/httpcommon"
	bomw "github.com/woffVienna/proteon-cursor/services/backoffice-gateway/internal/adapters/http/middleware"
)

// Config holds HTTP server configuration.
type Config struct {
	Port              string
	OpenAPIBundlePath string
	ServiceName       string
	Version           string
	// BasePath is the URL path prefix when behind ingress (e.g. /backoffice). Empty for direct access.
	BasePath string
}

// Server is the HTTP adapter for the backoffice gateway.
type Server struct {
	cfg       Config
	authProxy *httputil.ReverseProxy
	appKeyMW  func(http.Handler) http.Handler
	jwtAuthMW func(http.Handler) http.Handler
}

// NewServer creates a gateway HTTP server.
func NewServer(
	cfg Config,
	authProxy *httputil.ReverseProxy,
	appKeyMW func(http.Handler) http.Handler,
	jwtAuthMW func(http.Handler) http.Handler,
) *Server {
	return &Server{
		cfg:       cfg,
		authProxy: authProxy,
		appKeyMW:  appKeyMW,
		jwtAuthMW: jwtAuthMW,
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

	mount := func(r chi.Router) {
		specRoute := "/openapi.yaml"
		swaggerRoute := "/swagger"
		healthRoute := "/v1/health"
		if s.cfg.BasePath != "" {
			specRoute = s.cfg.BasePath + "/openapi.yaml"
			swaggerRoute = s.cfg.BasePath + "/swagger"
			healthRoute = s.cfg.BasePath + "/v1/health"
		}
		httpcommon.MountDocsAndHealth(r, httpcommon.DocsOptions{
			Title:        "Proteon Backoffice Gateway - Swagger UI",
			SpecFilePath: s.cfg.OpenAPIBundlePath,
			SpecRoute:    specRoute,
			SwaggerRoute: swaggerRoute,
			HealthRoute:  healthRoute,
		})

		r.Group(func(r chi.Router) {
			r.Use(s.appKeyMW)
			r.Post("/v1/auth/login", authLoginProxy(s.authProxy))
		})

		r.Group(func(r chi.Router) {
			r.Use(s.jwtAuthMW)
			_ = bomw.HeaderPlatformUserID
		})
	}

	if s.cfg.BasePath != "" {
		r.Route(s.cfg.BasePath, mount)
	} else {
		mount(r)
	}

	return r
}

// authLoginProxy wraps the auth proxy and rewrites the path from /v1/auth/login to /v1/login
// so the auth service (which exposes POST /v1/login) receives the correct path.
func authLoginProxy(proxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r2 := r.Clone(r.Context())
		r2.URL = &url.URL{Path: "/v1/login", RawPath: ""}
		if r.URL.RawQuery != "" {
			r2.URL.RawQuery = r.URL.RawQuery
		}
		proxy.ServeHTTP(w, r2)
	}
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(":"+s.cfg.Port, s.Router())
}

// String implements fmt.Stringer.
func (s *Server) String() string {
	return fmt.Sprintf("backoffice-gateway(port=%s)", s.cfg.Port)
}
