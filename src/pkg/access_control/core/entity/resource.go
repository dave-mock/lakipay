package entity

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Operations  []uuid.UUID `json:"operations" gorm:"type:uuid[]"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
