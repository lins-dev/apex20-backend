package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/infrastructure/adapter/inbound/http"
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

	campaignUC := http.CampaignUseCases{
		Create: usecase.NewCreateCampaignUseCase(campaignRepo),
		List:   usecase.NewListCampaignsUseCase(campaignRepo),
		Get:    usecase.NewGetCampaignUseCase(campaignRepo),
		Update: usecase.NewUpdateCampaignUseCase(campaignRepo),
		Delete: usecase.NewDeleteCampaignUseCase(campaignRepo),
	}
	http.RegisterCampaignHandler(server.GetAPI(), campaignUC)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := server.Start(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
