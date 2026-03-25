package http_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"connectrpc.com/connect"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/apex20/contracts/proto/apex20/v1"
	"github.com/apex20/contracts/proto/apex20/v1/apex20v1connect"

	"github.com/apex20/backend/internal/application/usecase"
	adapter "github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
	"github.com/apex20/backend/internal/infrastructure/adapter/inbound/http/middleware"
	"github.com/apex20/backend/internal/infrastructure/adapter/outbound/crypto"
	jwtinfra "github.com/apex20/backend/internal/infrastructure/adapter/outbound/jwt"
	"github.com/apex20/backend/internal/infrastructure/adapter/outbound/repository"
)

const (
	testUserName     = "Integration User"
	testUserEmail    = "integration@apex20.dev"
	testUserPassword = "senha1234"
)

// --- helpers ---

func openAuthIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	dsn := os.Getenv("DATABASE_URL")
	require.NotEmpty(t, dsn, "DATABASE_URL is required for integration tests")
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	require.NoError(t, db.Ping(), "failed to connect to database")
	t.Cleanup(func() { db.Close() })
	return db
}

func cleanUsersTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.ExecContext(context.Background(),
		"DELETE FROM campaign_members; DELETE FROM campaigns; DELETE FROM users;",
	)
	require.NoError(t, err)
}

type authTestStack struct {
	server *httptest.Server
	client apex20v1connect.AuthServiceClient
}

func buildAuthStack(t *testing.T, db *sql.DB) authTestStack {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	hasher := crypto.NewArgon2PasswordHasher()
	tokenGen := jwtinfra.NewRSATokenGenerator(privateKey, time.Hour)
	tokenValidator := jwtinfra.NewRSATokenValidator(&privateKey.PublicKey)
	userRepo := repository.NewPostgresUserRepository(db)

	chiServer := adapter.NewChiServer()
	adapter.RegisterAuthHandler(chiServer, adapter.AuthUseCases{
		SignUp: usecase.NewSignUpUseCase(userRepo, hasher, tokenGen),
		SignIn: usecase.NewSignInUseCase(userRepo, hasher, tokenGen),
	})

	// Rota protegida simples para validar o middleware JWT
	chiServer.RegisterRoute("GET", "/protected", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	jwtMw := middleware.JWTAuth(tokenValidator,
		"/apex20.v1.AuthService/",
		"/health",
		"/docs",
		"/openapi",
	)

	ts := httptest.NewServer(jwtMw(chiServer.GetHandler()))
	t.Cleanup(ts.Close)

	return authTestStack{
		server: ts,
		client: apex20v1connect.NewAuthServiceClient(ts.Client(), ts.URL),
	}
}

func signUpTestUser(t *testing.T, client apex20v1connect.AuthServiceClient) string {
	t.Helper()
	resp, err := client.SignUp(context.Background(), connect.NewRequest(&v1.SignUpRequest{
		Name:     testUserName,
		Email:    testUserEmail,
		Password: testUserPassword,
	}))
	require.NoError(t, err)
	return resp.Msg.AccessToken
}

// --- SignUp ---

func TestAuthIntegration_SignUp_CreatesUserAndReturnsToken(t *testing.T) {
	db := openAuthIntegrationDB(t)
	cleanUsersTable(t, db)
	stack := buildAuthStack(t, db)

	resp, err := stack.client.SignUp(context.Background(), connect.NewRequest(&v1.SignUpRequest{
		Name:     testUserName,
		Email:    testUserEmail,
		Password: testUserPassword,
	}))

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Msg.UserId)
	assert.NotEmpty(t, resp.Msg.AccessToken)
}

func TestAuthIntegration_SignUp_DuplicateEmail_ReturnsAlreadyExists(t *testing.T) {
	db := openAuthIntegrationDB(t)
	cleanUsersTable(t, db)
	stack := buildAuthStack(t, db)

	signUpTestUser(t, stack.client)

	_, err := stack.client.SignUp(context.Background(), connect.NewRequest(&v1.SignUpRequest{
		Name:     testUserName,
		Email:    testUserEmail,
		Password: testUserPassword,
	}))

	require.Error(t, err)
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeAlreadyExists, connectErr.Code())
}

// --- SignIn ---

func TestAuthIntegration_SignIn_ValidCredentials_ReturnsToken(t *testing.T) {
	db := openAuthIntegrationDB(t)
	cleanUsersTable(t, db)
	stack := buildAuthStack(t, db)
	signUpTestUser(t, stack.client)

	resp, err := stack.client.SignIn(context.Background(), connect.NewRequest(&v1.SignInRequest{
		Email:    testUserEmail,
		Password: testUserPassword,
	}))

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Msg.AccessToken)
}

func TestAuthIntegration_SignIn_WrongPassword_ReturnsUnauthenticated(t *testing.T) {
	db := openAuthIntegrationDB(t)
	cleanUsersTable(t, db)
	stack := buildAuthStack(t, db)
	signUpTestUser(t, stack.client)

	_, err := stack.client.SignIn(context.Background(), connect.NewRequest(&v1.SignInRequest{
		Email:    testUserEmail,
		Password: "senhaerrada",
	}))

	require.Error(t, err)
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeUnauthenticated, connectErr.Code())
}

func TestAuthIntegration_SignIn_UnknownEmail_ReturnsUnauthenticated(t *testing.T) {
	db := openAuthIntegrationDB(t)
	cleanUsersTable(t, db)
	stack := buildAuthStack(t, db)

	_, err := stack.client.SignIn(context.Background(), connect.NewRequest(&v1.SignInRequest{
		Email:    "ghost@apex20.dev",
		Password: testUserPassword,
	}))

	require.Error(t, err)
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeUnauthenticated, connectErr.Code())
}

// --- JWT Middleware ---

func TestAuthIntegration_ProtectedRoute_WithValidToken_Returns200(t *testing.T) {
	db := openAuthIntegrationDB(t)
	cleanUsersTable(t, db)
	stack := buildAuthStack(t, db)
	token := signUpTestUser(t, stack.client)

	req, err := http.NewRequest(http.MethodGet, stack.server.URL+"/protected", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := stack.server.Client().Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthIntegration_ProtectedRoute_WithoutToken_Returns401(t *testing.T) {
	db := openAuthIntegrationDB(t)
	cleanUsersTable(t, db)
	stack := buildAuthStack(t, db)

	req, err := http.NewRequest(http.MethodGet, stack.server.URL+"/protected", nil)
	require.NoError(t, err)

	resp, err := stack.server.Client().Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthIntegration_ProtectedRoute_WithInvalidToken_Returns401(t *testing.T) {
	db := openAuthIntegrationDB(t)
	cleanUsersTable(t, db)
	stack := buildAuthStack(t, db)

	req, err := http.NewRequest(http.MethodGet, stack.server.URL+"/protected", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")

	resp, err := stack.server.Client().Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
