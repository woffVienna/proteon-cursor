package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// New creates a reverse proxy that forwards requests to the given target URL.
// The original request path and query string are preserved.
func New(targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":{"code":"BAD_GATEWAY","message":"upstream unavailable"}}`))
	}

	return proxy, nil
}
