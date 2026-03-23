package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/apex20/backend/internal/application/port"
	apexv1 "github.com/apex20/contracts/proto/apex20/v1"
)

// CampaignMemberUseCases agrupa os use cases necessários para as rotas de membros.
type CampaignMemberUseCases struct {
	Invite port.MemberInviter
	Remove port.MemberRemover
}

type memberResponse struct {
	ID         uuid.UUID     `json:"id"`
	CampaignID uuid.UUID     `json:"campaign_id"`
	UserID     uuid.UUID     `json:"user_id"`
	Role       apexv1.Role   `json:"role"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type inviteMemberInput struct {
	CampaignID uuid.UUID `path:"id"`
	Body       struct {
		UserID uuid.UUID   `json:"user_id"`
		Role   apexv1.Role `json:"role"`
	}
}

type inviteMemberOutput struct {
	Body memberResponse
}

type removeMemberInput struct {
	CampaignID uuid.UUID `path:"id"`
	UserID     uuid.UUID `path:"userId"`
}

// RegisterCampaignMemberHandler registers /campaigns/{id}/members routes on the given API.
func RegisterCampaignMemberHandler(api huma.API, uc CampaignMemberUseCases) {
	huma.Register(api, huma.Operation{
		OperationID:   "invite-member",
		Method:        http.MethodPost,
		Path:          "/campaigns/{id}/members",
		Summary:       "Invite Member",
		Tags:          []string{"Campaign Members"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *inviteMemberInput) (*inviteMemberOutput, error) {
		m, err := uc.Invite.Execute(ctx, port.InviteMemberInput{
			CampaignID: input.CampaignID,
			UserID:     input.Body.UserID,
			Role:       input.Body.Role,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to invite member", err)
		}
		return &inviteMemberOutput{Body: memberResponse{
			ID:         m.ID,
			CampaignID: m.CampaignID,
			UserID:     m.UserID,
			Role:       m.Role,
			CreatedAt:  m.CreatedAt,
			UpdatedAt:  m.UpdatedAt,
		}}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "remove-member",
		Method:        http.MethodDelete,
		Path:          "/campaigns/{id}/members/{userId}",
		Summary:       "Remove Member",
		Tags:          []string{"Campaign Members"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *removeMemberInput) (*struct{}, error) {
		if err := uc.Remove.Execute(ctx, input.CampaignID, input.UserID); err != nil {
			if errors.Is(err, port.ErrNotFound) {
				return nil, huma.Error404NotFound("member not found")
			}
			return nil, huma.Error500InternalServerError("failed to remove member", err)
		}
		return nil, nil
	})
}
