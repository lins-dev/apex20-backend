package http

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
)

// CampaignUseCases agrupa os use cases necessários para as rotas de campanha.
type CampaignUseCases struct {
	Create port.CampaignCreator
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
		UserID      uuid.UUID `json:"user_id"`
		Name        string    `json:"name" minLength:"1" maxLength:"255"`
		Description string    `json:"description"`
	}
}

type createCampaignOutput struct {
	Body campaignResponse
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
		c, err := uc.Create.Execute(ctx, port.CreateCampaignInput{
			UserID:      input.Body.UserID,
			Name:        input.Body.Name,
			Description: input.Body.Description,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create campaign", err)
		}
		return &createCampaignOutput{Body: campaignResponse{
			ID:          c.ID,
			UserID:      c.UserID,
			Name:        c.Name,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		}}, nil
	})
}
