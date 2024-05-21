package entity

import (
	"time"

	"github.com/google/uuid"
)

type PreSession struct {
	Id        uuid.UUID
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
