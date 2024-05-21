package usecase

import (
	"auth/src/pkg/help/core/entity"
	"time"

	"github.com/google/uuid"
)

func (uc Usecase) CreateCurrency(name, currency, symbol string, rate float64, baseId uuid.UUID, countries []entity.CurrencyCountry) (*entity.Currency, error) {

	// Validate Input

	// Find base currency
	// _, err := uc.repo.FindCurrencyById(baseId)
	// if err != nil {
	// 	return nil, Error{
	// 		Type:    "CURRENCY_BASE_NOT_FOUND",
	// 		Message: err.Error(),
	// 	}
	// }

	id := uuid.New()

	if baseId == uuid.Nil {
		baseId = id
	}

	// Create currency
	curr := entity.Currency{
		Id:        id,
		Name:      name,
		Currency:  currency,
		Symbol:    symbol,
		Rate:      rate,
		BaseId:    baseId,
		Hidden:    false,
		CreatedAt: time.Now(),
		Countries: countries,
	}

	// Store currency
	err := uc.repo.StoreCurrency(curr)
	if err != nil {
		return nil, Error{
			Type:    "FAILED_TO_STORE_CURRENCY",
			Message: err.Error(),
		}
	}

	return &curr, nil
}

func (uc Usecase) GetCurrencies() ([]entity.Currency, error) {
	// Error
	var ErrCouldNotFindCurrencies string = "COULD_NOT_FIND_CURRENCIES"

	// Fetch currencies
	currencies, err := uc.repo.FindCurrencies()
	if err != nil {
		return make([]entity.Currency, 0), &Error{
			Type:    ErrCouldNotFindCurrencies,
			Message: err.Error(),
		}
	}

	// Return
	return currencies, nil
}

func (uc Usecase) GetSupportedCurrencies() ([]entity.Currency, error) {
	// Error
	var ErrCouldNotFindCurrencies string = "COULD_NOT_FIND_CURRENCIES"

	// Fetch currencies
	currencies, err := uc.repo.FindSupportedCurrencies()
	if err != nil {
		return make([]entity.Currency, 0), &Error{
			Type:    ErrCouldNotFindCurrencies,
			Message: err.Error(),
		}
	}

	// Return
	return currencies, nil
}
