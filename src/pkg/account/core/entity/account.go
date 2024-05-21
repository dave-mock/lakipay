package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	STORED AccountType = "STORED"
	BANK   AccountType = "BANK"
	CARD   AccountType = "CARD"
	WALLET AccountType = "WALLET"
)

type Transactor interface {
	// +
	Credit(from Account, amount float64)
	// -
	Debit(to Account, amount float64)
}

type Account struct {
	Id                 uuid.UUID
	Title              string
	Type               AccountType
	Default            bool
	Detail             Transactor
	User               User
	VerificationStatus struct {
		Verified   bool
		VerifiedBy *struct {
			Method  string
			Details interface{}
		}
	}
	CreatedAt time.Time
	UpdatedAt time.Time
}
