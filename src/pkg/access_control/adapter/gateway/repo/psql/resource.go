package repo

import (
	"auth/src/pkg/access_control/core/entity"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (repo PsqlRepo) CreateResource(name, description string, operations []uuid.UUID) (*entity.Resource, error) {
	const query = `
		INSERT INTO auth.resources (name, description, operations)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, operations, created_at, updated_at
	`
	var resource entity.Resource
	err := repo.db.QueryRow(query, name, description, pq.Array(operations)).Scan(
		&resource.ID,
		&resource.Name,
		&resource.Description,
		&resource.Operations,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %v", err)
	}
	return &resource, nil
}

func (repo PsqlRepo) UpdateResource(resourceID uuid.UUID, name, description string, operations []uuid.UUID) (*entity.Resource, error) {
	const resourceExistsQuery = `
		SELECT 1 FROM auth.resources WHERE id = $1
	`
	var resourceExists bool
	err := repo.db.QueryRow(resourceExistsQuery, resourceID).Scan(&resourceExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("resource with ID %v does not exist", resourceID)
		}
		return nil, fmt.Errorf("failed to check resource existence %v", err)
	}

	const updateQuery = `
		UPDATE auth.resources
		SET name = $2, description = $3, operations = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, operations, created_at, updated_at
	`

	var resource entity.Resource
	err = repo.db.QueryRow(updateQuery, resourceID, name, description, pq.Array(operations)).Scan(
		&resource.ID,
		&resource.Name,
		&resource.Description,
		pq.Array(&resource.Operations),
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update resource: %v", err)
	}

	return &resource, nil
}

func (repo PsqlRepo) GetResourceByID(resourceID uuid.UUID) (*entity.Resource, error) {
	const query = `
		SELECT id, name, description, operations, created_at, updated_at
		FROM auth.resources
		WHERE id = $1
	`
	var resource entity.Resource
	err := repo.db.QueryRow(query, resourceID).Scan(
		&resource.ID, &resource.Name, &resource.Description, pq.Array(&resource.Operations), &resource.CreatedAt, &resource.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("resource with ID %v not found", resourceID)
		}
		return nil, fmt.Errorf("failed to retrieve resource: %v", err)
	}
	return &resource, nil
}

func (repo PsqlRepo) ListResources() ([]entity.Resource, error) {
	const query = `
		SELECT id, name, description, operations, created_at, updated_at
		FROM auth.resources
		ORDER BY created_at DESC
	`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %v", err)
	}
	defer rows.Close()

	var resources []entity.Resource
	for rows.Next() {
		var resource entity.Resource
		err := rows.Scan(&resource.ID, &resource.Name, &resource.Description, pq.Array(&resource.Operations), &resource.CreatedAt, &resource.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan resource: %v", err)
		}
		resources = append(resources, resource)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading rows: %v", err)
	}

	return resources, nil
}

func (repo PsqlRepo) DeleteResource(resourceID uuid.UUID) error {
	const query = `
		DELETE FROM auth.resources
		WHERE id = $1
	`
	result, err := repo.db.Exec(query, resourceID)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no resource was deleted")
	}
	return nil
}
