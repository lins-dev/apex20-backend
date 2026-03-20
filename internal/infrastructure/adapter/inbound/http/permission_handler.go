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

type permissionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toPermissionResponse(p permission.Permission) permissionResponse {
	return permissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// List

type listPermissionsOutput struct {
	Body struct {
		Permissions []permissionResponse `json:"permissions"`
	}
}

// Get

type getPermissionInput struct {
	ID string `path:"id"`
}

type getPermissionOutput struct {
	Body permissionResponse
}

// Create

type createPermissionInput struct {
	Body struct {
		Name        string `json:"name" minLength:"1" maxLength:"100"`
		Description string `json:"description"`
	}
}

type createPermissionOutput struct {
	Body permissionResponse
}

// Update

type updatePermissionInput struct {
	ID   string `path:"id"`
	Body struct {
		Name        string `json:"name" minLength:"1" maxLength:"100"`
		Description string `json:"description"`
	}
}

type updatePermissionOutput struct {
	Body permissionResponse
}

// Delete

type deletePermissionInput struct {
	ID string `path:"id"`
}

// RegisterPermissionHandler registers all /admin/permissions routes on the given API.
func RegisterPermissionHandler(api huma.API, repo port.PermissionRepository) {
	huma.Register(api, huma.Operation{
		OperationID: "list-permissions",
		Method:      http.MethodGet,
		Path:        "/admin/permissions",
		Summary:     "List Permissions",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, _ *struct{}) (*listPermissionsOutput, error) {
		perms, err := repo.ListPermissions(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list permissions", err)
		}
		out := &listPermissionsOutput{}
		out.Body.Permissions = make([]permissionResponse, len(perms))
		for i, p := range perms {
			out.Body.Permissions[i] = toPermissionResponse(p)
		}
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-permission",
		Method:      http.MethodGet,
		Path:        "/admin/permissions/{id}",
		Summary:     "Get Permission",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, input *getPermissionInput) (*getPermissionOutput, error) {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid permission id", err)
		}
		p, err := repo.GetPermissionByID(ctx, id)
		if err != nil {
			if errors.Is(err, port.ErrNotFound) {
				return nil, huma.Error404NotFound("permission not found")
			}
			return nil, huma.Error500InternalServerError("failed to get permission", err)
		}
		return &getPermissionOutput{Body: toPermissionResponse(p)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:  "create-permission",
		Method:       http.MethodPost,
		Path:         "/admin/permissions",
		Summary:      "Create Permission",
		Tags:         []string{"Admin"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *createPermissionInput) (*createPermissionOutput, error) {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to generate id", err)
		}
		now := time.Now()
		p := permission.Permission{
			ID:          id,
			Name:        input.Body.Name,
			Description: input.Body.Description,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := repo.CreatePermission(ctx, p); err != nil {
			return nil, huma.Error500InternalServerError("failed to create permission", err)
		}
		return &createPermissionOutput{Body: toPermissionResponse(p)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-permission",
		Method:      http.MethodPut,
		Path:        "/admin/permissions/{id}",
		Summary:     "Update Permission",
		Tags:        []string{"Admin"},
	}, func(ctx context.Context, input *updatePermissionInput) (*updatePermissionOutput, error) {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid permission id", err)
		}
		existing, err := repo.GetPermissionByID(ctx, id)
		if err != nil {
			if errors.Is(err, port.ErrNotFound) {
				return nil, huma.Error404NotFound("permission not found")
			}
			return nil, huma.Error500InternalServerError("failed to get permission", err)
		}
		existing.Name = input.Body.Name
		existing.Description = input.Body.Description
		existing.UpdatedAt = time.Now()
		if err := repo.UpdatePermission(ctx, existing); err != nil {
			return nil, huma.Error500InternalServerError("failed to update permission", err)
		}
		return &updatePermissionOutput{Body: toPermissionResponse(existing)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:  "delete-permission",
		Method:       http.MethodDelete,
		Path:         "/admin/permissions/{id}",
		Summary:      "Delete Permission",
		Tags:         []string{"Admin"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *deletePermissionInput) (*struct{}, error) {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid permission id", err)
		}
		deleted, err := repo.DeletePermission(ctx, id, time.Now())
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to delete permission", err)
		}
		if !deleted {
			return nil, huma.Error404NotFound("permission not found")
		}
		return nil, nil
	})
}
