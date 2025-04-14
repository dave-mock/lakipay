package service

import "auth/src/pkg/checkout/core/entity"

type CheckoutInteractor interface {
	/// [GATEWAY]
	gatewayInteractor
	/// [TRANSACTION]
	transactionInteractor
}

type gatewayInteractor interface {
	GetGateways() ([]entity.Gateway, error)
	FindGatewayByKey(string) (*entity.Gateway, error)
}

type transactionInteractor interface {
	InitTransaction(to string, medium string, amount float64, redirects struct {
		Success string
		Cancel  string
		Decline string
	}, details map[string]interface{}) (*entity.Transaction, error)
	ConfirmTransaction(string) (any, error)
	InitDirectTransaction(to string, medium string, amount float64, details map[string]interface{}, callback string) (any, error)
	GetTransaction(string) (entity.Transaction, error)
	UpdatePaymentStatus(string, struct {
		Value   entity.TransactionStatus
		Message string
	}) error
}
