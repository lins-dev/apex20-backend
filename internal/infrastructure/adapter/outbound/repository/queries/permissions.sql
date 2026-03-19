-- name: ExistsAnyPermission :one
SELECT EXISTS (
    SELECT 1 FROM permissions WHERE deleted_at IS NULL
);

-- name: CreatePermission :exec
INSERT INTO permissions (id, name, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: CreateRolePermission :exec
INSERT INTO role_permissions (id, role, permission_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);
