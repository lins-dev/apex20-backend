package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/permission"
)

// RolePermissionUseCases agrupa os use cases necessários para as rotas de role-permission.
type RolePermissionUseCases struct {
	List   port.RolePermissionLister
	Get    port.RolePermissionGetter
	Create port.RolePermissionCreator
	Delete port.RolePermissionDeleter
}

type rolePermissionResponse struct {
	ID           uuid.UUID       `json:"id"`
	Role         permission.Role `json:"role"`
	PermissionID uuid.UUID       `json:"permission_id"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func toRolePermissionResponse(rp permission.RolePermission) rolePermissionResponse {
	return rolePermissionResponse{
		ID:           rp.ID,
		Role:         rp.Role,
		PermissionID: rp.PermissionID,
		CreatedAt:    rp.CreatedAt,
		UpdatedAt:    rp.UpdatedAt,
	}
}

// List

type listRolePermissionsOutput struct {
	Body struct {
		RolePermissions []rolePermissionResponse `json:"role_permissions"`
	}
}

// Get

type getRolePermissionInput struct {
	ID string `path:"id"`
}

type getRolePermissionOutput struct {
	Body rolePermissionResponse
}

// Create

type createRolePermissionInput struct {
	Body struct {
		Role         permission.Role `json:"role"`
		PermissionID uuid.UUID       `json:"permission_id"`
	}
}

type createRolePermissionOutput struct {
	Body rolePermissionResponse
}

// Delete

type deleteRolePermissionInput struct {
	ID string `path:"id"`
}

// RegisterRolePermissionHandler registers all /admin/role-permissions routes on the given API.
func RegisterRolePermissionHandler(api huma.API, uc RolePermissionUseCases) {
	huma.Register(api, huma.Operation{
		OperationID: "list-role-permissions",
		Method:      http.MethodGet,
		Path:        "/admin/role-permissions",
		Summary:     "List Role Permissions",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, _ *struct{}) (*listRolePermissionsOutput, error) {
		rps, err := uc.List.Execute(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list role permissions", err)
		}
		out := &listRolePermissionsOutput{}
		out.Body.RolePermissions = make([]rolePermissionResponse, len(rps))
		for i, rp := range rps {
			out.Body.RolePermissions[i] = toRolePermissionResponse(rp)
		}
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-role-permission",
		Method:      http.MethodGet,
		Path:        "/admin/role-permissions/{id}",
		Summary:     "Get Role Permission",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, input *getRolePermissionInput) (*getRolePermissionOutput, error) {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid role permission id", err)
		}
		rp, err := uc.Get.Execute(ctx, id)
		if err != nil {
			if errors.Is(err, port.ErrNotFound) {
				return nil, huma.Error404NotFound("role permission not found")
			}
			return nil, huma.Error500InternalServerError("failed to get role permission", err)
		}
		return &getRolePermissionOutput{Body: toRolePermissionResponse(rp)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "create-role-permission",
		Method:        http.MethodPost,
		Path:          "/admin/role-permissions",
		Summary:       "Create Role Permission",
		Tags:          []string{"Admin"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *createRolePermissionInput) (*createRolePermissionOutput, error) {
		rp, err := uc.Create.Execute(ctx, port.CreateRolePermissionInput{
			Role:         input.Body.Role,
			PermissionID: input.Body.PermissionID,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create role permission", err)
		}
		return &createRolePermissionOutput{Body: toRolePermissionResponse(rp)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-role-permission",
		Method:        http.MethodDelete,
		Path:          "/admin/role-permissions/{id}",
		Summary:       "Delete Role Permission",
		Tags:          []string{"Admin"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *deleteRolePermissionInput) (*struct{}, error) {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid role permission id", err)
		}
		deleted, err := uc.Delete.Execute(ctx, id)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to delete role permission", err)
		}
		if !deleted {
			return nil, huma.Error404NotFound("role permission not found")
		}
		return nil, nil
	})
}
