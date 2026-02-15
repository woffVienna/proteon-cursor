package httpcommon

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
)

// DocsOptions controls the mounted docs/health endpoints.
type DocsOptions struct {
	// URL path where the raw spec is served (default: /openapi.yaml)
	SpecRoute string
	// Filesystem path to the spec file, relative to the working directory.
	// Example: "api/openapi/identity-service/openapi.yml" (or openapi.bundle.yaml if you prefer)
	SpecFilePath string

	// URL path where Swagger UI HTML is served (default: /swagger)
	SwaggerRoute string
	// Title shown in the Swagger UI page
	Title string

	// URL path where health is served (default: /health)
	HealthRoute string
}

func (o *DocsOptions) withDefaults() DocsOptions {
	out := *o
	if out.SpecRoute == "" {
		out.SpecRoute = "/openapi.yaml"
	}
	if out.SwaggerRoute == "" {
		out.SwaggerRoute = "/swagger"
	}
	if out.HealthRoute == "" {
		out.HealthRoute = "/health"
	}
	if out.Title == "" {
		out.Title = "Proteon API - Swagger UI"
	}
	return out
}

// MountDocsAndHealth mounts:
//   - GET {SpecRoute}   -> serves SpecFilePath
//   - GET {SwaggerRoute}-> minimal Swagger UI page pointing at {SpecRoute}
//   - GET {HealthRoute} -> plain "ok"
//
// Safe to call for internal services too (you can keep it private via network controls).
func MountDocsAndHealth(r chi.Router, opts DocsOptions) {
	o := opts.withDefaults()

	// 1) Serve the raw OpenAPI spec
	r.Get(o.SpecRoute, func(w http.ResponseWriter, req *http.Request) {
		if o.SpecFilePath == "" {
			http.Error(w, "spec file path not configured", http.StatusInternalServerError)
			return
		}

		// If you serve .yml or .yaml, keep the route stable (/openapi.yaml) to simplify swagger UI.
		// Optionally set a better content-type than ServeFile guesses.
		ext := strings.ToLower(path.Ext(o.SpecFilePath))
		switch ext {
		case ".yml", ".yaml":
			w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		case ".json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}

		http.ServeFile(w, req, o.SpecFilePath)
	})

	// 2) Health
	r.Get(o.HealthRoute, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// 3) Simple Swagger UI page
	r.Get(o.SwaggerRoute, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <title>%s</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: "%s",
        dom_id: "#swagger-ui",
      });
    };
  </script>
</body>
</html>`, htmlEscape(o.Title), o.SpecRoute)
	})
}

// tiny HTML escaper (good enough for title)
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
