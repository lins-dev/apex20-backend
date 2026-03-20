-- name: ListRolePermissions :many
SELECT * FROM role_permissions WHERE deleted_at IS NULL ORDER BY role;

-- name: GetRolePermissionByID :one
SELECT * FROM role_permissions WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateRolePermission :exec
INSERT INTO role_permissions (id, role, permission_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: SoftDeleteRolePermission :execrows
UPDATE role_permissions SET deleted_at = $2, updated_at = $2
WHERE id = $1 AND deleted_at IS NULL;
