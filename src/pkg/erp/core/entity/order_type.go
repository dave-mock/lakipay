package entity

import (
	"github.com/google/uuid"
)

type OrderType struct {
	ID          uuid.UUID `json:"id"`
	TypeName    string    `json:"type_name"`
	MerchantID  string    `json:"merchant_id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	CreatedBy   uuid.UUID `json:"created_by"`
	UpdatedBy   uuid.UUID `json:"updated_by"`
}
