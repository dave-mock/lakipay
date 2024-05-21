package entity

import (
	"time"

	"github.com/google/uuid"
)

type PasswordAuth struct {
	Id        uuid.UUID
	Token     string
	Password  PasswordIdentity
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
