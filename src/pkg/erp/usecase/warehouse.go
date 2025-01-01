package usecase

import (
	"auth/src/pkg/erp/core/entity"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	ErrFailedToCreateWarehouse     = "FAILED_TO_CREATE_WAREHOUSE"
	ErrFailedToUpdateWarehouse     = "FAILED_TO_UPDATE_WAREHOUSE"
	ErrFailedToDeleteWarehouse     = "FAILED_TO_DELETE_WAREHOUSE"
	ErrFailedToListWarehouses      = "FAILED_TO_LIST_WAREHOUSES"
	ErrFailedToGetWarehouse        = "FAILED_TO_GET_WAREHOUSE"
	ErrFailedToDeactivateWarehouse = "FAILED_TO_DEACTIVATE_WAREHOUSE"
)

// creates a new warehouse
func (uc Usecase) CreateWarehouse(name, location string, capacity int, description string, userId uuid.UUID) (*entity.Warehouse, error) {
	// Validate input
	if err := validateWarehouseInput(name, location, capacity, description); err != nil {
		uc.log.Println("ERROR VALIDATING WAREHOUSE INPUT:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToCreateWarehouse, err)
	}

	// Create warehouse entity without an ID and IsActive
	warehouse := &entity.Warehouse{
		Name:      name,
		Location:  location,
		Capacity:  capacity,
		CreatedBy: userId,
		UpdatedBy: userId,
		IsActive:  true,
	}

	// Log action
	uc.log.Println("CREATING WAREHOUSE")

	// Call repository method to create the warehouse
	warehouse, err := uc.repo.CreateWarehouse(
		warehouse.Name, warehouse.Location, warehouse.Capacity, description, warehouse.CreatedBy,
	)
	if err != nil {
		uc.log.Println("ERROR CREATING WAREHOUSE:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToCreateWarehouse, err)
	}

	uc.log.Println("WAREHOUSE CREATED SUCCESSFULLY")
	return warehouse, nil
}

// |||||| retrieves a warehouse by ID  |||||
func (uc Usecase) GetWarehouse(warehouseId, userId uuid.UUID) (*entity.Warehouse, error) {
	uc.log.Println("FETCHING WAREHOUSE")

	warehouse, err := uc.repo.GetWarehouse(warehouseId, userId)
	if err != nil {
		uc.log.Println("ERROR FETCHING WAREHOUSE:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToGetWarehouse, err)
	}

	uc.log.Println("WAREHOUSE FETCHED SUCCESSFULLY")
	return warehouse, nil
}

// |||||||||| updates an existing warehouse's details
func (uc Usecase) UpdateWarehouse(warehouseIDParsed uuid.UUID, name, location string, capacity int, description string, userId uuid.UUID) (*entity.Warehouse, error) {
	warehouse, err := uc.repo.GetWarehouse(warehouseIDParsed, userId)
	if err != nil {
		uc.log.Println("ERROR FETCHING WAREHOUSE FOR UPDATE:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToUpdateWarehouse, err)
	}

	if name != "" {
		warehouse.Name = name
	}
	if location != "" {
		warehouse.Location = location
	}
	if capacity > 0 {
		warehouse.Capacity = capacity
	}
	if description != "" {
		warehouse.Description = description
	}

	warehouse.UpdatedAt = time.Now()
	warehouse.UpdatedBy = userId

	uc.log.Println("UPDATING WAREHOUSE")

	// pass updated values to the repository
	warehouse, err = uc.repo.UpdateWarehouse(
		warehouse.Id, warehouse.Name, warehouse.Location, warehouse.Capacity, description, warehouse.UpdatedBy,
	)
	if err != nil {
		uc.log.Println("ERROR UPDATING WAREHOUSE:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToUpdateWarehouse, err)
	}

	uc.log.Println("WAREHOUSE UPDATED SUCCESSFULLY")
	return warehouse, nil
}

func (uc Usecase) DeleteWarehouse(warehouseID uuid.UUID, userId uuid.UUID) error {
	uc.log.Println("DELETING WAREHOUSE")

	err := uc.repo.DeleteWarehouse(warehouseID, userId)
	if err != nil {
		uc.log.Println("ERROR DELETING WAREHOUSE:", err)
		return fmt.Errorf("%s: %w", ErrFailedToDeleteWarehouse, err)
	}

	uc.log.Println("WAREHOUSE DELETED SUCCESSFULLY")
	return nil
}

func (uc Usecase) ListWarehouses(userId uuid.UUID) ([]entity.Warehouse, error) {
	uc.log.Println("FETCHING WAREHOUSE LIST")

	warehouses, err := uc.repo.ListWarehouses(userId)
	if err != nil {
		uc.log.Println("ERROR FETCHING WAREHOUSE LIST:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToListWarehouses, err)
	}

	uc.log.Println("WAREHOUSE LIST FETCHED SUCCESSFULLY")
	return warehouses, nil
}

// retrieves a list of warehouses for a specific merchant
func (uc Usecase) ListMerchantWarehouses(userId uuid.UUID) ([]entity.Warehouse, error) {
	uc.log.Println("FETCHING MERCHANT WAREHOUSE LIST")

	warehouses, err := uc.repo.ListMerchantWarehouses(userId)
	if err != nil {
		uc.log.Println("ERROR FETCHING MERCHANT WAREHOUSE LIST:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToListWarehouses, err)
	}

	uc.log.Println("MERCHANT WAREHOUSE LIST FETCHED SUCCESSFULLY")
	return warehouses, nil
}

// deactivates a warehouse by ID
func (uc Usecase) DeactivateWarehouse(warehouseID uuid.UUID, userId uuid.UUID) (*entity.Warehouse, error) {
	uc.log.Println("DEACTIVATING WAREHOUSE")

	warehouse, err := uc.repo.GetWarehouse(warehouseID, userId)
	if err != nil {
		uc.log.Println("ERROR FETCHING WAREHOUSE FOR DEACTIVATION:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToDeactivateWarehouse, err)
	}

	// Set IsActive to false
	warehouse.IsActive = false
	warehouse.UpdatedAt = time.Now()
	warehouse.UpdatedBy = userId
	warehouse, err = uc.repo.UpdateWarehouse(
		warehouse.Id, warehouse.Name, warehouse.Location, warehouse.Capacity, warehouse.Description, warehouse.UpdatedBy,
	)
	if err != nil {
		uc.log.Println("ERROR DEACTIVATING WAREHOUSE:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToDeactivateWarehouse, err)
	}

	uc.log.Println("WAREHOUSE DEACTIVATED SUCCESSFULLY")
	return warehouse, nil
}

// validates the input fields when creating or updating a warehouse
func validateWarehouseInput(name, location string, capacity int, description string) error {
	if name == "" {
		return errors.New("warehouse name cannot be empty")
	}

	if location == "" {
		return errors.New("warehouse location cannot be empty")
	}

	if capacity <= 0 {
		return errors.New("warehouse capacity must be greater than zero")
	}

	if description == "" {
		return errors.New("warehouse description cannot be empty")
	}

	return nil
}
