package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ChiServer struct {
	router *chi.Mux
	api    huma.API
}

// HealthResponse define a estrutura de resposta para o Huma.
type HealthResponse struct {
	Body struct {
		Status string `json:"status" example:"operational"`
	}
}

func NewChiServer() *ChiServer {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Inicializa o Huma v2 sobre o Chi
	config := huma.DefaultConfig("Apex20 API", "0.0.1")
	api := humachi.New(r, config)

	server := &ChiServer{
		router: r,
		api:    api,
	}

	// Registra o Health Check via Huma
	huma.Register(api, huma.Operation{
		OperationID: "get-health",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health Check",
		Description: "Verifica se o servidor está operacional.",
		Tags:        []string{"System"},
	}, func(ctx context.Context, input *struct{}) (*HealthResponse, error) {
		resp := &HealthResponse{}
		resp.Body.Status = "Apex20 Backend is operational"
		return resp, nil
	})

	return server
}

func (s *ChiServer) RegisterRoute(method, path string, handler http.HandlerFunc) {
	s.router.MethodFunc(method, path, handler)
}

func (s *ChiServer) Start(port string) error {
	fmt.Printf("Backend server starting on :%s...\n", port)
	fmt.Printf("OpenAPI documentation available at :%s/docs\n", port)
	return http.ListenAndServe(":"+port, s.router)
}

func (s *ChiServer) GetHandler() http.Handler {
	return s.router
}

func (s *ChiServer) GetAPI() huma.API {
	return s.api
}
