package entity

import (
	"time"

	"github.com/google/uuid"
)

type SubCatalog struct {
	ID           uuid.UUID `json:"id"`
	CatalogID    uuid.UUID `json:"catalog_id"`
	SubCatalogID uuid.UUID `json:"sub_catalog_id"`
	MerchantID   uuid.UUID `json:"merchant_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    uuid.UUID `json:"created_by"`
	UpdatedBy    uuid.UUID `json:"updated_by"`
}
