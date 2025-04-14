package entity

import "time"

type AccountType string

const (
	LAKIPAY AccountType = "LAKIPAY"
	BANK    AccountType = "BANK"
	CARD    AccountType = "CARD"
	WALLET  AccountType = "WALLET"
)

type Gateway struct {
	Id      string
	Key     string
	Name    string
	Acronym string

	Icon string
	Type AccountType

	CanProcess bool
	CanSettle  bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
