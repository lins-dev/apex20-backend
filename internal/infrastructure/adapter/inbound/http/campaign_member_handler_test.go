package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/campaign"
	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
)

// --- stubs ---

type stubMemberCreator struct{}

func (s *stubMemberCreator) CreateCampaignMember(_ context.Context, _ campaign.Member) error {
	return nil
}

type stubMemberDeleter struct {
	err error
}

func (s *stubMemberDeleter) DeleteCampaignMember(_ context.Context, _, _ uuid.UUID) error {
	return s.err
}

func newServerWithMembers(deleterErr error) *adapter.ChiServer {
	server := adapter.NewChiServer()
	adapter.RegisterCampaignMemberHandler(server.GetAPI(), adapter.CampaignMemberUseCases{
		Invite: usecase.NewInviteMemberUseCase(&stubMemberCreator{}),
		Remove: usecase.NewRemoveMemberUseCase(&stubMemberDeleter{err: deleterErr}),
	})
	return server
}

// --- POST /campaigns/{id}/members ---

func TestCampaignMemberHandler_Invite_ReturnsCreatedMember(t *testing.T) {
	server := newServerWithMembers(nil)
	campaignID := uuid.New()
	userID := uuid.New()

	body := `{"user_id":"` + userID.String() + `","role":1}`
	req := httptest.NewRequest(http.MethodPost, "/campaigns/"+campaignID.String()+"/members", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestCampaignMemberHandler_Invite_ReturnsBadRequestOnMissingUserID(t *testing.T) {
	server := newServerWithMembers(nil)

	body := `{"role":1}`
	req := httptest.NewRequest(http.MethodPost, "/campaigns/"+uuid.New().String()+"/members", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

// --- DELETE /campaigns/{id}/members/{userId} ---

func TestCampaignMemberHandler_Remove_ReturnsNoContent(t *testing.T) {
	server := newServerWithMembers(nil)

	req := httptest.NewRequest(http.MethodDelete,
		"/campaigns/"+uuid.New().String()+"/members/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestCampaignMemberHandler_Remove_ReturnsNotFound(t *testing.T) {
	server := newServerWithMembers(port.ErrNotFound)

	req := httptest.NewRequest(http.MethodDelete,
		"/campaigns/"+uuid.New().String()+"/members/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
