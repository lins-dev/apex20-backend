package http_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/apex20/contracts/proto/apex20/v1"
	"github.com/apex20/contracts/proto/apex20/v1/apex20v1connect"
	"github.com/apex20/backend/internal/application/port"
	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
	"github.com/apex20/backend/internal/domain/user"
)

// --- stubs ---

type stubSignUpper struct {
	output port.SignUpOutput
	err    error
}

func (s *stubSignUpper) Execute(_ context.Context, _ port.SignUpInput) (port.SignUpOutput, error) {
	return s.output, s.err
}

type stubSignIner struct {
	output port.SignInOutput
	err    error
}

func (s *stubSignIner) Execute(_ context.Context, _ port.SignInInput) (port.SignInOutput, error) {
	return s.output, s.err
}

// --- helpers ---

func newAuthServer(signUpper port.UserSignUpper, signIner port.UserSignIner) *httptest.Server {
	server := adapter.NewChiServer()
	adapter.RegisterAuthHandler(server, adapter.AuthUseCases{
		SignUp: signUpper,
		SignIn: signIner,
	})
	return httptest.NewServer(server.GetHandler())
}

// --- tests ---

func TestAuthHandler_SignUp_ReturnsAccessToken(t *testing.T) {
	userID := uuid.New()
	ts := newAuthServer(
		&stubSignUpper{output: port.SignUpOutput{
			User:        user.User{ID: userID, Email: "hero@apex20.com", Name: "Hero"},
			AccessToken: "jwt.token.here",
		}},
		&stubSignIner{},
	)
	defer ts.Close()

	client := apex20v1connect.NewAuthServiceClient(ts.Client(), ts.URL)
	resp, err := client.SignUp(context.Background(), connect.NewRequest(&v1.SignUpRequest{
		Email:    "hero@apex20.com",
		Password: "secret123",
		Name:     "Hero",
	}))

	require.NoError(t, err)
	assert.Equal(t, userID.String(), resp.Msg.UserId)
	assert.Equal(t, "jwt.token.here", resp.Msg.AccessToken)
}

func TestAuthHandler_SignUp_ReturnsAlreadyExistsOnDuplicateEmail(t *testing.T) {
	ts := newAuthServer(
		&stubSignUpper{err: port.ErrEmailAlreadyExists},
		&stubSignIner{},
	)
	defer ts.Close()

	client := apex20v1connect.NewAuthServiceClient(ts.Client(), ts.URL)
	_, err := client.SignUp(context.Background(), connect.NewRequest(&v1.SignUpRequest{
		Email:    "existing@apex20.com",
		Password: "secret123",
		Name:     "Hero",
	}))

	require.Error(t, err)
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeAlreadyExists, connectErr.Code())
}

func TestAuthHandler_SignIn_ReturnsAccessToken(t *testing.T) {
	userID := uuid.New()
	ts := newAuthServer(
		&stubSignUpper{},
		&stubSignIner{output: port.SignInOutput{
			User:        user.User{ID: userID, Email: "hero@apex20.com"},
			AccessToken: "jwt.token.here",
		}},
	)
	defer ts.Close()

	client := apex20v1connect.NewAuthServiceClient(ts.Client(), ts.URL)
	resp, err := client.SignIn(context.Background(), connect.NewRequest(&v1.SignInRequest{
		Email:    "hero@apex20.com",
		Password: "secret123",
	}))

	require.NoError(t, err)
	assert.Equal(t, userID.String(), resp.Msg.UserId)
	assert.Equal(t, "jwt.token.here", resp.Msg.AccessToken)
}

func TestAuthHandler_SignIn_ReturnsUnauthenticatedOnInvalidCredentials(t *testing.T) {
	ts := newAuthServer(
		&stubSignUpper{},
		&stubSignIner{err: port.ErrInvalidCredentials},
	)
	defer ts.Close()

	client := apex20v1connect.NewAuthServiceClient(ts.Client(), ts.URL)
	_, err := client.SignIn(context.Background(), connect.NewRequest(&v1.SignInRequest{
		Email:    "hero@apex20.com",
		Password: "wrong",
	}))

	require.Error(t, err)
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeUnauthenticated, connectErr.Code())
}
