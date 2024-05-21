package usecase

import (
	"auth/src/pkg/org/core/entity"
	"time"

	"github.com/google/uuid"
)

func (uc Usecase) CreateTax(name, description string, rate float64, from entity.TaxableEntity, countryWhitelist []string, countryBlacklist []string, hidden bool) (*entity.Tax, error) {
	tax := entity.Tax{
		Id:               uuid.New(),
		Name:             name,
		Description:      description,
		Rate:             rate,
		From:             from,
		CountryWhitelist: countryWhitelist,
		CountryBlacklist: countryBlacklist,
		Hidden:           hidden,
		CreatedAt:        time.Now(),
	}

	// Store
	err := uc.repo.StoreTax(tax)
	if err != nil {
		return nil, Error{
			Type:    "CREATE_TAX",
			Message: err.Error(),
		}
	}

	// Return
	return &tax, nil
}

func (uc Usecase) GetTaxes() ([]entity.Tax, error) {
	var taxes []entity.Tax

	taxes, err := uc.repo.FindTaxes()
	if err != nil {
		return nil, Error{
			Type:    "FIND_TAXES",
			Message: err.Error(),
		}
	}

	return taxes, nil
}
