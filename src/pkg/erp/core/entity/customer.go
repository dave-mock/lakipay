package entity

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	Id            uuid.UUID `json:"id"`
	CustomerID    uuid.UUID `json:"customer_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	LoyaltyPoints int       `json:"loyalty_points"`
	DateOfBirth   string    `json:"date_of_birth"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     uuid.UUID `json:"created_by"`
	UpdatedBy     uuid.UUID `json:"updated_by"`
	MerchantID    uuid.UUID `json:"merchant_id"`
}
