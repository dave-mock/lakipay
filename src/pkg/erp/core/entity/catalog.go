package entity

import (
	"time"

	"github.com/google/uuid"
)

type Catalog struct {
	Id          uuid.UUID
	MerchantId  uuid.UUID
	Name        string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   uuid.UUID
	UpdatedBy   uuid.UUID
}
