package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
)

func TestHTTPServer_HealthCheck(t *testing.T) {
	server := adapter.NewChiServer()

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "operational")
}
