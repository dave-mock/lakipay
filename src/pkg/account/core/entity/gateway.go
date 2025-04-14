package entity

import "time"

type Gateway struct {
	Id      string
	Key     string
	Name    string
	Acronym string

	Icon string
	Type AccountType

	CanProcess bool
	CanSettle  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
