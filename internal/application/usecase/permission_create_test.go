package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/permission"
)

type mockPermissionCreator struct {
	err     error
	created []permission.Permission
}

func (m *mockPermissionCreator) CreatePermission(_ context.Context, p permission.Permission) error {
	if m.err != nil {
		return m.err
	}
	m.created = append(m.created, p)
	return nil
}

func TestCreatePermissionUseCase_Execute_ReturnsCreatedPermission(t *testing.T) {
	repo := &mockPermissionCreator{}
	uc := usecase.NewCreatePermissionUseCase(repo)

	got, err := uc.Execute(context.Background(), port.CreatePermissionInput{
		Name:        "scene.manage",
		Description: "Gerenciar cenas",
	})

	require.NoError(t, err)
	assert.NotEqual(t, "00000000-0000-0000-0000-000000000000", got.ID.String())
	assert.Equal(t, "scene.manage", got.Name)
	assert.Equal(t, "Gerenciar cenas", got.Description)
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())
	require.Len(t, repo.created, 1)
	assert.Equal(t, got.ID, repo.created[0].ID)
}

func TestCreatePermissionUseCase_Execute_ReturnsErrorOnRepoFailure(t *testing.T) {
	uc := usecase.NewCreatePermissionUseCase(&mockPermissionCreator{err: errors.New("insert failed")})

	_, err := uc.Execute(context.Background(), port.CreatePermissionInput{Name: "x"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "creating permission")
}
