package service

import "auth/src/pkg/checkout/core/entity"

type CheckoutRepo interface {
	FindGateways() ([]entity.Gateway, error)
	FindGatewayByKey(string) (*entity.Gateway, error)

	/// [TRANSACTION]
	StoreTransaction(entity.Transaction) error
	FindTransaction(string) (*entity.Transaction, error)
	UpdateTransaction(entity.Transaction) error
}
