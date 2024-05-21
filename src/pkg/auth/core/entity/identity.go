package entity

import (
	"time"

	"github.com/google/uuid"
)

type Identity interface{}

type PhoneIdentity struct {
	Id        uuid.UUID
	User      User
	Phone     Phone
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PasswordIdentity struct {
	Id        uuid.UUID
	User      User
	Password  string
	Hint      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// [TODO] - will go on
