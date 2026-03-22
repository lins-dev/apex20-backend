package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/permission"
)

func TestListRolesUseCase_Execute_ReturnsAllRoles(t *testing.T) {
	uc := usecase.NewListRolesUseCase()

	got, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Contains(t, got, permission.RoleGM)
	assert.Contains(t, got, permission.RolePlayer)
	assert.Contains(t, got, permission.RoleTrusted)
}
