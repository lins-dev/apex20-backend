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

type mockPermissionGetterUpdater struct {
	getResult   permission.Permission
	getErr      error
	updateErr   error
	updatedPerm *permission.Permission
}

func (m *mockPermissionGetterUpdater) GetPermissionByID(_ context.Context, _ uuid.UUID) (permission.Permission, error) {
	return m.getResult, m.getErr
}

func (m *mockPermissionGetterUpdater) UpdatePermission(_ context.Context, p permission.Permission) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.updatedPerm = &p
	return nil
}

func TestUpdatePermissionUseCase_Execute_ReturnsUpdatedPermission(t *testing.T) {
	id := uuid.New()
	existing := permission.Permission{ID: id, Name: "old.name", Description: "Old", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	repo := &mockPermissionGetterUpdater{getResult: existing}
	uc := usecase.NewUpdatePermissionUseCase(repo)

	got, err := uc.Execute(context.Background(), port.UpdatePermissionInput{
		ID:          id,
		Name:        "new.name",
		Description: "New",
	})

	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "new.name", got.Name)
	assert.Equal(t, "New", got.Description)
	assert.True(t, got.UpdatedAt.After(existing.UpdatedAt) || got.UpdatedAt.Equal(existing.UpdatedAt))
	require.NotNil(t, repo.updatedPerm)
	assert.Equal(t, "new.name", repo.updatedPerm.Name)
}

func TestUpdatePermissionUseCase_Execute_ReturnsErrNotFound(t *testing.T) {
	uc := usecase.NewUpdatePermissionUseCase(&mockPermissionGetterUpdater{getErr: port.ErrNotFound})

	_, err := uc.Execute(context.Background(), port.UpdatePermissionInput{ID: uuid.New(), Name: "x"})

	require.Error(t, err)
	assert.True(t, errors.Is(err, port.ErrNotFound))
}

func TestUpdatePermissionUseCase_Execute_ReturnsGetError(t *testing.T) {
	uc := usecase.NewUpdatePermissionUseCase(&mockPermissionGetterUpdater{getErr: errors.New("db error")})

	_, err := uc.Execute(context.Background(), port.UpdatePermissionInput{ID: uuid.New(), Name: "x"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting permission")
}

func TestUpdatePermissionUseCase_Execute_ReturnsUpdateError(t *testing.T) {
	repo := &mockPermissionGetterUpdater{
		getResult: permission.Permission{ID: uuid.New(), Name: "x", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		updateErr: errors.New("update failed"),
	}
	uc := usecase.NewUpdatePermissionUseCase(repo)

	_, err := uc.Execute(context.Background(), port.UpdatePermissionInput{ID: uuid.New(), Name: "x"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "updating permission")
}
