package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id        uuid.UUID
	Token     string
	User      User
	Device    Device
	CreatedAt time.Time
	UpdatedAt time.Time
}
