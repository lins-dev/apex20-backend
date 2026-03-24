package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/user"
)

// UserUseCases groups the use cases needed for the user routes.
type UserUseCases struct {
	Get    port.UserGetter
	Update port.UserUpdater
	Delete port.UserDeleter
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Nick      string    `json:"nick,omitempty"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type getUserInput struct {
	ID uuid.UUID `path:"id"`
}

type getUserOutput struct {
	Body userResponse
}

type updateUserInput struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Name string  `json:"name" minLength:"1" maxLength:"255"`
		Nick *string `json:"nick,omitempty" maxLength:"255"`
	}
}

type updateUserOutput struct {
	Body userResponse
}

type deleteUserInput struct {
	ID uuid.UUID `path:"id"`
}

// RegisterUserHandler registers all /users routes on the given API.
func RegisterUserHandler(api huma.API, uc UserUseCases) {
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/users/{id}",
		Summary:     "Get User",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *getUserInput) (*getUserOutput, error) {
		u, err := uc.Get.Execute(ctx, input.ID)
		if err != nil {
			if errors.Is(err, port.ErrNotFound) {
				return nil, huma.Error404NotFound("user not found")
			}
			return nil, huma.Error500InternalServerError("failed to get user", err)
		}
		return &getUserOutput{Body: toUserResponse(u)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodPatch,
		Path:        "/users/{id}",
		Summary:     "Update User",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *updateUserInput) (*updateUserOutput, error) {
		u, err := uc.Update.Execute(ctx, port.UpdateUserInput{
			ID:   input.ID,
			Name: input.Body.Name,
			Nick: input.Body.Nick,
		})
		if err != nil {
			if errors.Is(err, port.ErrNotFound) {
				return nil, huma.Error404NotFound("user not found")
			}
			return nil, huma.Error500InternalServerError("failed to update user", err)
		}
		return &updateUserOutput{Body: toUserResponse(u)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-user",
		Method:        http.MethodDelete,
		Path:          "/users/{id}",
		Summary:       "Delete User",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *deleteUserInput) (*struct{}, error) {
		if err := uc.Delete.Execute(ctx, input.ID); err != nil {
			if errors.Is(err, port.ErrNotFound) {
				return nil, huma.Error404NotFound("user not found")
			}
			return nil, huma.Error500InternalServerError("failed to delete user", err)
		}
		return nil, nil
	})
}

func toUserResponse(u user.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Nick:      u.Nick,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}