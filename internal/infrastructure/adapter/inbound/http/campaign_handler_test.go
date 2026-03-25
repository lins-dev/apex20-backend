package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/domain/campaign"
	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
	"github.com/apex20/backend/internal/infrastructure/adapter/inbound/http/middleware"
)

// --- stubs for create ---

type stubCampaignWithMemberCreator struct{}

func (s *stubCampaignWithMemberCreator) CreateCampaignWithMember(_ context.Context, _ campaign.Campaign, _ campaign.Member) error {
	return nil
}

// --- stubs for list/get/update/delete ---

type stubCampaignsByUserIDLister struct {
	campaigns []campaign.Campaign
}

func (s *stubCampaignsByUserIDLister) ListCampaignsByUserID(_ context.Context, _ uuid.UUID) ([]campaign.Campaign, error) {
	return s.campaigns, nil
}

type stubCampaignByIDGetter struct {
	c   campaign.Campaign
	err error
}

func (s *stubCampaignByIDGetter) GetCampaignByID(_ context.Context, _ uuid.UUID) (campaign.Campaign, error) {
	return s.c, s.err
}

type stubCampaignUpdater struct {
	err error
}

func (s *stubCampaignUpdater) UpdateCampaign(_ context.Context, id uuid.UUID, name string, description *string) (campaign.Campaign, error) {
	if s.err != nil {
		return campaign.Campaign{}, s.err
	}
	desc := ""
	if description != nil {
		desc = *description
	}
	return campaign.Campaign{ID: id, Name: name, Description: desc, UpdatedAt: time.Now()}, nil
}

type stubCampaignDeleter struct {
	err error
}

func (s *stubCampaignDeleter) DeleteCampaign(_ context.Context, _ uuid.UUID) error {
	return s.err
}

// --- server builder ---

func newServerWithCampaigns() *adapter.ChiServer {
	return newServerWithCampaignStubs(
		&stubCampaignsByUserIDLister{},
		&stubCampaignByIDGetter{},
		&stubCampaignUpdater{},
		&stubCampaignDeleter{},
	)
}

func newServerWithCampaignStubs(
	lister *stubCampaignsByUserIDLister,
	getter *stubCampaignByIDGetter,
	updater *stubCampaignUpdater,
	deleter *stubCampaignDeleter,
) *adapter.ChiServer {
	server := adapter.NewChiServer()
	adapter.RegisterCampaignHandler(server.GetAPI(), adapter.CampaignUseCases{
		Create: usecase.NewCreateCampaignUseCase(&stubCampaignWithMemberCreator{}),
		List:   usecase.NewListCampaignsUseCase(lister),
		Get:    usecase.NewGetCampaignUseCase(getter),
		Update: usecase.NewUpdateCampaignUseCase(updater),
		Delete: usecase.NewDeleteCampaignUseCase(deleter),
	})
	return server
}

// withAuth injects auth claims into the request context, simulating the JWT middleware in unit tests.
func withAuth(r *http.Request, userID uuid.UUID) *http.Request {
	return r.WithContext(middleware.WithAuthClaims(r.Context(), port.AuthClaims{UserID: userID}))
}

// --- tests ---

func TestCampaignHandler_Create_ReturnsCreatedCampaign(t *testing.T) {
	server := newServerWithCampaigns()
	userID := uuid.New()

	body := `{"name":"Campanha Teste","description":"Uma desc"}`
	req := httptest.NewRequest(http.MethodPost, "/campaigns", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAuth(req, userID)
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

	body := `{"name":""}`
	req := httptest.NewRequest(http.MethodPost, "/campaigns", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAuth(req, uuid.New())
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestCampaignHandler_Create_Returns401WhenNotAuthenticated(t *testing.T) {
	server := newServerWithCampaigns()

	body := `{"name":"Campanha","description":""}`
	req := httptest.NewRequest(http.MethodPost, "/campaigns", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCampaignHandler_List_ReturnsCampaigns(t *testing.T) {
	userID := uuid.New()
	campaigns := []campaign.Campaign{
		{ID: uuid.New(), UserID: userID, Name: "C1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), UserID: userID, Name: "C2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	server := newServerWithCampaignStubs(
		&stubCampaignsByUserIDLister{campaigns: campaigns},
		&stubCampaignByIDGetter{},
		&stubCampaignUpdater{},
		&stubCampaignDeleter{},
	)

	req := httptest.NewRequest(http.MethodGet, "/campaigns", nil)
	req = withAuth(req, userID)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var resp []map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Len(t, resp, 2)
}

func TestCampaignHandler_List_Returns401WhenNotAuthenticated(t *testing.T) {
	server := newServerWithCampaigns()

	req := httptest.NewRequest(http.MethodGet, "/campaigns", nil)
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCampaignHandler_Get_ReturnsCampaign(t *testing.T) {
	id := uuid.New()
	server := newServerWithCampaignStubs(
		&stubCampaignsByUserIDLister{},
		&stubCampaignByIDGetter{c: campaign.Campaign{ID: id, Name: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		&stubCampaignUpdater{},
		&stubCampaignDeleter{},
	)

	req := httptest.NewRequest(http.MethodGet, "/campaigns/"+id.String(), nil)
	req = withAuth(req, uuid.New())
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, id.String(), resp["id"])
	assert.Equal(t, "Test", resp["name"])
}

func TestCampaignHandler_Get_ReturnsNotFound(t *testing.T) {
	server := newServerWithCampaignStubs(
		&stubCampaignsByUserIDLister{},
		&stubCampaignByIDGetter{err: port.ErrNotFound},
		&stubCampaignUpdater{},
		&stubCampaignDeleter{},
	)

	req := httptest.NewRequest(http.MethodGet, "/campaigns/"+uuid.New().String(), nil)
	req = withAuth(req, uuid.New())
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCampaignHandler_Update_ReturnsUpdatedCampaign(t *testing.T) {
	id := uuid.New()
	server := newServerWithCampaigns()

	body := `{"name":"Novo Nome","description":"Nova Desc"}`
	req := httptest.NewRequest(http.MethodPut, "/campaigns/"+id.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAuth(req, uuid.New())
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "Novo Nome", resp["name"])
	assert.Equal(t, "Nova Desc", resp["description"])
}

func TestCampaignHandler_Update_ReturnsNotFound(t *testing.T) {
	server := newServerWithCampaignStubs(
		&stubCampaignsByUserIDLister{},
		&stubCampaignByIDGetter{},
		&stubCampaignUpdater{err: port.ErrNotFound},
		&stubCampaignDeleter{},
	)

	body := `{"name":"X"}`
	req := httptest.NewRequest(http.MethodPut, "/campaigns/"+uuid.New().String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAuth(req, uuid.New())
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCampaignHandler_Update_ReturnsBadRequestOnMissingName(t *testing.T) {
	server := newServerWithCampaigns()

	body := `{"name":""}`
	req := httptest.NewRequest(http.MethodPut, "/campaigns/"+uuid.New().String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withAuth(req, uuid.New())
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestCampaignHandler_Delete_ReturnsNoContent(t *testing.T) {
	server := newServerWithCampaigns()

	req := httptest.NewRequest(http.MethodDelete, "/campaigns/"+uuid.New().String(), nil)
	req = withAuth(req, uuid.New())
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestCampaignHandler_Delete_ReturnsNotFound(t *testing.T) {
	server := newServerWithCampaignStubs(
		&stubCampaignsByUserIDLister{},
		&stubCampaignByIDGetter{},
		&stubCampaignUpdater{},
		&stubCampaignDeleter{err: port.ErrNotFound},
	)

	req := httptest.NewRequest(http.MethodDelete, "/campaigns/"+uuid.New().String(), nil)
	req = withAuth(req, uuid.New())
	rr := httptest.NewRecorder()
	server.GetHandler().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
