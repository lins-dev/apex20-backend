-- 004_create_role_permissions.sql
CREATE TABLE role_permissions (
    id              UUID         PRIMARY KEY,
    role            user_role    NOT NULL,
    permission_id   UUID         NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ  NULL,
    UNIQUE (role, permission_id)
);

---- create above / drop below ----

DROP TABLE role_permissions;
