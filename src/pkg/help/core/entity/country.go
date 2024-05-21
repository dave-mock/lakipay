package entity

import (
	"time"

	"github.com/google/uuid"
)

type PhonePrefix struct {
	Prefix    string
	Pattern   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Country struct {
	Id          uuid.UUID
	Name        string
	DefaultName string
	Iso2        string
	Flag        string
	PhonePrefix PhonePrefix
	Hidden      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
