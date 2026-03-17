package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	// Alias para evitar conflito com net/http
	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
)

func TestHTTPServer_HealthCheck(t *testing.T) {
	// Setup do Servidor
	server := adapter.NewChiServer()
	
	// Criando uma requisição de teste
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	// Ação
	server.GetHandler().ServeHTTP(rr, req)

	// Asserção
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "operational")
}
