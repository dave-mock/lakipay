package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID
	MerchantID  uuid.UUID
	Name        string
	Description string
	Price       float64
	Currency    string
	SKU         string
	Weight      float64
	Dimensions  string
	ImageURL    string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   uuid.UUID
	UpdatedBy   uuid.UUID
}
