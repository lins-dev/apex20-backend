package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
)

// stubPermissionRepository is an in-memory stub for handler tests.
type stubPermissionRepository struct {
	permissions []permission.Permission
}

func (s *stubPermissionRepository) ExistsAny(_ context.Context) (bool, error) {
	return len(s.permissions) > 0, nil
}

func (s *stubPermissionRepository) ListPermissions(_ context.Context) ([]permission.Permission, error) {
	return s.permissions, nil
}

func (s *stubPermissionRepository) GetPermissionByID(_ context.Context, id uuid.UUID) (permission.Permission, error) {
	for _, p := range s.permissions {
		if p.ID == id {
			return p, nil
		}
	}
	return permission.Permission{}, port.ErrNotFound
}

func (s *stubPermissionRepository) CreatePermission(_ context.Context, p permission.Permission) error {
	s.permissions = append(s.permissions, p)
	return nil
}

func (s *stubPermissionRepository) UpdatePermission(_ context.Context, p permission.Permission) error {
	for i, existing := range s.permissions {
		if existing.ID == p.ID {
			s.permissions[i] = p
			return nil
		}
	}
	return port.ErrNotFound
}

func (s *stubPermissionRepository) DeletePermission(_ context.Context, id uuid.UUID, _ time.Time) (bool, error) {
	for i, p := range s.permissions {
		if p.ID == id {
			s.permissions = append(s.permissions[:i], s.permissions[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

func newServerWithPermissions(perms []permission.Permission) *adapter.ChiServer {
	repo := &stubPermissionRepository{permissions: perms}
	server := adapter.NewChiServer()
	adapter.RegisterPermissionHandler(server.GetAPI(), repo)
	return server
}

func TestPermissionHandler_List(t *testing.T) {
	id := uuid.New()
	server := newServerWithPermissions([]permission.Permission{
		{ID: id, Name: "chat.send", Description: "Enviar mensagens", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/permissions", nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "chat.send")
}

func TestPermissionHandler_Get(t *testing.T) {
	id := uuid.New()
	server := newServerWithPermissions([]permission.Permission{
		{ID: id, Name: "token.move.own", Description: "Mover token", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/permissions/"+id.String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "token.move.own")
}

func TestPermissionHandler_Get_NotFound(t *testing.T) {
	server := newServerWithPermissions(nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/permissions/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPermissionHandler_Create(t *testing.T) {
	server := newServerWithPermissions(nil)

	body := `{"name":"new.permission","description":"Nova permissão"}`
	req := httptest.NewRequest(http.MethodPost, "/admin/permissions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "new.permission")
}

func TestPermissionHandler_Update(t *testing.T) {
	id := uuid.New()
	server := newServerWithPermissions([]permission.Permission{
		{ID: id, Name: "old.name", Description: "Old", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	})

	body := `{"name":"new.name","description":"New"}`
	req := httptest.NewRequest(http.MethodPut, "/admin/permissions/"+id.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "new.name")
}

func TestPermissionHandler_Delete(t *testing.T) {
	id := uuid.New()
	server := newServerWithPermissions([]permission.Permission{
		{ID: id, Name: "to.delete", Description: "Delete me", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	})

	req := httptest.NewRequest(http.MethodDelete, "/admin/permissions/"+id.String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestPermissionHandler_Delete_NotFound(t *testing.T) {
	server := newServerWithPermissions(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/permissions/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPermissionHandler_Create_InvalidBody(t *testing.T) {
	server := newServerWithPermissions(nil)

	req := httptest.NewRequest(http.MethodPost, "/admin/permissions", strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func permissionsFromBody(t *testing.T, body string) []map[string]any {
	t.Helper()
	var resp struct {
		Permissions []map[string]any `json:"permissions"`
	}
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	return resp.Permissions
}
