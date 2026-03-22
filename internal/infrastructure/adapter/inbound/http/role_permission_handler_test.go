package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/permission"
	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
)

// stubRolePermissionRepository is an in-memory stub satisfying all role-permission use case interfaces.
type stubRolePermissionRepository struct {
	rolePermissions []permission.RolePermission
}

func (s *stubRolePermissionRepository) ListRolePermissions(_ context.Context) ([]permission.RolePermission, error) {
	return s.rolePermissions, nil
}

func (s *stubRolePermissionRepository) GetRolePermissionByID(_ context.Context, id uuid.UUID) (permission.RolePermission, error) {
	for _, rp := range s.rolePermissions {
		if rp.ID == id {
			return rp, nil
		}
	}
	return permission.RolePermission{}, port.ErrNotFound
}

func (s *stubRolePermissionRepository) CreateRolePermission(_ context.Context, rp permission.RolePermission) error {
	s.rolePermissions = append(s.rolePermissions, rp)
	return nil
}

func (s *stubRolePermissionRepository) DeleteRolePermission(_ context.Context, id uuid.UUID, _ time.Time) (bool, error) {
	for i, rp := range s.rolePermissions {
		if rp.ID == id {
			s.rolePermissions = append(s.rolePermissions[:i], s.rolePermissions[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

func newServerWithRolePermissions(rps []permission.RolePermission) *adapter.ChiServer {
	repo := &stubRolePermissionRepository{rolePermissions: rps}
	uc := adapter.RolePermissionUseCases{
		List:   usecase.NewListRolePermissionsUseCase(repo),
		Get:    usecase.NewGetRolePermissionUseCase(repo),
		Create: usecase.NewCreateRolePermissionUseCase(repo),
		Delete: usecase.NewDeleteRolePermissionUseCase(repo),
	}
	server := adapter.NewChiServer()
	adapter.RegisterRolePermissionHandler(server.GetAPI(), uc)
	return server
}

func TestRolePermissionHandler_List(t *testing.T) {
	permID := uuid.New()
	rpID := uuid.New()
	server := newServerWithRolePermissions([]permission.RolePermission{
		{ID: rpID, Role: permission.RolePlayer, PermissionID: permID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/role-permissions", nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), rpID.String())
}

func TestRolePermissionHandler_Get(t *testing.T) {
	permID := uuid.New()
	rpID := uuid.New()
	server := newServerWithRolePermissions([]permission.RolePermission{
		{ID: rpID, Role: permission.RoleGM, PermissionID: permID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/role-permissions/"+rpID.String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), rpID.String())
}

func TestRolePermissionHandler_Get_NotFound(t *testing.T) {
	server := newServerWithRolePermissions(nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/role-permissions/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRolePermissionHandler_Create(t *testing.T) {
	permID := uuid.New()
	server := newServerWithRolePermissions(nil)

	body := `{"role":2,"permission_id":"` + permID.String() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/admin/role-permissions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), permID.String())
}

func TestRolePermissionHandler_Delete(t *testing.T) {
	permID := uuid.New()
	rpID := uuid.New()
	server := newServerWithRolePermissions([]permission.RolePermission{
		{ID: rpID, Role: permission.RoleTrusted, PermissionID: permID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	})

	req := httptest.NewRequest(http.MethodDelete, "/admin/role-permissions/"+rpID.String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestRolePermissionHandler_Delete_NotFound(t *testing.T) {
	server := newServerWithRolePermissions(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/role-permissions/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
