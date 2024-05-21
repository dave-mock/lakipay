package usecase

import (
	"auth/src/pkg/help/core/entity"

	"github.com/google/uuid"
)

type Repo interface {
	// Phone Prefix
	StorePhonePrefix(entity.PhonePrefix) error
	FindPhoneprefixByPrefix(prefix string) (*entity.PhonePrefix, error)
	// Country
	StoreCountry(entity.Country) error
	FindCountryById(uuid.UUID) (*entity.Country, error)
	FindCountries() ([]entity.Country, error)
	FindSupportedCountries() ([]entity.Country, error)
	// Currency
	StoreCurrency(entity.Currency) error
	FindCurrencyById(uuid.UUID) (*entity.Currency, error)
	FindCurrencies() ([]entity.Currency, error)
	FindSupportedCurrencies() ([]entity.Currency, error)
}
