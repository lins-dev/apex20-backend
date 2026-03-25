package main

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"log"
	nethttp "net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/cors"

	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
	"github.com/apex20/backend/internal/infrastructure/adapter/inbound/http/middleware"
	"github.com/apex20/backend/internal/infrastructure/adapter/outbound/crypto"
	jwtinfra "github.com/apex20/backend/internal/infrastructure/adapter/outbound/jwt"
	"github.com/apex20/backend/internal/infrastructure/adapter/outbound/repository"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("connecting to database: %v", err)
	}

	privateKey := mustLoadRSAPrivateKey(os.Getenv("JWT_PRIVATE_KEY_PEM"))
	publicKey := mustLoadRSAPublicKey(os.Getenv("JWT_PUBLIC_KEY_PEM"))

	hasher := crypto.NewArgon2PasswordHasher()
	tokenGen := jwtinfra.NewRSATokenGenerator(privateKey, 24*time.Hour)
	tokenValidator := jwtinfra.NewRSATokenValidator(publicKey)
	userRepo := repository.NewPostgresUserRepository(db)

	authUC := http.AuthUseCases{
		SignUp: usecase.NewSignUpUseCase(userRepo, hasher, tokenGen),
		SignIn: usecase.NewSignInUseCase(userRepo, hasher, tokenGen),
	}

	permRepo := repository.NewPostgresPermissionRepository(db)
	rolePermRepo := repository.NewPostgresRolePermissionRepository(db)
	campaignRepo := repository.NewPostgresCampaignRepository(db)

	permUC := http.PermissionUseCases{
		List:   usecase.NewListPermissionsUseCase(permRepo),
		Get:    usecase.NewGetPermissionUseCase(permRepo),
		Create: usecase.NewCreatePermissionUseCase(permRepo),
		Update: usecase.NewUpdatePermissionUseCase(permRepo),
		Delete: usecase.NewDeletePermissionUseCase(permRepo),
	}

	rolePermUC := http.RolePermissionUseCases{
		List:   usecase.NewListRolePermissionsUseCase(rolePermRepo),
		Get:    usecase.NewGetRolePermissionUseCase(rolePermRepo),
		Create: usecase.NewCreateRolePermissionUseCase(rolePermRepo),
		Delete: usecase.NewDeleteRolePermissionUseCase(rolePermRepo),
	}

	server := http.NewChiServer()
	http.RegisterPermissionHandler(server.GetAPI(), permUC)
	http.RegisterRolePermissionHandler(server.GetAPI(), rolePermUC)
	roleUC := http.RoleUseCases{
		List: usecase.NewListRolesUseCase(),
	}
	http.RegisterRoleHandler(server.GetAPI(), roleUC)

	campaignMemberRepo := repository.NewPostgresCampaignMemberRepository(db)

	campaignUC := http.CampaignUseCases{
		Create: usecase.NewCreateCampaignUseCase(campaignRepo),
		List:   usecase.NewListCampaignsUseCase(campaignRepo),
		Get:    usecase.NewGetCampaignUseCase(campaignRepo),
		Update: usecase.NewUpdateCampaignUseCase(campaignRepo),
		Delete: usecase.NewDeleteCampaignUseCase(campaignRepo),
	}
	http.RegisterCampaignHandler(server.GetAPI(), campaignUC)

	campaignMemberUC := http.CampaignMemberUseCases{
		Invite: usecase.NewInviteMemberUseCase(campaignMemberRepo),
		Remove: usecase.NewRemoveMemberUseCase(campaignMemberRepo),
	}
	http.RegisterCampaignMemberHandler(server.GetAPI(), campaignMemberUC)
	http.RegisterAuthHandler(server, authUC)

	userUC := http.UserUseCases{
		Get:    usecase.NewGetUserUseCase(userRepo),
		Update: usecase.NewUpdateUserUseCase(userRepo),
		Delete: usecase.NewDeleteUserUseCase(userRepo),
	}
	http.RegisterUserHandler(server.GetAPI(), userUC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtMw := middleware.JWTAuth(tokenValidator,
		"/health",
		"/apex20.v1.AuthService/",
		"/docs",
		"/openapi",
	)

	// CORS — permite requisições do frontend (browser) para o backend.
	// AllowedHeaders inclui os headers exigidos pelo protocolo ConnectRPC.
	corsMw := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
			"Connect-Protocol-Version",
			"Connect-Timeout-Ms",
			"X-User-Agent",
		},
		ExposedHeaders: []string{"Content-Type"},
		MaxAge:         7200,
	})

	log.Printf("Backend server starting on :%s...", port)
	log.Printf("OpenAPI documentation available at :%s/docs", port)
	if err := nethttp.ListenAndServe(":"+port, corsMw.Handler(jwtMw(server.GetHandler()))); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func mustLoadRSAPrivateKey(pemStr string) *rsa.PrivateKey {
	if pemStr == "" {
		log.Fatal("JWT_PRIVATE_KEY_PEM is required")
	}
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		log.Fatal("JWT_PRIVATE_KEY_PEM: invalid PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("JWT_PRIVATE_KEY_PEM: %v", err)
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		log.Fatal("JWT_PRIVATE_KEY_PEM: not an RSA key")
	}
	return rsaKey
}

func mustLoadRSAPublicKey(pemStr string) *rsa.PublicKey {
	if pemStr == "" {
		log.Fatal("JWT_PUBLIC_KEY_PEM is required")
	}
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		log.Fatal("JWT_PUBLIC_KEY_PEM: invalid PEM block")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("JWT_PUBLIC_KEY_PEM: %v", err)
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		log.Fatal("JWT_PUBLIC_KEY_PEM: not an RSA key")
	}
	return rsaKey
}
