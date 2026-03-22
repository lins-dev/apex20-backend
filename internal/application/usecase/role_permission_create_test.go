package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/permission"
)

type mockRolePermissionCreator struct {
	err     error
	created []permission.RolePermission
}

func (m *mockRolePermissionCreator) CreateRolePermission(_ context.Context, rp permission.RolePermission) error {
	if m.err != nil {
		return m.err
	}
	m.created = append(m.created, rp)
	return nil
}

func TestCreateRolePermissionUseCase_Execute_ReturnsCreatedRolePermission(t *testing.T) {
	permID := uuid.New()
	repo := &mockRolePermissionCreator{}
	uc := usecase.NewCreateRolePermissionUseCase(repo)

	got, err := uc.Execute(context.Background(), port.CreateRolePermissionInput{
		Role:         permission.RolePlayer,
		PermissionID: permID,
	})

	require.NoError(t, err)
	assert.NotEqual(t, "00000000-0000-0000-0000-000000000000", got.ID.String())
	assert.Equal(t, permission.RolePlayer, got.Role)
	assert.Equal(t, permID, got.PermissionID)
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())
	require.Len(t, repo.created, 1)
	assert.Equal(t, got.ID, repo.created[0].ID)
}

func TestCreateRolePermissionUseCase_Execute_ReturnsErrorOnRepoFailure(t *testing.T) {
	uc := usecase.NewCreateRolePermissionUseCase(&mockRolePermissionCreator{err: errors.New("insert failed")})

	_, err := uc.Execute(context.Background(), port.CreateRolePermissionInput{
		Role:         permission.RoleGM,
		PermissionID: uuid.New(),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "creating role permission")
}
