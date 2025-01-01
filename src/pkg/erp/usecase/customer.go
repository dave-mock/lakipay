package usecase

import (
	"auth/src/pkg/erp/core/entity"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// custom error struct to handle structured error responses
type JSONError struct {
	ErrorType string `json:"errorType"`
	Message   string `json:"message"`
}

// creates a new customer with validation
func (uc Usecase) CreateCustomer(customerID uuid.UUID, name, email, phone, address string, loyaltyPoints int, dateOfBirth, status string, createdBy, merchantID uuid.UUID) (*entity.Customer, error) {
	const ErrFailedToStoreCustomer = "FAILED_TO_STORE_CUSTOMER"
	const ErrInvalidInput = "INVALID_INPUT"

	// Validate inputs
	if err := validateCustomerInput(
		name,
		email,
		phone,
		address,
		loyaltyPoints,
		dateOfBirth,
		status,
	); err != nil {
		return nil, err // Return standard error
	}

	// Create the customer entity
	customer := &entity.Customer{
		Id:            uuid.New(),
		Name:          name,
		Email:         email,
		Phone:         phone,
		Address:       address,
		LoyaltyPoints: loyaltyPoints,
		DateOfBirth:   dateOfBirth,
		Status:        status,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		CreatedBy:     createdBy,
		UpdatedBy:     createdBy,
		MerchantID:    merchantID,
	}

	// Log the creation attempt
	uc.log.Println("CREATING CUSTOMER")
	customer, err := uc.repo.CreateCustomer(
		customerID,
		name,
		email,
		phone,
		address,
		loyaltyPoints,
		dateOfBirth,
		status, createdBy,
		merchantID)
	if err != nil {
		uc.log.Println("ERROR STORING CUSTOMER")
		return nil, fmt.Errorf("%s: %v", ErrFailedToStoreCustomer, err)
	}

	uc.log.Println("CUSTOMER CREATED SUCCESSFULLY")
	return customer, nil
}

// updates an existing customer with validation
func (uc Usecase) UpdateCustomer(customerID uuid.UUID, name, email, phone, address string, loyaltyPoints int, dateOfBirth, status string, updatedBy, merchantID uuid.UUID) (*entity.Customer, error) {
	const ErrFailedToUpdateCustomer = "FAILED_TO_UPDATE_CUSTOMER"
	const ErrInvalidInput = "INVALID_INPUT"

	// Validate inputs
	if err := validateCustomerInput(
		name,
		email,
		phone,
		address,
		loyaltyPoints,
		dateOfBirth,
		status,
	); err != nil {
		return nil, err // Return standard error
	}

	// Update customer in the repository
	customer, err := uc.repo.UpdateCustomer(
		customerID,
		name, email,
		phone,
		address,
		loyaltyPoints,
		dateOfBirth,
		status, updatedBy,
		merchantID,
	)
	if err != nil {
		uc.log.Println("ERROR UPDATING CUSTOMER")
		return nil, fmt.Errorf("%s: %v", ErrFailedToUpdateCustomer, err)
	}

	uc.log.Println("CUSTOMER UPDATED SUCCESSFULLY")
	return customer, nil
}

// validates the customer input and returns an error if invalid
func validateCustomerInput(name, email, phone, address string, loyaltyPoints int, dateOfBirth, status string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}
	if strings.TrimSpace(phone) == "" {
		return fmt.Errorf("phone cannot be empty")
	}
	if strings.TrimSpace(address) == "" {
		return fmt.Errorf("address cannot be empty")
	}
	if loyaltyPoints < 0 {
		return fmt.Errorf("loyalty points cannot be negative")
	}
	if strings.TrimSpace(dateOfBirth) == "" {
		return fmt.Errorf("date of birth cannot be empty")
	}
	if strings.TrimSpace(status) == "" {
		return fmt.Errorf("status cannot be empty")
	}
	return nil
}

// checks if the email format is valid
func isValidEmail(email string) bool {
	const emailRegex = `^[a-z0-9]+@[a-z0-9]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
