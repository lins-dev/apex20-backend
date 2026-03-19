-- 003_create_permissions.sql
CREATE TABLE permissions (
    id          UUID         PRIMARY KEY,
    name        VARCHAR(100) UNIQUE NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ  NULL
);

---- create above / drop below ----

DROP TABLE permissions;
