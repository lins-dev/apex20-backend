package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

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

	server := http.NewChiServer()
	http.RegisterPermissionHandler(server.GetAPI(), permRepo)
	http.RegisterRolePermissionHandler(server.GetAPI(), rolePermRepo)
	http.RegisterRoleHandler(server.GetAPI())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := server.Start(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
