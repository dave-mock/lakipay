package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductCatalog struct {
	Id           uuid.UUID
	MerchantID   uuid.UUID
	CatalogID    uuid.UUID
	ProductID    uuid.UUID
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedBy    uuid.UUID
	UpdatedBy    uuid.UUID
}
