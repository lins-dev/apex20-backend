# apex20-backend

API HTTP/gRPC do **Apex20** — Virtual Tabletop para RPG.

Construído em Go com Chi, arquitetura hexagonal, sqlc e ConnectRPC. Persiste no PostgreSQL e publica eventos via Redis Pub/Sub.

## Pré-requisitos

- Go v1.26+
- PostgreSQL v16
- Redis v7

## Instalação

```bash
# Ativar git hooks locais
make setup

# Copiar e preencher variáveis de ambiente
cp .env.example .env
```

## Comandos

```bash
make run     # inicia o servidor
make build   # compila o binário em bin/api
make test    # executa os testes
make lint    # golangci-lint
```

## Estrutura

```
cmd/api/                              Entrypoint
internal/
  application/port/                   Interfaces (ports)
  infrastructure/adapter/inbound/     Handlers HTTP
  infrastructure/adapter/outbound/    Redis, repositórios SQL
  infrastructure/adapter/outbound/repository/migrations/  Migrações (tern)
```

## Variáveis de ambiente

Copie `.env.example` para `.env` e preencha os valores.

## Documentação

Consulte o submodule `docs/` ou o repositório [apex20-docs](https://github.com/lins-dev/apex20-docs).
