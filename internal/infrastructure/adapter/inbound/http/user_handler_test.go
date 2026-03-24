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
	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
	"github.com/apex20/backend/internal/domain/user"
)

// --- stubs ---

type stubUserGetter struct {
	u   user.User
	err error
}

func (s *stubUserGetter) Execute(_ context.Context, _ uuid.UUID) (user.User, error) {
	return s.u, s.err
}

type stubUserUpdater struct {
	u   user.User
	err error
}

func (s *stubUserUpdater) Execute(_ context.Context, _ port.UpdateUserInput) (user.User, error) {
	return s.u, s.err
}

type stubUserDeleter struct{ err error }

func (s *stubUserDeleter) Execute(_ context.Context, _ uuid.UUID) error { return s.err }

// --- helper ---

func newServerWithUsers(getter port.UserGetter, updater port.UserUpdater, deleter port.UserDeleter) *adapter.ChiServer {
	server := adapter.NewChiServer()
	adapter.RegisterUserHandler(server.GetAPI(), adapter.UserUseCases{
		Get:    getter,
		Update: updater,
		Delete: deleter,
	})
	return server
}

// --- tests ---

func TestUserHandler_Get_ReturnsUser(t *testing.T) {
	id := uuid.New()
	server := newServerWithUsers(
		&stubUserGetter{u: user.User{ID: id, Name: "Hero", Email: "hero@apex20.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		&stubUserUpdater{},
		&stubUserDeleter{},
	)

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, id.String(), resp["id"])
	assert.Equal(t, "Hero", resp["name"])
}

func TestUserHandler_Get_ReturnsNotFound(t *testing.T) {
	server := newServerWithUsers(
		&stubUserGetter{err: port.ErrNotFound},
		&stubUserUpdater{},
		&stubUserDeleter{},
	)

	req := httptest.NewRequest(http.MethodGet, "/users/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUserHandler_Update_ReturnsUpdatedUser(t *testing.T) {
	id := uuid.New()
	nick := "dragonslayer"
	server := newServerWithUsers(
		&stubUserGetter{},
		&stubUserUpdater{u: user.User{ID: id, Name: "Hero Updated", Nick: nick, UpdatedAt: time.Now()}},
		&stubUserDeleter{},
	)

	body := `{"name":"Hero Updated","nick":"dragonslayer"}`
	req := httptest.NewRequest(http.MethodPatch, "/users/"+id.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "Hero Updated", resp["name"])
	assert.Equal(t, "dragonslayer", resp["nick"])
}

func TestUserHandler_Update_ReturnsNotFound(t *testing.T) {
	server := newServerWithUsers(
		&stubUserGetter{},
		&stubUserUpdater{err: port.ErrNotFound},
		&stubUserDeleter{},
	)

	body := `{"name":"Hero"}`
	req := httptest.NewRequest(http.MethodPatch, "/users/"+uuid.New().String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUserHandler_Delete_ReturnsNoContent(t *testing.T) {
	server := newServerWithUsers(
		&stubUserGetter{},
		&stubUserUpdater{},
		&stubUserDeleter{},
	)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestUserHandler_Delete_ReturnsNotFound(t *testing.T) {
	server := newServerWithUsers(
		&stubUserGetter{},
		&stubUserUpdater{},
		&stubUserDeleter{err: port.ErrNotFound},
	)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
