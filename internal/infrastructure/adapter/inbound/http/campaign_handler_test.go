package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/campaign"
	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
)

type stubCampaignWithMemberCreator struct{}

func (s *stubCampaignWithMemberCreator) CreateCampaignWithMember(_ context.Context, _ campaign.Campaign, _ campaign.Member) error {
	return nil
}

func newServerWithCampaigns() *adapter.ChiServer {
	server := adapter.NewChiServer()
	adapter.RegisterCampaignHandler(server.GetAPI(), adapter.CampaignUseCases{
		Create: usecase.NewCreateCampaignUseCase(&stubCampaignWithMemberCreator{}),
	})
	return server
}

func TestCampaignHandler_Create_ReturnsCreatedCampaign(t *testing.T) {
	server := newServerWithCampaigns()
	userID := uuid.New()

	body := `{"user_id":"` + userID.String() + `","name":"Campanha Teste","description":"Uma desc"}`
	req := httptest.NewRequest(http.MethodPost, "/campaigns", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "Campanha Teste", resp["name"])
	assert.Equal(t, "Uma desc", resp["description"])
	assert.Equal(t, userID.String(), resp["user_id"])
	assert.NotEmpty(t, resp["id"])
}

func TestCampaignHandler_Create_ReturnsBadRequestOnMissingName(t *testing.T) {
	server := newServerWithCampaigns()

	body := `{"user_id":"` + uuid.New().String() + `","name":""}`
	req := httptest.NewRequest(http.MethodPost, "/campaigns", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}
