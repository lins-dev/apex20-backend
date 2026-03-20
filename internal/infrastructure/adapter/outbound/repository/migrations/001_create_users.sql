-- 001_create_users.sql
CREATE EXTENSION IF NOT EXISTS citext;

-- UUIDv7 gerado na camada de aplicação (PG16 não possui gen_random_uuid_v7 nativo)
-- is_admin controla acesso administrativo de plataforma (ADR-002)
-- roles de campanha são gerenciadas em campaign_members (005)
CREATE TABLE users (
    id            UUID         PRIMARY KEY,
    email         CITEXT       UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL,
    password_hash TEXT         NOT NULL,
    is_admin      BOOLEAN      NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ  NULL
);

---- create above / drop below ----

DROP TABLE users;
DROP EXTENSION citext;
