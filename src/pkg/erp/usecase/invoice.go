package usecase

import (
	"auth/src/pkg/erp/core/entity"
	"fmt"

	"github.com/google/uuid"
)

const (
	ErrFailedToCreateMerchantInvoice = "Failed to generate invoice"
)

func (uc Usecase) CreateMerchantInvoice(userId uuid.UUID, orderID uuid.UUID) (*entity.Order, error) {
	uc.log.Println("Fetching merchant order for invoice", "userId", userId, "orderID", orderID)
	order, err := uc.repo.CreateMerchantInvoice(userId, orderID)
	if err != nil {
		uc.log.Println(ErrFailedToCreateMerchantInvoice, err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToCreateMerchantInvoice, err)
	}

	if order == nil {
		errMsg := "no order found for the provided userId and orderID"
		uc.log.Println(ErrFailedToCreateMerchantInvoice, errMsg)
		return nil, fmt.Errorf("%s: %s", ErrFailedToCreateMerchantInvoice, errMsg)
	}

	uc.log.Println("Merchant Invoice created successfully", "userId", userId, "orderID", orderID)
	return order, nil
}
