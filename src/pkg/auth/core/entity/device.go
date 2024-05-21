package entity

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type Device struct {
	Id        uuid.UUID
	IP        net.IPAddr
	Name      string
	Agent     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeviceAuth struct {
	Id        uuid.UUID
	Token     string
	Device    Device
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
