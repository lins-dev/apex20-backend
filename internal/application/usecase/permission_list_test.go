package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/permission"
)

type mockPermissionLister struct {
	result []permission.Permission
	err    error
}

func (m *mockPermissionLister) ListPermissions(_ context.Context) ([]permission.Permission, error) {
	return m.result, m.err
}

func TestListPermissionsUseCase_Execute_ReturnsList(t *testing.T) {
	perms := []permission.Permission{
		{ID: uuid.New(), Name: "chat.send", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "chat.roll", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	uc := usecase.NewListPermissionsUseCase(&mockPermissionLister{result: perms})

	got, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "chat.send", got[0].Name)
}

func TestListPermissionsUseCase_Execute_ReturnsError(t *testing.T) {
	uc := usecase.NewListPermissionsUseCase(&mockPermissionLister{err: errors.New("db error")})

	got, err := uc.Execute(context.Background())

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "listing permissions")
}
