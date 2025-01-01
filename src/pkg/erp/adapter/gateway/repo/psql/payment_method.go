package psql

import (
	"auth/src/pkg/erp/core/entity"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (repo PsqlRepo) CreatePaymentMethod(name, methodType string, commission float64, details string, isActive bool, createdBy uuid.UUID) error {
	const ErrFailedToStorePaymentMethod = "FAILED_TO_STORE_PAYMENT_METHOD"
	query := `
		INSERT INTO erp.payment_methods (name, method_type, commission, details, is_active, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, method_type, commission, details, is_active, created_at, updated_at, created_by, updated_by`

	var paymentMethod entity.PaymentMethod
	err := repo.db.QueryRow(
		query,
		name,
		methodType,
		commission,
		details,
		isActive,
		time.Now(),
		time.Now(),
		createdBy,
		createdBy).Scan(
		&paymentMethod.Id,
		&paymentMethod.Name,
		&paymentMethod.Type,
		&paymentMethod.Comission,
		&paymentMethod.Details,
		&paymentMethod.IsActive,
		&paymentMethod.CreatedAt,
		&paymentMethod.UpdatedAt,
		&paymentMethod.CreatedBy,
		&paymentMethod.UpdatedBy,
	)
	if err != nil {
		repo.log.Println("ERROR STORING PAYMENT METHOD:", err)
		return fmt.Errorf("%s: %v", ErrFailedToStorePaymentMethod, err)
	}

	repo.log.Println("PAYMENT METHOD CREATED SUCCESSFULLY")
	return nil
}

func (repo PsqlRepo) UpdatePaymentMethod(paymentMethodID uuid.UUID, name, methodType string, commission float64, details string, isActive bool, updatedBy uuid.UUID) error {
	const ErrFailedToUpdatePaymentMethod = "FAILED_TO_UPDATE_PAYMENT_METHOD"
	query := `
		UPDATE erp.payment_methods
		SET name = $1, method_type = $2, commission = $3, details = $4, is_active = $5, updated_at = $6, updated_by = $7
		WHERE id = $8
		RETURNING id`

	_, err := repo.db.Exec(query, name, methodType, commission, details, isActive, time.Now(), updatedBy, paymentMethodID)
	if err != nil {
		repo.log.Println("ERROR UPDATING PAYMENT METHOD:", err)
		return fmt.Errorf("%s: %v", ErrFailedToUpdatePaymentMethod, err)
	}

	repo.log.Println("PAYMENT METHOD UPDATED SUCCESSFULLY")
	return nil
}

func (repo PsqlRepo) GetPaymentMethod(paymentMethodID uuid.UUID, userID uuid.UUID) (*entity.PaymentMethod, error) {
	const ErrFailedToFetchPaymentMethod = "FAILED_TO_FETCH_PAYMENT_METHOD"
	query := `
		SELECT id, name, method_type, commission, details, is_active, created_at, updated_at, created_by, updated_by, merchant_id
		FROM erp.payment_methods
		WHERE id = $1 AND user_id = $2`

	var paymentMethod entity.PaymentMethod
	err := repo.db.QueryRow(query, paymentMethodID, userID).Scan(
		&paymentMethod.Id, &paymentMethod.Name, &paymentMethod.Type, &paymentMethod.Comission,
		&paymentMethod.Details, &paymentMethod.IsActive, &paymentMethod.CreatedAt, &paymentMethod.UpdatedAt,
		&paymentMethod.CreatedBy, &paymentMethod.UpdatedBy, &paymentMethod.MerchantID,
	)
	if err != nil {
		repo.log.Println("ERROR FETCHING PAYMENT METHOD:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToFetchPaymentMethod, err)
	}

	repo.log.Println("PAYMENT METHOD FETCHED SUCCESSFULLY")
	return &paymentMethod, nil
}

// fetches all payment methods for a user
func (repo PsqlRepo) ListPaymentMethods(userID uuid.UUID) ([]entity.PaymentMethod, error) {
	const ErrFailedToFetchPaymentMethods = "FAILED_TO_FETCH_PAYMENT_METHODS"
	query := `
		SELECT id, name, method_type, commission, details, is_active, created_at, updated_at, created_by, updated_by, user_id
		FROM erp.payment_methods
		WHERE user_id = $1`

	rows, err := repo.db.Query(query, userID)
	if err != nil {
		repo.log.Println("ERROR FETCHING PAYMENT METHODS:", err)
		return nil, fmt.Errorf("%s: %v", ErrFailedToFetchPaymentMethods, err)
	}
	defer rows.Close()

	var paymentMethods []entity.PaymentMethod
	for rows.Next() {
		var paymentMethod entity.PaymentMethod
		if err := rows.Scan(
			&paymentMethod.Id,
			&paymentMethod.Name,
			&paymentMethod.Type,
			&paymentMethod.Comission,
			&paymentMethod.Details,
			&paymentMethod.IsActive,
			&paymentMethod.CreatedAt,
			&paymentMethod.UpdatedAt,
			&paymentMethod.CreatedBy,
			&paymentMethod.UpdatedBy,
			&paymentMethod.MerchantID,
		); err != nil {
			repo.log.Println("ERROR SCANNING PAYMENT METHOD:", err)
			return nil, fmt.Errorf("%s: %v", ErrFailedToFetchPaymentMethods, err)
		}
		paymentMethods = append(paymentMethods, paymentMethod)
	}

	repo.log.Println("PAYMENT METHODS FETCHED SUCCESSFULLY")
	return paymentMethods, nil
}

// removes a payment method from the database
func (repo PsqlRepo) DeletePaymentMethod(paymentMethodID uuid.UUID, userID uuid.UUID) error {
	const ErrFailedToDeletePaymentMethod = "FAILED_TO_DELETE_PAYMENT_METHOD"
	query := `DELETE FROM erp.payment_methods WHERE id = $1 AND user_id = $2`
	_, err := repo.db.Exec(query, paymentMethodID, userID)
	if err != nil {
		repo.log.Println("ERROR DELETING PAYMENT METHOD:", err)
		return fmt.Errorf("%s: %v", ErrFailedToDeletePaymentMethod, err)
	}

	repo.log.Println("PAYMENT METHOD DELETED SUCCESSFULLY")
	return nil
}

func (repo PsqlRepo) DeactivatePaymentMethod(id string, userId uuid.UUID) error {
	const ErrFailedToDeactivatePaymentMethod = "FAILED_TO_DEACTIVATE_PAYMENT_METHOD"
	_, err := repo.db.Exec(`
		UPDATE erp.payment_methods
		SET status = $1, updated_by = $2
		WHERE id = $3
	`, "inactive", userId, id)
	if err != nil {
		repo.log.Println("ERROR DEACTIVATING PAYMENT METHOD:", err)
		return fmt.Errorf("%s: %v", ErrFailedToDeactivatePaymentMethod, err)
	}

	repo.log.Println("PAYMENT METHOD DEACTIVATED SUCCESSFULLY")
	return nil
}
