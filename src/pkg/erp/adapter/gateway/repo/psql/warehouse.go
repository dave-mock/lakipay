package psql

import (
	"auth/src/pkg/erp/core/entity"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (repo PsqlRepo) CreateWarehouse(name, location string, capacity int, description string, userId uuid.UUID) (*entity.Warehouse, error) {
	const ErrFailedToStoreWarehouse = "FAILED_TO_STORE_WAREHOUSE"
	warehouseID := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt

	_, err := repo.db.Exec(`
		INSERT INTO erp.warehouses (id, merchant_id, name, location, capacity, description, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, warehouseID, userId, name, location, capacity, description, createdAt, updatedAt, userId, userId)

	if err != nil {

		// Log the detailed error for internal purposes
		repo.log.Println("CREATE WAREHOUSE ERROR:", err)
		return nil, fmt.Errorf("We encountered an issue while creating the warehouse. Please try again later.")
	}

	warehouse := &entity.Warehouse{
		Id:          warehouseID,
		MerchantID:  userId,
		Name:        name,
		Location:    location,
		Capacity:    capacity,
		Description: description,
		CreatedBy:   userId,
		UpdatedBy:   userId,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		IsActive:    true,
	}
	repo.log.Println("WAREHOUSE CREATED SUCCESSFULLY")
	return warehouse, nil
}

func (repo PsqlRepo) GetWarehouse(warehouseId, userId uuid.UUID) (*entity.Warehouse, error) {
	const ErrFailedToGetWarehouse = "FAILED_TO_GET_WAREHOUSE"
	var warehouse entity.Warehouse

	err := repo.db.QueryRow(`
		SELECT id, name, location, capacity, description, created_at, updated_at, created_by, updated_by, is_active
		FROM erp.warehouses
		WHERE id = $1 AND created_by = $2
	`, warehouseId, userId).Scan(
		&warehouse.Id, &warehouse.Name, &warehouse.Location, &warehouse.Capacity, &warehouse.Description,
		&warehouse.CreatedAt, &warehouse.UpdatedAt, &warehouse.CreatedBy, &warehouse.UpdatedBy, &warehouse.IsActive,
	)

	if err != nil {
		repo.log.Println("ERROR FETCHING WAREHOUSE:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToGetWarehouse, err)
	}

	repo.log.Println("WAREHOUSE FETCHED SUCCESSFULLY")
	return &warehouse, nil
}

func (repo PsqlRepo) UpdateWarehouse(id uuid.UUID, name, location string, capacity int, description string, userId uuid.UUID) (*entity.Warehouse, error) {
	const ErrFailedToUpdateWarehouse = "FAILED_TO_UPDATE_WAREHOUSE"
	updatedAt := time.Now()

	_, err := repo.db.Exec(`
		UPDATE erp.warehouses
		SET name = $1, location = $2, capacity = $3, description = $4, updated_at = $5, updated_by = $6
		WHERE id = $7
	`, name, location, capacity, description, updatedAt, userId, id)
	if err != nil {
		repo.log.Println("ERROR UPDATING WAREHOUSE:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToUpdateWarehouse, err)
	}

	// Return the updated warehouse
	warehouse := &entity.Warehouse{
		Id:          id,
		Name:        name,
		Location:    location,
		Capacity:    capacity,
		Description: description,
		UpdatedBy:   userId,
		UpdatedAt:   updatedAt,
		IsActive:    true,
	}
	repo.log.Println("WAREHOUSE UPDATED SUCCESSFULLY")
	return warehouse, nil
}

func (repo PsqlRepo) DeleteWarehouse(warehouseId, userId uuid.UUID) error {
	const ErrFailedToDeleteWarehouse = "FAILED_TO_DELETE_WAREHOUSE"
	_, err := repo.db.Exec(`
		DELETE FROM erp.warehouses WHERE id = $1 AND created_by = $2
	`, warehouseId, userId)
	if err != nil {
		repo.log.Println("ERROR DELETING WAREHOUSE:", err)
		return fmt.Errorf("%s: %v", ErrFailedToDeleteWarehouse, err)
	}

	repo.log.Println("WAREHOUSE DELETED SUCCESSFULLY")
	return nil
}

func (repo PsqlRepo) ListWarehouses(userId uuid.UUID) ([]entity.Warehouse, error) {
	const ErrFailedToListWarehouses = "FAILED_TO_LIST_WAREHOUSES"
	var warehouses []entity.Warehouse

	rows, err := repo.db.Query(`
		SELECT id, name, location, capacity, description, created_at, updated_at, created_by, updated_by, is_active
		FROM erp.warehouses
		WHERE created_by = $1
	`, userId)
	if err != nil {
		repo.log.Println("ERROR FETCHING WAREHOUSE LIST:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToListWarehouses, err)
	}
	defer rows.Close()

	for rows.Next() {
		var warehouse entity.Warehouse
		if err := rows.Scan(
			&warehouse.Id,
			&warehouse.Name,
			&warehouse.Location,
			&warehouse.Capacity,
			&warehouse.Description,
			&warehouse.CreatedAt,
			&warehouse.UpdatedAt,
			&warehouse.CreatedBy,
			&warehouse.UpdatedBy,
			&warehouse.IsActive,
		); err != nil {
			repo.log.Println("ERROR SCANNING WAREHOUSE:", err)
			continue
		}
		warehouses = append(warehouses, warehouse)
	}

	repo.log.Println("WAREHOUSE LIST FETCHED SUCCESSFULLY")
	return warehouses, nil
}

func (repo PsqlRepo) ListMerchantWarehouses(userId uuid.UUID) ([]entity.Warehouse, error) {
	const ErrFailedToListMerchantWarehouses = "FAILED_TO_LIST_MERCHANT_WAREHOUSES"
	var warehouses []entity.Warehouse

	rows, err := repo.db.Query(`
		SELECT id, name, location, capacity, description, created_at, updated_at, created_by, updated_by, is_active
		FROM erp.warehouses
		WHERE created_by = $1
	`, userId)
	if err != nil {
		repo.log.Println("ERROR FETCHING MERCHANT WAREHOUSE LIST:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToListMerchantWarehouses, err)
	}
	defer rows.Close()

	for rows.Next() {
		var warehouse entity.Warehouse
		if err := rows.Scan(
			&warehouse.Id,
			&warehouse.Name,
			&warehouse.Location,
			&warehouse.Capacity,
			&warehouse.Description,
			&warehouse.CreatedAt,
			&warehouse.UpdatedAt,
			&warehouse.CreatedBy,
			&warehouse.UpdatedBy,
			&warehouse.IsActive,
		); err != nil {
			repo.log.Println("ERROR SCANNING WAREHOUSE:", err)
			continue
		}
		warehouses = append(warehouses, warehouse)
	}

	repo.log.Println("MERCHANT WAREHOUSE LIST FETCHED SUCCESSFULLY")
	return warehouses, nil
}

func (repo PsqlRepo) DeactivateWarehouse(warehouseId, userId uuid.UUID) (*entity.Warehouse, error) {
	const ErrFailedToDeactivateWarehouse = "FAILED_TO_DEACTIVATE_WAREHOUSE"
	updatedAt := time.Now()
	_, err := repo.db.Exec(`
		UPDATE erp.warehouses
		SET is_active = false, updated_at = $1, updated_by = $2
		WHERE id = $3 AND created_by = $2
	`, updatedAt, userId, warehouseId)
	if err != nil {
		repo.log.Println("ERROR DEACTIVATING WAREHOUSE:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToDeactivateWarehouse, err)
	}

	warehouse, err := repo.GetWarehouse(
		warehouseId,
		userId)
	if err != nil {
		repo.log.Println("ERROR FETCHING DEACTIVATED WAREHOUSE:", err)
		return nil, err
	}

	repo.log.Println("WAREHOUSE DEACTIVATED SUCCESSFULLY")
	return warehouse, nil
}
