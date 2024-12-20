package repo

import (
	"auth/src/pkg/access_control/core/entity"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (repo PsqlRepo) CreateResource(name, description string) (*entity.Resource, error) {
	const query = `
		INSERT INTO auth.resources (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`
	var resource entity.Resource
	err := repo.db.QueryRow(query, name, description).Scan(
		&resource.ID, &resource.Name, &resource.Description, &resource.CreatedAt, &resource.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %v", err)
	}
	return &resource, nil
}

func (repo PsqlRepo) UpdateResource(resourceID uuid.UUID, name, description string) (*entity.Resource, error) {

	_, err := repo.GetResourceByID(resourceID)
	if err != nil {
		if err.Error() == "resource not found" {
			return nil, fmt.Errorf("resource with ID %v does not exist", resourceID)
		}
		return nil, fmt.Errorf("failed to retrieve resource: %v", err)
	}

	const query = `
		UPDATE auth.resources
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`
	var resource entity.Resource
	err = repo.db.QueryRow(query, name, description, resourceID).Scan(
		&resource.ID, &resource.Name, &resource.Description, &resource.CreatedAt, &resource.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update resource: %v", err)
	}

	return &resource, nil
}

func (repo PsqlRepo) GetResourceByID(resourceID uuid.UUID) (*entity.Resource, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM auth.resources
		WHERE id = $1
	`
	var resource entity.Resource
	err := repo.db.QueryRow(query, resourceID).Scan(
		&resource.ID, &resource.Name, &resource.Description, &resource.CreatedAt, &resource.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("resource with ID %v not found", resourceID)
		}
		return nil, fmt.Errorf("failed to get resource with ID %v: %v", resourceID, err)
	}

	return &resource, nil
}

func (repo PsqlRepo) ListResources() ([]entity.Resource, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
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
		err := rows.Scan(&resource.ID, &resource.Name, &resource.Description, &resource.CreatedAt, &resource.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan resource: %v", err)
		}
		resources = append(resources, resource)
	}

	return resources, nil
}
func (repo PsqlRepo) DeleteResource(resourceID uuid.UUID) error {

	_, err := repo.GetResourceByID(resourceID)
	if err != nil {
		if err.Error() == "resource not found" {
			return fmt.Errorf("resource with ID %v does not exist", resourceID)
		}
		return fmt.Errorf("Resource %v", err)
	}

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
