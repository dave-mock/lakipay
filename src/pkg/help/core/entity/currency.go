package entity

import (
	"time"

	"github.com/google/uuid"
)

type CurrencyCountry struct {
	Id          uuid.UUID
	Default     bool
	Name        string
	DefaultName string
	Iso2        string
	Flag        string
}

type Currency struct {
	Id uuid.UUID
	// Name - dollar
	Name string
	// Currency - USD
	Currency string
	// Symbol - $
	Symbol    string
	Rate      float64
	BaseId    uuid.UUID
	Hidden    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Countries []CurrencyCountry
}
