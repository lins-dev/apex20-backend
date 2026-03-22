package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/permission"
)

type mockPermissionGetter struct {
	result permission.Permission
	err    error
}

func (m *mockPermissionGetter) GetPermissionByID(_ context.Context, _ uuid.UUID) (permission.Permission, error) {
	return m.result, m.err
}

func TestGetPermissionUseCase_Execute_ReturnsPermission(t *testing.T) {
	id := uuid.New()
	expected := permission.Permission{ID: id, Name: "token.move.own", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	uc := usecase.NewGetPermissionUseCase(&mockPermissionGetter{result: expected})

	got, err := uc.Execute(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.Name, got.Name)
}

func TestGetPermissionUseCase_Execute_ReturnsErrNotFound(t *testing.T) {
	uc := usecase.NewGetPermissionUseCase(&mockPermissionGetter{err: port.ErrNotFound})

	_, err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.True(t, errors.Is(err, port.ErrNotFound))
}

func TestGetPermissionUseCase_Execute_ReturnsRepoError(t *testing.T) {
	uc := usecase.NewGetPermissionUseCase(&mockPermissionGetter{err: errors.New("db error")})

	_, err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting permission")
}
