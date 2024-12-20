package entity

import (
	"github.com/google/uuid"
)

type PaymentMethod struct {
	Id         uuid.UUID
	MerchantID string `json:"merchant_id"`
	Name       string
	Type       string
	Comission  float64
	Details    string
	IsActive   bool
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
	CreatedBy  uuid.UUID `json:"created_by"`
	UpdatedBy  uuid.UUID `json:"updated_by"`
}
