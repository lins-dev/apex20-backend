.PHONY: run build test test-integration migrate-test lint setup

MIGRATIONS_DIR=./internal/infrastructure/adapter/outbound/repository/migrations

run:
	go run ./cmd/api/...

build:
	go build -o bin/api ./cmd/api/...

test:
	go test -short ./...

migrate-test:
	tern migrate --migrations $(MIGRATIONS_DIR) --conn-string "$(DATABASE_URL)"

test-integration: migrate-test
	go test ./...

lint:
	golangci-lint run ./...

setup:
	git config core.hooksPath .githooks
