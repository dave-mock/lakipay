package usecase

import (
	"auth/src/pkg/org/core/entity"
	"time"

	"github.com/google/uuid"
)

func (uc Usecase) CreateLegalCondition(name string, countryWhitelist []string, countryBlacklist []string) (*entity.LegalCondition, error) {

	// Create entity
	legalCondition := entity.LegalCondition{
		Id:               uuid.New(),
		Name:             name,
		CountryWhitelist: countryWhitelist,
		CountryBlacklist: countryBlacklist,
		CreatedAt:        time.Now(),
	}

	// Store entity
	err := uc.repo.StoreLegalCondition(legalCondition)
	if err != nil {
		return nil, Error{
			Type:    "ErrCreateLegalCondition",
			Message: err.Error(),
		}
	}

	return &legalCondition, nil
}

func (uc Usecase) GetLegalConditions() ([]entity.LegalCondition, error) {
	var legalConditions []entity.LegalCondition = make([]entity.LegalCondition, 0)

	legalConditions, err := uc.repo.FindLegalConditions()
	if err != nil {
		return nil, Error{
			Type:    "",
			Message: err.Error(),
		}
	}

	return legalConditions, nil
}

func (uc Usecase) GetLegalConditionByName(name string) (*entity.LegalCondition, error) {
	return uc.repo.FindLegalConditionByName(name)
}
