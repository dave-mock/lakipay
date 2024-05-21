package usecase

import "auth/src/pkg/org/core/entity"

type Repo interface {
	StoreCategory(entity.Category) error
	FindCategories() ([]entity.Category, error)
	FindCategoryByName(name string) (*entity.Category, error)

	StoreLegalCondition(entity.LegalCondition) error
	FindLegalConditions() ([]entity.LegalCondition, error)
	FindLegalConditionByName(name string) (*entity.LegalCondition, error)

	// Taxes
	StoreTax(v entity.Tax) error
	FindTaxes() ([]entity.Tax, error)
}
