package port

import "net/http"

// HTTPServer define a interface para o servidor de entrega web.
type HTTPServer interface {
	RegisterRoute(method, path string, handler http.HandlerFunc)
	Start(port string) error
	GetHandler() http.Handler
}
