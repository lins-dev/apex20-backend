-- name: ExistsAnyPermission :one
SELECT EXISTS (
    SELECT 1 FROM permissions WHERE deleted_at IS NULL
);

-- name: ListPermissions :many
SELECT * FROM permissions WHERE deleted_at IS NULL ORDER BY name;

-- name: GetPermissionByID :one
SELECT * FROM permissions WHERE id = $1 AND deleted_at IS NULL;

-- name: CreatePermission :exec
INSERT INTO permissions (id, name, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdatePermission :exec
UPDATE permissions SET name = $2, description = $3, updated_at = $4
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeletePermission :execrows
UPDATE permissions SET deleted_at = $2, updated_at = $2
WHERE id = $1 AND deleted_at IS NULL;
