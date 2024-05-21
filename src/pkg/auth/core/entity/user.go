package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id            uuid.UUID
	SirName       string
	FirstName     string
	LastName      string
	Gender        Gender
	DateOfBirth   time.Time
	Nationalities []Nationality
	Addresses     []Address
	Identities    []Identity
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Gender string

var (
	MALE   Gender = "MALE"
	FEMALE Gender = "FEMALE"
)
