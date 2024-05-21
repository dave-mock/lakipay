package rest

import (
	"auth/src/pkg/org/core/entity"
	"auth/src/pkg/org/usecase"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Category struct {
	Id                 uuid.UUID   `json:"id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	Icon               string      `json:"icon"`
	Parents            []uuid.UUID `json:"parents"`
	CountriesWhitelist []string    `json:"country_whitelist"`
	CountriesBlacklist []string    `json:"country_blacklist"`
	Hidden             bool        `json:"hidden"`
	Options            []Option    `json:"options"`
}

func NewCategoryFromJson(i entity.Category) Category {

	var options []Option = make([]Option, 0)

	for j := 0; j < len(i.Options); j++ {
		options = append(options, OptionFromEntity(i.Options[j]))
	}

	return Category{
		Id:                 i.Id,
		Name:               i.Name,
		Description:        i.Description,
		Icon:               i.Icon,
		Parents:            i.Parents,
		CountriesWhitelist: i.CountryWhitelist,
		CountriesBlacklist: i.CountryBlacklist,
		Hidden:             i.Hidden,
		Options:            options,
	}
}

type Option struct {
	Id               string          `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	DataType         string          `json:"data_type"`
	RepresentedIn    string          `json:"represented_in"`
	Values           []interface{}   `json:"values"`
	AllowCustomValue bool            `json:"allow_custom_value"`
	Validator        OptionValidator `json:"validator"`
}

func OptionFromEntity(i entity.Option) Option {
	return Option{
		Id:               i.Id.String(),
		Name:             i.Name,
		Description:      i.Description,
		DataType:         string(i.DataType),
		RepresentedIn:    i.RepresentedIn,
		Values:           i.Values,
		AllowCustomValue: i.AllowCustomValue,
		Validator:        OptionValidator(i.Validator),
	}
}

type OptionValidator map[string]struct {
	Value   interface{} `json:"value"`
	Message interface{} `json:"message"`
}

func CategoryFromEntity(i entity.Category) Category {

	var options []Option = make([]Option, 0)
	for k := 0; k < len(i.Options); k++ {
		options = append(options, OptionFromEntity(i.Options[k]))
	}

	return Category{
		Id:                 i.Id,
		Name:               i.Name,
		Description:        i.Description,
		Icon:               i.Icon,
		Parents:            i.Parents,
		CountriesWhitelist: i.CountryWhitelist,
		CountriesBlacklist: i.CountryBlacklist,
		Options:            options,
		Hidden:             i.Hidden,
	}
}

func (controller Controller) GetCategories(w http.ResponseWriter, r *http.Request) {
	types, err := controller.interactor.GetCategorys()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	var serTypes []Category = make([]Category, 0)
	for i := 0; i < len(types); i++ {
		serTypes = append(serTypes, CategoryFromEntity(types[i]))
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    serTypes,
	}, http.StatusOK)
}

func (controller Controller) GetAddCategory(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Name               string      `json:"name"`
		Description        string      `json:"description"`
		Icon               string      `json:"icon"`
		Parents            []uuid.UUID `json:"parents"`
		CountriesWhitelist []string    `json:"country_whitelist"`
		CountriesBlacklist []string    `json:"country_blacklist"`
		Options            []struct {
			Name             string        `json:"name"`
			Description      string        `json:"description"`
			DataType         string        `json:"data_type"`
			RepresentedIn    string        `json:"represented_in"`
			Values           []interface{} `json:"values"`
			AllowCustomValue bool          `json:"allow_custom_value"`
			Validator        map[string]struct {
				Value   interface{} `json:"value"`
				Message string      `json:"message"`
			} `json:"validator"`
		} `json:"options"`
		Hidden bool `json:"hidden"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var options []struct {
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
	} = make([]struct {
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
	}, 0)

	for i := 0; i < len(req.Options); i++ {

		var validator map[string]struct {
			Value   interface{}
			Message string
		} = make(map[string]struct {
			Value   interface{}
			Message string
		})

		for k, v := range req.Options[i].Validator {
			validator[k] = struct {
				Value   interface{}
				Message string
			}{
				Value:   v.Value,
				Message: v.Message,
			}
		}

		options = append(options, struct {
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
		}{
			Name:             req.Options[i].Name,
			Description:      req.Options[i].Description,
			DataType:         entity.DataType(req.Options[i].DataType),
			RepresentedIn:    req.Options[i].RepresentedIn,
			Values:           req.Options[i].Values,
			AllowCustomValue: req.Options[i].AllowCustomValue,
			Validator:        validator,
		})
	}

	category, err := controller.interactor.CreateCategory(
		req.Name,
		req.Description,
		req.Icon,
		req.Parents,
		req.CountriesWhitelist,
		req.CountriesBlacklist,
		req.Hidden,
		options,
	)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	SendJSONResponse(w, Response{
		Success: true,
		Data:    CategoryFromEntity(*category),
	}, http.StatusOK)

	// SendJSONResponse(w, Response{
	// 	Success: true,
	// 	Data:    CategoryFromEntity(*_type),
	// }, http.StatusOK)
}
