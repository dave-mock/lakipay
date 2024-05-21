package entity

import (
	"time"

	"github.com/google/uuid"
)

type Bank struct {
	Id        uuid.UUID
	Name      string
	ShortName string
	BIN       string
	SwiftCode string
	Logo      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
