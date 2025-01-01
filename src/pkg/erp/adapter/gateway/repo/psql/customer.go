package psql

import (
	"auth/src/pkg/erp/core/entity"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (repo PsqlRepo) CreateCustomer(customerID uuid.UUID, name, email, phone, address string, loyaltyPoints int, dateOfBirth, status string, createdBy, merchantID uuid.UUID) (*entity.Customer, error) {
	const ErrFailedToStoreCustomer = "FAILED_TO_STORE_CUSTOMER"
	query := `
		INSERT INTO erp.customers (id, name, email, phone, address, loyalty_points, date_of_birth, status, created_at, updated_at, created_by, updated_by, merchant_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, name, email, phone, address, loyalty_points, date_of_birth, status, created_at, updated_at, created_by, updated_by, merchant_id`

	var customer entity.Customer
	err := repo.db.QueryRow(query, customerID, name, email, phone, address, loyaltyPoints, dateOfBirth, status, time.Now(), time.Now(), createdBy, createdBy, merchantID).Scan(
		&customer.Id, &customer.Name, &customer.Email, &customer.Phone, &customer.Address, &customer.LoyaltyPoints,
		&customer.DateOfBirth, &customer.Status, &customer.CreatedAt, &customer.UpdatedAt, &customer.CreatedBy, &customer.UpdatedBy, &customer.MerchantID,
	)
	if err != nil {
		repo.log.Println("ERROR STORING CUSTOMER:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToStoreCustomer, err)
	}

	repo.log.Println("CUSTOMER CREATED SUCCESSFULLY")
	return &customer, nil
}

func (repo PsqlRepo) UpdateCustomer(customerID uuid.UUID, name, email, phone, address string, loyaltyPoints int, dateOfBirth, status string, updatedBy, merchantID uuid.UUID) (*entity.Customer, error) {
	const ErrFailedToUpdateCustomer = "FAILED_TO_UPDATE_CUSTOMER"
	query := `
		UPDATE erp.customers
		SET name = $1, email = $2, phone = $3, address = $4, loyalty_points = $5, date_of_birth = $6, status = $7, updated_at = $8, updated_by = $9
		WHERE id = $10
		RETURNING id, name, email, phone, address, loyalty_points, date_of_birth, status, created_at, updated_at, created_by, updated_by, merchant_id`

	var customer entity.Customer
	err := repo.db.QueryRow(query, name, email, phone, address, loyaltyPoints, dateOfBirth, status, time.Now(), updatedBy, customerID).Scan(
		&customer.Id, &customer.Name, &customer.Email, &customer.Phone, &customer.Address, &customer.LoyaltyPoints,
		&customer.DateOfBirth, &customer.Status, &customer.CreatedAt, &customer.UpdatedAt, &customer.CreatedBy, &customer.UpdatedBy, &customer.MerchantID,
	)
	if err != nil {
		repo.log.Println("ERROR UPDATING CUSTOMER:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToUpdateCustomer, err)
	}

	repo.log.Println("CUSTOMER UPDATED SUCCESSFULLY")
	return &customer, nil
}

func (repo PsqlRepo) GetCustomerByID(customerID uuid.UUID) (*entity.Customer, error) {
	const ErrFailedToFetchCustomer = "FAILED_TO_FETCH_CUSTOMER"
	query := `
		SELECT id, name, email, phone, address, loyalty_points, date_of_birth, status, created_at, updated_at, created_by, updated_by, merchant_id
		FROM erp.customers
		WHERE id = $1`

	var customer entity.Customer
	err := repo.db.QueryRow(query, customerID).Scan(
		&customer.Id,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.LoyaltyPoints,
		&customer.DateOfBirth,
		&customer.Status,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.CreatedBy,
		&customer.UpdatedBy,
		&customer.MerchantID,
	)
	if err != nil {
		repo.log.Println("ERROR FETCHING CUSTOMER:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToFetchCustomer, err)
	}

	repo.log.Println("CUSTOMER FETCHED SUCCESSFULLY")
	return &customer, nil
}

func (repo PsqlRepo) DeleteCustomer(customerID uuid.UUID) error {
	const ErrFailedToDeleteCustomer = "FAILED_TO_DELETE_CUSTOMER"
	query := `DELETE FROM erp.customers WHERE id = $1`
	_, err := repo.db.Exec(query, customerID)
	if err != nil {
		repo.log.Println("ERROR DELETING CUSTOMER:", err)
		return fmt.Errorf("%s: %v", ErrFailedToDeleteCustomer, err)
	}

	repo.log.Println("CUSTOMER DELETED SUCCESSFULLY")
	return nil
}
