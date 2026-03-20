package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
)

func newServerWithRoles() *adapter.ChiServer {
	server := adapter.NewChiServer()
	adapter.RegisterRoleHandler(server.GetAPI())
	return server
}

func TestRoleHandler_List(t *testing.T) {
	server := newServerWithRoles()

	req := httptest.NewRequest(http.MethodGet, "/admin/roles", nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.Contains(t, body, "ROLE_GM")
	assert.Contains(t, body, "ROLE_PLAYER")
	assert.Contains(t, body, "ROLE_TRUSTED")
	assert.NotContains(t, body, "ROLE_ADMIN")
}
