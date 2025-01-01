package repo

import (
	"auth/src/pkg/access_control/core/entity"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (repo PsqlRepo) CreatePermission(resourceID string, resource uuid.UUID, operations []uuid.UUID, effect string) (*entity.Permission, error) {
	if operations == nil {
		operations = []uuid.UUID{}
	}

	const checkQuery = `
		SELECT id
		FROM auth.permissions
		WHERE resource_id = $1 AND operations @> $2 AND effect = $3
		LIMIT 1
	`
	var existingID uuid.UUID
	err := repo.db.QueryRow(checkQuery, resourceID, pq.Array(operations), effect).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check for duplicate permission %v", err)
	}

	if existingID != uuid.Nil {
		return nil, fmt.Errorf("permission already exists for resource_id %v with the same operations", resourceID)
	}

	const insertQuery = `
		INSERT INTO auth.permissions (resource_id, resource, operations, effect, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, resource_id, resource, operations, effect, created_at, updated_at
	`
	var permission entity.Permission
	err = repo.db.QueryRow(insertQuery, resourceID, resource, pq.Array(operations), effect).Scan(
		&permission.ID,
		&permission.ResourceID,
		&permission.Resource,
		pq.Array(&permission.Operations),
		&permission.Effect,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission %v", err)
	}

	const fetchResourceNameQuery = `
		SELECT name
		FROM auth.resources
		WHERE id = $1
	`
	var resourceName string
	err = repo.db.QueryRow(fetchResourceNameQuery, permission.Resource).Scan(&resourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resource name %v", err)
	}

	permission.ResourceName = resourceName
	operationNames := make([]string, len(permission.Operations))
	for i, operationID := range permission.Operations {
		const fetchOperationNameQuery = `
			SELECT name
			FROM auth.operations
			WHERE id = $1
		`
		err := repo.db.QueryRow(fetchOperationNameQuery, operationID).Scan(&operationNames[i])
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operation name for operation_id %v %v", operationID, err)
		}
	}

	permission.OperationName = strings.Join(operationNames, ", ")

	return &permission, nil
}

func (repo PsqlRepo) UpdatePermission(permissionID uuid.UUID, resourceID string, resource uuid.UUID, operations []uuid.UUID, effect string) (*entity.Permission, error) {
	if operations == nil {
		operations = []uuid.UUID{}
	}

	const resourceExistsQuery = `
        SELECT EXISTS(SELECT 1 FROM auth.resources WHERE id = $1)
    `
	var resourceExists bool
	err := repo.db.QueryRow(resourceExistsQuery, resource).Scan(&resourceExists)
	if err != nil {
		return nil, fmt.Errorf("error checking resource existence: %v", err)
	}
	if !resourceExists {
		return nil, fmt.Errorf("resource with id %v does not exist", resource)
	}

	const updatePermissionQuery = `
        UPDATE auth.permissions
        SET resource = $2, resource_id = $3, operations = $4, effect = $5, updated_at = NOW()
        WHERE id = $1
        RETURNING id, resource_id, resource, operations, effect, created_at, updated_at
    `
	var permission entity.Permission
	err = repo.db.QueryRow(updatePermissionQuery, permissionID, resource, resourceID, pq.Array(operations), effect).Scan(
		&permission.ID,
		&permission.ResourceID,
		&permission.Resource,
		pq.Array(&permission.Operations),
		&permission.Effect,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no permission found with id %v", permissionID)
		}
		return nil, fmt.Errorf("error updating permission: %v", err)
	}

	const fetchResourceNameQuery = `
		SELECT name
		FROM auth.resources
		WHERE id = $1
	`
	var resourceName string
	err = repo.db.QueryRow(fetchResourceNameQuery, permission.Resource).Scan(&resourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resource name: %v", err)
	}
	permission.ResourceName = resourceName
	operationNames := make([]string, len(permission.Operations))
	for i, operationID := range permission.Operations {
		const fetchOperationNameQuery = `
			SELECT name
			FROM auth.operations
			WHERE id = $1
		`
		err := repo.db.QueryRow(fetchOperationNameQuery, operationID).Scan(&operationNames[i])
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operation name for operation_id %v: %v", operationID, err)
		}
	}

	permission.OperationName = strings.Join(operationNames, ", ")

	return &permission, nil
}

func (repo PsqlRepo) ListPermissions() ([]entity.Permission, error) {
	const query = `
		SELECT 
			p.id,
			r.name AS resource,
			p.resource_id,
			p.operations,
			string_agg(o.name, ', ') AS operation_names,  -- Aggregate operation names into a comma-separated string
			p.effect,
			p.created_at,
			p.updated_at
		FROM auth.permissions p
		JOIN auth.resources r ON p.resource_id = r.id
		LEFT JOIN auth.operations o ON o.id = ANY(p.operations)  -- Join on operations using ANY operator
		GROUP BY p.id, r.name, p.resource_id, p.effect, p.created_at, p.updated_at
	`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions %v", err)
	}
	defer rows.Close()

	var permissions []entity.Permission
	for rows.Next() {
		var permission entity.Permission
		if err := rows.Scan(
			&permission.ID,
			&permission.Resource,
			&permission.ResourceID,
			pq.Array(&permission.Operations),
			&permission.OperationName,
			&permission.Effect,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan permission %v", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over permissions %v", err)
	}

	return permissions, nil
}

func (repo PsqlRepo) DeletePermission(ResourceID uuid.UUID, permissionID uuid.UUID) error {
	const checkQuery = `
		SELECT COUNT(*) 
		FROM auth.permissions
		WHERE resource_id = $1 AND id = $2
	`
	var count int
	err := repo.db.QueryRow(checkQuery, ResourceID, permissionID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check permission existence %v", err)
	}
	if count == 0 {
		return fmt.Errorf("no permission found for resource_id %v and permission_id %v", ResourceID, permissionID)
	}

	const deleteQuery = `
		DELETE FROM auth.permissions
		WHERE resource_id = $1 AND id = $2
	`
	_, err = repo.db.Exec(deleteQuery, ResourceID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to delete permission %v", err)
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
		return nil, fmt.Errorf("failed to list user permissions %v", err)
	}
	defer rows.Close()

	var permissions []entity.Permission
	for rows.Next() {
		var permission entity.Permission
		if err := rows.Scan(
			&permission.ID,
			&permission.Resource,
			&permission.ResourceID,
			&permission.Operations,
			&permission.Effect,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan permission %v", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over user permissions %v", err)
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
		return nil, fmt.Errorf("failed to list group permissions %v", err)
	}
	defer rows.Close()

	var permissions []entity.Permission
	for rows.Next() {
		var permission entity.Permission
		if err := rows.Scan(
			&permission.ID,
			&permission.Resource,
			&permission.ResourceID,
			&permission.Operations,
			&permission.Effect,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan permission %v", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over group permissions %v", err)
	}

	if len(permissions) == 0 {
		return nil, fmt.Errorf("no permissions found for the provided group ID")
	}

	return permissions, nil
}
