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
type Catalog struct {
	Id          uuid.UUID  `json:"id" db:"id"`
	MerchantId  uuid.UUID  `json:"merchant_id" db:"merchant_id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy   *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy   *uuid.UUID `json:"updated_by" db:"updated_by"`
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
