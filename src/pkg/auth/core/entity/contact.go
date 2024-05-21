package entity

import (
	"auth/src/pkg/account/core/entity"
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	Id        uuid.UUID
	SirName   string
	FirstName string
	LastName  string
	User      struct {
		SirName       string
		FirstName     string
		LastName      string
		Gender        Gender
		DateOfBirth   time.Time
		Nationalities []Nationality
		Addresses     []Address
		Identities    []Identity
		Accounts      []struct {
			Id      uuid.UUID
			Type    entity.AccountType
			Default bool
			Detail  interface{}
			User    User
		}
		CreatedAt time.Time
	}
}
