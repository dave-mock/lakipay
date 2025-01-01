package entity

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	Id          uuid.UUID
	MerchantID  uuid.UUID
	Name        string
	Location    string
	Capacity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   uuid.UUID
	UpdatedBy   uuid.UUID
	IsActive    bool
	Description string
}
