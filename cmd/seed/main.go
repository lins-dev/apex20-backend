package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/apex20/backend/internal/application/usecase"
	"github.com/apex20/backend/internal/infrastructure/adapter/outbound/crypto"
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

	userRepo := repository.NewPostgresUserRepository(db)
	permRepo := repository.NewPostgresPermissionRepository(db)
	rolePermRepo := repository.NewPostgresRolePermissionRepository(db)

	seedPerms := usecase.NewSeedPermissionsUseCase(permRepo)
	seedRolePerms := usecase.NewSeedRolePermissionsUseCase(rolePermRepo)
	seedAdmin := usecase.NewSeedAdminUserUseCase(userRepo, userRepo, crypto.NewArgon2PasswordHasher())

	permissionIDs, err := seedPerms.Execute(context.Background())
	if err != nil {
		log.Fatalf("seeding permissions: %v", err)
	}

	if err := seedRolePerms.Execute(context.Background(), permissionIDs); err != nil {
		log.Fatalf("seeding role permissions: %v", err)
	}

	if err := seedAdmin.Execute(context.Background()); err != nil {
		log.Fatalf("seeding admin user: %v", err)
	}

	log.Println("seed completed successfully")
}
