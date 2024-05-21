package entity

import (
	"time"

	"github.com/google/uuid"
)

type TaxableEntity string

const (
	PAYER TaxableEntity = "PAYER"
	PAYEE TaxableEntity = "PAYEE"
)

type Tax struct {
	Id               uuid.UUID
	Name             string
	Rate             float64
	From             TaxableEntity
	Description      string
	CountryWhitelist []string
	CountryBlacklist []string
	Hidden           bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
