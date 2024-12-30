package repo

import (
	"auth/src/pkg/access_control/core/entity"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (repo PsqlRepo) CreateOperations(name, description string) (*entity.Operation, error) {
	const query = `
		INSERT INTO auth.Operations (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`
	var operations entity.Operation
	err := repo.db.QueryRow(query, name, description).Scan(
		&operations.ID, &operations.Name, &operations.Description, &operations.CreatedAt, &operations.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create operation: %v", err)
	}
	return &operations, nil
}

func (repo PsqlRepo) UpdateOperations(operationsID uuid.UUID, name, description string) (*entity.Operation, error) {

	_, err := repo.GetOperationsByID(operationsID)
	if err != nil {
		if err.Error() == "operation not found" {
			return nil, fmt.Errorf("operation with ID %v does not exist", operationsID)
		}
		return nil, fmt.Errorf("failed to retrieve operation: %v", err)
	}

	const query = `
		UPDATE auth.Operations
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`
	var operations entity.Operation
	err = repo.db.QueryRow(query, name, description, operationsID).Scan(
		&operations.ID, &operations.Name, &operations.Description, &operations.CreatedAt, &operations.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update operation: %v", err)
	}

	return &operations, nil
}

func (repo PsqlRepo) GetOperationsByID(operationsID uuid.UUID) (*entity.Operation, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM auth.Operations
		WHERE id = $1
	`
	var operations entity.Operation
	err := repo.db.QueryRow(query, operationsID).Scan(
		&operations.ID, &operations.Name, &operations.Description, &operations.CreatedAt, &operations.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("operation with ID %v not found", operationsID)
		}
		return nil, fmt.Errorf("failed to get operation with ID %v: %v", operationsID, err)
	}

	return &operations, nil
}

func (repo PsqlRepo) ListOperations() ([]*entity.Operation, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM auth.Operations
		ORDER BY created_at DESC
	`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list operations: %v", err)
	}
	defer rows.Close()

	var operations []*entity.Operation
	for rows.Next() {
		var operation entity.Operation
		err := rows.Scan(&operation.ID, &operation.Name, &operation.Description, &operation.CreatedAt, &operation.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan operation: %v", err)
		}
		operations = append(operations, &operation)
	}

	return operations, nil
}

func (repo PsqlRepo) DeleteOperations(operationsID uuid.UUID) error {
	_, err := repo.GetOperationsByID(operationsID)
	if err != nil {
		if err.Error() == "operation not found" {
			return fmt.Errorf("operation with ID %v does not exist", operationsID)
		}
		return fmt.Errorf("failed to get operation: %v", err)
	}

	const query = `
		DELETE FROM auth.Operations
		WHERE id = $1
	`
	result, err := repo.db.Exec(query, operationsID)
	if err != nil {
		return fmt.Errorf("failed to delete operation: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no operation was deleted")
	}

	return nil
}
