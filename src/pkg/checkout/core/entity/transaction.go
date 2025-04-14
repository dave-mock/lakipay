package entity

import "time"

type TransactionStatus string

const (
	TxnPending    TransactionStatus = "pending"
	TxnProcessing TransactionStatus = "processing"
	TxnCompleted  TransactionStatus = "completed"
	TxnCanceled   TransactionStatus = "canceled"
	TxnDeclined   TransactionStatus = "declined"
)

type Transaction struct {
	Id string

	To string

	For string

	Ttl int

	Pricing struct {
		Amount float64
		Fees   []map[string]float64
	}

	Status struct {
		Value   TransactionStatus
		Message string
	}
	GateWay string
	Type    string

	Details map[string]interface{}

	NotifyUrl string

	CreatedAt time.Time
	UpdatedAt time.Time
}
