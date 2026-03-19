-- 001_create_users.sql
CREATE EXTENSION IF NOT EXISTS citext;

-- UUIDv7 gerado na camada de aplicação (PG16 não possui gen_random_uuid_v7 nativo)
-- role é SMALLINT mapeado via enum no Protobuf (apex20-contracts/proto/apex20/v1/role.proto)
-- Valores: 1=GM, 2=Player, 3=Trusted, 4=Admin
CREATE TABLE users (
    id            UUID         PRIMARY KEY,
    email         CITEXT       UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL,
    password_hash TEXT         NOT NULL,
    role          SMALLINT     NOT NULL DEFAULT 2,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ  NULL
);

---- create above / drop below ----

DROP TABLE users;
DROP EXTENSION citext;
