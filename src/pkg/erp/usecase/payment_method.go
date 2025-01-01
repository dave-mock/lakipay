package usecase

import (
	"auth/src/pkg/erp/core/entity"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Error constants
const (
	ErrFailedToCreatePaymentMethod     = "Failed to create payment method"
	ErrFailedToUpdatePaymentMethod     = "Failed to update payment method"
	ErrFailedToGetPaymentMethod        = "Failed to get payment method"
	ErrFailedToDeactivatePaymentMethod = "Failed to deactivate payment method"
	ErrFailedToListPaymentMethods      = "Failed to list payment methods"
)

func (uc Usecase) CreatePaymentMethod(name string, methodType string, commission float64, details string, isActive bool, userId uuid.UUID) error {
	if err := validatePaymentMethodInput(name, methodType, commission, details); err != nil {
		uc.log.Println(ErrFailedToCreatePaymentMethod, err)
		return fmt.Errorf("%s: %w", ErrFailedToCreatePaymentMethod, err)
	}

	paymentMethod := &entity.PaymentMethod{
		Id:        uuid.New(),
		Name:      name,
		Type:      methodType,
		Comission: commission,
		Details:   details,
		IsActive:  isActive,
	}

	// Log the action
	uc.log.Print("Creating payment method", "name", name, "userId", userId)
	err := uc.repo.CreatePaymentMethod(
		paymentMethod.Name,
		paymentMethod.Type,
		paymentMethod.Comission,
		paymentMethod.Details,
		paymentMethod.IsActive,
		userId,
	)
	if err != nil {
		uc.log.Println("Error creating payment method", "error", err)
		return fmt.Errorf("%s: %w", ErrFailedToCreatePaymentMethod, err)
	}

	// Log success
	uc.log.Println("Payment method created successfully", "paymentMethodId", paymentMethod.Id)
	return nil
}

func (uc Usecase) UpdatePaymentMethod(id uuid.UUID, name string, methodType string, commission float64, details string, isActive *bool, userId uuid.UUID) error {
	paymentMethod, err := uc.repo.GetPaymentMethod(id, userId)
	if err != nil {
		uc.log.Println(ErrFailedToUpdatePaymentMethod, err)
		return fmt.Errorf("%s: %w", ErrFailedToUpdatePaymentMethod, err)
	}

	if name != "" {
		paymentMethod.Name = name
	}
	if methodType != "" {
		paymentMethod.Type = methodType
	}
	if commission < 0 {
		paymentMethod.Comission = commission
	}
	if details != "" {
		paymentMethod.Details = details
	}
	if isActive != nil {
		paymentMethod.IsActive = *isActive
	}

	uc.log.Println("Updating payment method", "paymentMethodId", paymentMethod.Id)
	err = uc.repo.UpdatePaymentMethod(
		paymentMethod.Id,
		paymentMethod.Name,
		paymentMethod.Type,
		paymentMethod.Comission,
		paymentMethod.Details,
		paymentMethod.IsActive, userId)
	if err != nil {
		uc.log.Println("Error updating payment method", "error", err)
		return fmt.Errorf("%s: %w", ErrFailedToUpdatePaymentMethod, err)
	}

	uc.log.Println("Payment method updated successfully", "paymentMethodId", paymentMethod.Id)
	return nil
}

func (uc Usecase) ListPaymentMethods(userId uuid.UUID) ([]entity.PaymentMethod, error) {
	uc.log.Println("Fetching list of payment methods", "userId", userId)
	paymentMethods, err := uc.repo.ListPaymentMethods(userId)
	if err != nil {
		uc.log.Println(ErrFailedToListPaymentMethods, err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToListPaymentMethods, err)
	}

	uc.log.Println("Payment methods fetched successfully", "userId", userId)
	return paymentMethods, nil
}

func (uc Usecase) GetPaymentMethod(id uuid.UUID, userId uuid.UUID) (*entity.PaymentMethod, error) {
	uc.log.Println("Fetching payment method", "paymentMethodId", id, "userId", userId)
	paymentMethod, err := uc.repo.GetPaymentMethod(id, userId)
	if err != nil {
		uc.log.Println(ErrFailedToGetPaymentMethod, err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToGetPaymentMethod, err)
	}

	uc.log.Println("Payment method fetched successfully", "paymentMethodId", id)
	return paymentMethod, nil
}

func (uc Usecase) DeactivatePaymentMethod(id string, userId uuid.UUID) error {
	uc.log.Println("Deactivating payment method", "paymentMethodId", id, "userId", userId)
	err := uc.repo.DeactivatePaymentMethod(id, userId)
	if err != nil {
		uc.log.Println(ErrFailedToDeactivatePaymentMethod, err)
		return fmt.Errorf("%s: %w", ErrFailedToDeactivatePaymentMethod, err)
	}

	uc.log.Println("Payment method deactivated successfully", "paymentMethodId", id)
	return nil
}

func validatePaymentMethodInput(name string, methodType string, commission float64, details string) error {
	if name == "" || methodType == "" || commission < 0 || details == "" {
		return errors.New("payment method details are missing or invalid")
	}
	return nil
}
