package usecase

import (
	"auth/src/pkg/help/core/entity"
	"time"

	"github.com/google/uuid"
)

func (uc Usecase) CreateCountry(name, defaultName, iso2, flag string, phonePrefix entity.PhonePrefix, hidden bool) (*entity.Country, error) {

	// Error
	var ErrFailedToCreateCountry string = "FAILED_TO_CREATE_COUNTRY"

	// Validate

	// Create
	var country entity.Country

	id := uuid.New()

	country = entity.Country{
		Id:          id,
		Name:        name,
		DefaultName: defaultName,
		Iso2:        iso2,
		Flag:        flag,
		PhonePrefix: phonePrefix,
		CreatedAt:   time.Now(),
	}

	// Store
	err := uc.repo.StoreCountry(country)
	if err != nil {
		return nil, &Error{
			Type:    ErrFailedToCreateCountry,
			Message: err.Error(),
		}
	}

	// Return
	return &country, nil
}

func (uc Usecase) GetCountries() ([]entity.Country, error) {
	// Error
	var ErrCouldNotFindCountries string = "COULD_NOT_FIND_COUNTRIES"

	// Fetch countries
	countries, err := uc.repo.FindCountries()
	if err != nil {
		return make([]entity.Country, 0), &Error{
			Type:    ErrCouldNotFindCountries,
			Message: err.Error(),
		}
	}

	// Return
	return countries, nil
}

func (uc Usecase) GetSupportedCountries() ([]entity.Country, error) {
	// Error
	var ErrCouldNotFindCountries string = "COULD_NOT_FIND_COUNTRIES"

	// Fetch countries
	countries, err := uc.repo.FindSupportedCountries()
	if err != nil {
		return make([]entity.Country, 0), &Error{
			Type:    ErrCouldNotFindCountries,
			Message: err.Error(),
		}
	}

	// Return
	return countries, nil
}

func (uc Usecase) SearchAndFilterCountries(query string, filterKey string) ([]entity.Country, error) {
	var countries []entity.Country = make([]entity.Country, 0)

	// Search countries
	// countries, err := uc.repo.SearchCountries

	// Return

	return countries, nil
}
