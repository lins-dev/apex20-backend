.PHONY: run build test lint setup

run:
	go run ./cmd/api/...

build:
	go build -o bin/api ./cmd/api/...

test:
	go test ./...

lint:
	golangci-lint run ./...

setup:
	git config core.hooksPath .githooks
