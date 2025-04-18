package common

import (
	"article-service/pkg/gateway"
	"net/http"
	"strings"
)

// HandlerMux handles the API requests, applying CORS and Content-Type checks.
func HandlerMux(g gateway.GatewayInterface, allowedContentTypes []string) http.Handler {
	// Wrap the handler with CORS and Content-Type validation middleware
	return CORS(StrictContentTypeMiddleware(allowedContentTypes)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request path is for the API (starts with /api)
		if strings.HasPrefix(r.URL.Path, "/api") {
			// Serve the API route using the gateway's runtime mux
			g.GetRuntimeMux().ServeHTTP(w, r)
			return
		}

		// Serve other routes using the gateway's regular mux
		g.GetMux().ServeHTTP(w, r)
	})))
}

// StrictContentTypeMiddleware enforces strict checks on the Content-Type header.
func StrictContentTypeMiddleware(allowedContentTypes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip methods that don't usually have a body (e.g., GET, HEAD, OPTIONS)
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodDelete || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			// Get the Content-Type header
			contentType := r.Header.Get("Content-Type")
			if contentType == "" {
				http.Error(w, "Missing Content-Type", http.StatusBadRequest)
				return
			}

			// Check if the Content-Type is allowed
			valid := false
			for _, allowed := range allowedContentTypes {
				if strings.HasPrefix(contentType, allowed) {
					valid = true
					break
				}
			}

			if !valid {
				http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
				return
			}

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}
