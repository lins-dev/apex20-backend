package main

import (
	"context"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"database/sql"

	"github.com/apex20/backend/internal/application/usecase"
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

	repo := repository.NewPostgresPermissionRepository(db)
	uc := usecase.NewSeedPermissionsUseCase(repo)

	if err := uc.Execute(context.Background()); err != nil {
		log.Fatalf("seeding permissions: %v", err)
	}

	log.Println("permissions seeded successfully")
}
