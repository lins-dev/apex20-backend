package port

import "net/http"

// HTTPServer defines the interface for the inbound HTTP server adapter.
type HTTPServer interface {
	RegisterRoute(method, path string, handler http.HandlerFunc)
	GetHandler() http.Handler
}
