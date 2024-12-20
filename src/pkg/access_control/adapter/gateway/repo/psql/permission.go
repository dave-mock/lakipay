package repo

import (
	"auth/src/pkg/access_control/core/entity"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (repo PsqlRepo) CreatePermission(resourceID uuid.UUID, operation, effect string) (*entity.Permission, error) {
	const query = `
		INSERT INTO auth.permissions (resource_id, operation, effect)
		VALUES ($1, $2, $3)
		RETURNING id, resource_id, operation, effect, created_at, updated_at
	`
	var permission entity.Permission
	err := repo.db.QueryRow(query, resourceID, operation, effect).Scan(
		&permission.ID,
		&permission.ResourceID,
		&permission.Operation,
		&permission.Effect,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %v", err)
	}

	const fetchResourceNameQuery = `
		SELECT r.name
		FROM auth.resources r
		WHERE r.id = $1
	`
	var resourceName string
	err = repo.db.QueryRow(fetchResourceNameQuery, resourceID).Scan(&resourceName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("resource with ID %v not found", resourceID)
		}
		return nil, fmt.Errorf("failed to fetch resource name: %v", err)
	}

	permission.Resource = resourceName
	return &permission, nil
}

func (repo PsqlRepo) UpdatePermission(permissionID uuid.UUID, operation, effect string) (*entity.Permission, error) {
	const query = `
        UPDATE auth.permissions
        SET operation = $2, effect = $3, updated_at = NOW()
        WHERE id = $1
        RETURNING id, resource_id, operation, effect, created_at, updated_at
    `

	var permission entity.Permission
	err := repo.db.QueryRow(query, permissionID, operation, effect).Scan(
		&permission.ID,
		&permission.ResourceID,
		&permission.Operation,
		&permission.Effect,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update permission: %v", err)
	}

	return &permission, nil
}

func (repo PsqlRepo) ListPermissions() ([]entity.Permission, error) {
	const query = `
		SELECT id, resource_id, operation, effect, created_at, updated_at
		FROM auth.permissions
	`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %v", err)
	}
	defer rows.Close()

	var permissions []entity.Permission
	for rows.Next() {
		var permission entity.Permission
		if err := rows.Scan(&permission.ID, &permission.ResourceID, &permission.Operation, &permission.Effect, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %v", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over permissions: %v", err)
	}

	return permissions, nil
}

func (repo PsqlRepo) DeletePermission(resourceID, permissionID uuid.UUID) error {
	const resourceExistsQuery = `
		SELECT 1 FROM auth.resources WHERE id = $1
	`
	var resourceExists bool
	err := repo.db.QueryRow(resourceExistsQuery, resourceID).Scan(&resourceExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("resource with ID %v does not exist", resourceID)
		}
		return fmt.Errorf("failed to check resource existence: %v", err)
	}

	const permissionExistsQuery = `
		SELECT 1 FROM auth.permissions WHERE permission_id = $1
	`
	var permissionExists bool
	err = repo.db.QueryRow(permissionExistsQuery, permissionID).Scan(&permissionExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("permission with ID %v does not exist", permissionID)
		}
		return fmt.Errorf("failed to check permission existence: %v", err)
	}

	const deleteQuery = `
		DELETE FROM auth.permissions
		WHERE resource_id = $1 AND permission_id = $2
	`
	_, err = repo.db.Exec(deleteQuery, resourceID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %v", err)
	}

	return nil
}

func (repo PsqlRepo) ListUserPermissions(userID uuid.UUID) ([]entity.Permission, error) {
	const query = `
        SELECT 
            p.id, 
            r.name AS resource, 
            p.resource_id AS resource_identifier, 
            p.operation, 
            p.effect, 
            p.created_at, 
            p.updated_at
        FROM auth.permissions p
        JOIN auth.user_permissions up ON p.id = up.permission_id
        JOIN auth.resources r ON p.resource_id = r.id
        WHERE up.user_id = $1
    `
	rows, err := repo.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user permissions: %v", err)
	}
	defer rows.Close()

	var permissions []entity.Permission
	for rows.Next() {
		var permission entity.Permission
		if err := rows.Scan(
			&permission.ID,
			&permission.Resource,
			&permission.ResourceID,
			&permission.Operation,
			&permission.Effect,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %v", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over user permissions: %v", err)
	}

	if len(permissions) == 0 {
		return nil, fmt.Errorf("no permissions found for the provided user ID")
	}

	return permissions, nil
}

func (repo PsqlRepo) ListGroupPermissions(groupID uuid.UUID) ([]entity.Permission, error) {
	const query = `
		SELECT 
			p.id, 
			r.name AS resource, 
			p.resource_id AS resource_identifier, 
			p.operation, 
			p.effect, 
			p.created_at, 
			p.updated_at
		FROM auth.permissions p
		JOIN auth.group_permissions gp ON p.id = gp.permission_id
		JOIN auth.resources r ON p.resource_id = r.id
		WHERE gp.group_id = $1
	`
	rows, err := repo.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list group permissions: %v", err)
	}
	defer rows.Close()

	var permissions []entity.Permission
	for rows.Next() {
		var permission entity.Permission
		if err := rows.Scan(
			&permission.ID,
			&permission.Resource,
			&permission.ResourceID,
			&permission.Operation,
			&permission.Effect,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %v", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over group permissions: %v", err)
	}

	if len(permissions) == 0 {
		return nil, fmt.Errorf("no permissions found for the provided group ID")
	}

	return permissions, nil
}
