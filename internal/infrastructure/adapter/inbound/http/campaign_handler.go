package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/campaign"
	"github.com/apex20/backend/internal/infrastructure/adapter/inbound/http/middleware"
)

// CampaignUseCases agrupa os use cases necessários para as rotas de campanha.
type CampaignUseCases struct {
	Create port.CampaignCreator
	List   port.CampaignLister
	Get    port.CampaignGetter
	Update port.CampaignUpdater
	Delete port.CampaignDeleter
}

type campaignResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type createCampaignInput struct {
	Body struct {
		Name        string `json:"name" minLength:"1" maxLength:"255"`
		Description string `json:"description"`
	}
}

type createCampaignOutput struct {
	Body campaignResponse
}

type listCampaignsInput struct{}

type listCampaignsOutput struct {
	Body []campaignResponse
}

type getCampaignInput struct {
	ID uuid.UUID `path:"id"`
}

type getCampaignOutput struct {
	Body campaignResponse
}

type updateCampaignInput struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Name        string  `json:"name" minLength:"1" maxLength:"255"`
		Description *string `json:"description,omitempty"`
	}
}

type updateCampaignOutput struct {
	Body campaignResponse
}

type deleteCampaignInput struct {
	ID uuid.UUID `path:"id"`
}

// RegisterCampaignHandler registers all /campaigns routes on the given API.
func RegisterCampaignHandler(api huma.API, uc CampaignUseCases) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-campaign",
		Method:        http.MethodPost,
		Path:          "/campaigns",
		Summary:       "Create Campaign",
		Tags:          []string{"Campaigns"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *createCampaignInput) (*createCampaignOutput, error) {
		claims, ok := middleware.ClaimsFromContext(ctx)
		if !ok {
			return nil, huma.Error401Unauthorized("authentication required")
		}
		c, err := uc.Create.Execute(ctx, port.CreateCampaignInput{
			UserID:      claims.UserID,
			Name:        input.Body.Name,
			Description: input.Body.Description,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create campaign", err)
		}
		return &createCampaignOutput{Body: fromCampaign(c)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-campaigns",
		Method:      http.MethodGet,
		Path:        "/campaigns",
		Summary:     "List Campaigns",
		Tags:        []string{"Campaigns"},
	}, func(ctx context.Context, _ *listCampaignsInput) (*listCampaignsOutput, error) {
		claims, ok := middleware.ClaimsFromContext(ctx)
		if !ok {
			return nil, huma.Error401Unauthorized("authentication required")
		}
		campaigns, err := uc.List.Execute(ctx, claims.UserID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list campaigns", err)
		}
		resp := make([]campaignResponse, len(campaigns))
		for i, c := range campaigns {
			resp[i] = fromCampaign(c)
		}
		return &listCampaignsOutput{Body: resp}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-campaign",
		Method:      http.MethodGet,
		Path:        "/campaigns/{id}",
		Summary:     "Get Campaign",
		Tags:        []string{"Campaigns"},
	}, func(ctx context.Context, input *getCampaignInput) (*getCampaignOutput, error) {
		if _, ok := middleware.ClaimsFromContext(ctx); !ok {
			return nil, huma.Error401Unauthorized("authentication required")
		}
		c, err := uc.Get.Execute(ctx, input.ID)
		if err != nil {
			if errors.Is(err, port.ErrNotFound) {
				return nil, huma.Error404NotFound("campaign not found")
			}
			return nil, huma.Error500InternalServerError("failed to get campaign", err)
		}
		return &getCampaignOutput{Body: fromCampaign(c)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-campaign",
		Method:      http.MethodPut,
		Path:        "/campaigns/{id}",
		Summary:     "Update Campaign",
		Tags:        []string{"Campaigns"},
	}, func(ctx context.Context, input *updateCampaignInput) (*updateCampaignOutput, error) {
		if _, ok := middleware.ClaimsFromContext(ctx); !ok {
			return nil, huma.Error401Unauthorized("authentication required")
		}
		c, err := uc.Update.Execute(ctx, port.UpdateCampaignInput{
			ID:          input.ID,
			Name:        input.Body.Name,
			Description: input.Body.Description,
		})
		if err != nil {
			if errors.Is(err, port.ErrNotFound) {
				return nil, huma.Error404NotFound("campaign not found")
			}
			return nil, huma.Error500InternalServerError("failed to update campaign", err)
		}
		return &updateCampaignOutput{Body: fromCampaign(c)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-campaign",
		Method:        http.MethodDelete,
		Path:          "/campaigns/{id}",
		Summary:       "Delete Campaign",
		Tags:          []string{"Campaigns"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *deleteCampaignInput) (*struct{}, error) {
		if _, ok := middleware.ClaimsFromContext(ctx); !ok {
			return nil, huma.Error401Unauthorized("authentication required")
		}
		if err := uc.Delete.Execute(ctx, input.ID); err != nil {
			if errors.Is(err, port.ErrNotFound) {
				return nil, huma.Error404NotFound("campaign not found")
			}
			return nil, huma.Error500InternalServerError("failed to delete campaign", err)
		}
		return nil, nil
	})
}

func fromCampaign(c campaign.Campaign) campaignResponse {
	return campaignResponse{
		ID:          c.ID,
		UserID:      c.UserID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
