package usecase

import (
	"auth/src/pkg/org/core/entity"

	"github.com/google/uuid"
)

type Interactor interface {
	CheckTIN(tin string) (*entity.Organization, error)

	CreateCategory(
		name string,
		description string,
		icon string,
		parents []uuid.UUID,
		countryWhitelist []string,
		countryBlacklist []string,
		hidden bool,
		options []struct {
			Name             string
			Description      string
			DataType         entity.DataType
			RepresentedIn    string
			Values           []interface{}
			AllowCustomValue bool
			Validator        map[string]struct {
				Value   interface{}
				Message string
			}
		},
	) (*entity.Category, error)
	GetCategorys() ([]entity.Category, error)
	GetCategoryByName(name string) (*entity.Category, error)

	CreateLegalCondition(string, []string, []string) (*entity.LegalCondition, error)
	GetLegalConditions() ([]entity.LegalCondition, error)
	GetLegalConditionByName(name string) (*entity.LegalCondition, error)

	// Taxes
	CreateTax(
		name,
		description string,
		rate float64,
		from entity.TaxableEntity,
		countryWhitelist []string,
		countryBlacklist []string,
		hidden bool,
	) (*entity.Tax, error)
	GetTaxes() ([]entity.Tax, error)

	// Organization
	OrgInteractor
}
