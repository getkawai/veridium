-- Roles

-- name: GetRole :one
SELECT * FROM rbac_roles WHERE id = ?;

-- name: GetRoleByName :one
SELECT * FROM rbac_roles WHERE name = ?;

-- name: ListRoles :many
SELECT * FROM rbac_roles
WHERE is_active = 1
ORDER BY name ASC;

-- name: CreateRole :one
INSERT INTO rbac_roles (
    name, display_name, description, is_system, is_active, metadata,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateRole :one
UPDATE rbac_roles
SET display_name = ?,
    description = ?,
    is_active = ?,
    metadata = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM rbac_roles WHERE id = ? AND is_system = 0;

-- Permissions

-- name: GetPermission :one
SELECT * FROM rbac_permissions WHERE id = ?;

-- name: GetPermissionByCode :one
SELECT * FROM rbac_permissions WHERE code = ?;

-- name: ListPermissions :many
SELECT * FROM rbac_permissions
WHERE is_active = 1
ORDER BY category ASC, name ASC;

-- name: ListPermissionsByCategory :many
SELECT * FROM rbac_permissions
WHERE category = ? AND is_active = 1
ORDER BY name ASC;

-- name: CreatePermission :one
INSERT INTO rbac_permissions (
    code, name, description, category, is_active, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdatePermission :one
UPDATE rbac_permissions
SET name = ?,
    description = ?,
    is_active = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeletePermission :exec
DELETE FROM rbac_permissions WHERE id = ?;

-- Role Permissions

-- name: LinkRoleToPermission :exec
INSERT INTO rbac_role_permissions (role_id, permission_id, created_at)
VALUES (?, ?, ?);

-- name: UnlinkRoleFromPermission :exec
DELETE FROM rbac_role_permissions
WHERE role_id = ? AND permission_id = ?;

-- name: GetRolePermissions :many
SELECT p.* FROM rbac_permissions p
INNER JOIN rbac_role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = ?;

-- User Roles

-- name: LinkUserToRole :exec
INSERT INTO rbac_user_roles (user_id, role_id, created_at, expires_at)
VALUES (?, ?, ?, ?);

-- name: UnlinkUserFromRole :exec
DELETE FROM rbac_user_roles
WHERE user_id = ? AND role_id = ?;

-- name: GetUserRoles :many
SELECT r.* FROM rbac_roles r
INNER JOIN rbac_user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = ? AND (ur.expires_at IS NULL OR ur.expires_at > ?);

-- name: GetUserPermissions :many
SELECT DISTINCT p.* FROM rbac_permissions p
INNER JOIN rbac_role_permissions rp ON p.id = rp.permission_id
INNER JOIN rbac_user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = ? AND (ur.expires_at IS NULL OR ur.expires_at > ?);

