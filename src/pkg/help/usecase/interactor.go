package usecase

import (
	"auth/src/pkg/help/core/entity"

	"github.com/google/uuid"
)

type Interactor interface {
	CreatePhonePrefix(prefix, pattern string) (*entity.PhonePrefix, error)
	GetPhonePrefixByPrefix(prefix string) (*entity.PhonePrefix, error)

	// Country
	CreateCountry(name, defaultName, iso2, flag string, phonePrefix entity.PhonePrefix, hidden bool) (*entity.Country, error)
	GetCountries() ([]entity.Country, error)
	GetSupportedCountries() ([]entity.Country, error)
	SearchAndFilterCountries(query string, filterKey string) ([]entity.Country, error)

	// Currency
	CreateCurrency(name, currency, symbol string, rate float64, baseId uuid.UUID, countries []entity.CurrencyCountry) (*entity.Currency, error)
	GetCurrencies() ([]entity.Currency, error)
	GetSupportedCurrencies() ([]entity.Currency, error)
}
