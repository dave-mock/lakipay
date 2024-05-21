package usecase

import (
	"auth/src/pkg/org/core/entity"
	"time"

	"github.com/google/uuid"
)

func (uc Usecase) GetCategorys() ([]entity.Category, error) {
	var categories []entity.Category = make([]entity.Category, 0)

	// Fetch from repo
	categories, err := uc.repo.FindCategories()
	uc.log.Println(categories)
	if err != nil {
		return nil, Error{
			Type:    "ErrFetchingcategories",
			Message: err.Error(),
		}
	}

	return categories, nil
}

func (uc Usecase) CreateCategory(
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
	}) (*entity.Category, error) {

	var _options []entity.Option = make([]entity.Option, 0)

	now := time.Now()

	for i := 0; i < len(options); i++ {

		var validator entity.OptionValidator = entity.OptionValidator{}

		for k, v := range options[i].Validator {
			validator[k] = struct {
				Value   interface{}
				Message interface{}
			}{
				Value:   v.Value,
				Message: v.Message,
			}
		}

		_options = append(_options, entity.Option{
			Id:               uuid.New(),
			Name:             options[i].Name,
			Description:      options[i].Description,
			DataType:         options[i].DataType,
			RepresentedIn:    options[i].RepresentedIn,
			Values:           options[i].Values,
			AllowCustomValue: options[i].AllowCustomValue,
			Validator:        validator,
			CreatedAt:        now,
		})
	}

	var category entity.Category = entity.Category{
		Id:               uuid.New(),
		Name:             name,
		Description:      description,
		Icon:             icon,
		Parents:          parents,
		CountryWhitelist: countryWhitelist,
		CountryBlacklist: countryBlacklist,
		Hidden:           hidden,
		Options:          _options,
		CreatedAt:        now,
	}

	// Store
	err := uc.repo.StoreCategory(category)
	if err != nil {
		return nil, Error{
			Type:    "ErrAddingOrgType",
			Message: err.Error(),
		}
	}

	return &category, nil
}

func (uc Usecase) GetCategoryByName(name string) (*entity.Category, error) {
	return uc.repo.FindCategoryByName(name)
}
