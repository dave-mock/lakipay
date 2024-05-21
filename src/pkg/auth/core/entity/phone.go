package entity

import (
	"time"

	"github.com/google/uuid"
)

type Phone struct {
	Id        uuid.UUID
	Prefix    string
	Number    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (phone Phone) String() string {
	return "+" + phone.Prefix + phone.Number
}

type PhoneAuth struct {
	Id        uuid.UUID
	Token     string
	Phone     Phone
	Code      string
	Status    bool
	Method    string
	Length    int64
	Timeout   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
