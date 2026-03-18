-- 001_create_users.sql
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE user_role AS ENUM ('gm', 'player', 'trusted');

-- UUIDv7 gerado na camada de aplicação (PG16 não possui gen_random_uuid_v7 nativo)
CREATE TABLE users (
    id            UUID         PRIMARY KEY,
    email         CITEXT       UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL,
    password_hash TEXT         NOT NULL,
    role          user_role    NOT NULL DEFAULT 'player',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ  NULL
);

---- create above / drop below ----

DROP TABLE users;
DROP TYPE user_role;
DROP EXTENSION citext;
