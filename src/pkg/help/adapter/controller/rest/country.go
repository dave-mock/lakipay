package rest

import (
	"auth/src/pkg/help/core/entity"
	"encoding/json"
	"net/http"
)

type PhonePrefix struct {
	Prefix  string `json:"prefix"`
	Pattern string `json:"pattern"`
}

type Country struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	DefaultName string      `json:"default_name"`
	Iso2        string      `json:"iso2"`
	Flag        string      `json:"flag"`
	PhonePrefix PhonePrefix `json:"phone_prefix"`
	Hidden      bool        `json:"hidden"`
}

func (controller Controller) GetAddCountry(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Name        string `json:"name"`
		DefaultName string `json:"default_name"`
		Iso2        string `json:"iso2"`
		Flag        string `json:"flag"`
		PhonePrefix struct {
			Prefix  string `json:"prefix"`
			Pattern string `json:"pattern"`
		} `json:"phone_prefix"`
	}

	// Parse Request
	var req Request

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "FAILED_TO_DECODE_REQUEST_BODY",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Call Usecases

	var phonePrefix *entity.PhonePrefix

	// Find / Create Phone Prefix
	phonePrefix, err = controller.interactor.GetPhonePrefixByPrefix(req.PhonePrefix.Prefix)
	if err != nil {
		phonePrefix, err = controller.interactor.CreatePhonePrefix(req.PhonePrefix.Prefix, req.PhonePrefix.Pattern)
		if err != nil {
			SendJSONResponse(w, Response{
				Success: false,
				Error:   err,
			}, http.StatusBadRequest)
			return
		}

	}

	// Create Country
	country, err := controller.interactor.CreateCountry(req.Name, req.DefaultName, req.Iso2, req.Flag, *phonePrefix, true)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: Country{
			Id:          country.Id.String(),
			Name:        country.Name,
			DefaultName: country.DefaultName,
			Iso2:        country.Iso2,
			Flag:        country.Flag,
			PhonePrefix: PhonePrefix{
				Prefix:  phonePrefix.Prefix,
				Pattern: phonePrefix.Pattern,
			},
		},
	}, http.StatusOK)
}

func (controller Controller) GetCountries(w http.ResponseWriter, r *http.Request) {
	// type Request struct{}

	countries, err := controller.interactor.GetCountries()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, http.StatusBadRequest)
		return
	}

	// Parse response
	var parsedCountries []Country = make([]Country, 0)

	for i := 0; i < len(countries); i++ {
		parsedCountries = append(parsedCountries, Country{
			Id:          countries[i].Id.String(),
			Name:        countries[i].Name,
			DefaultName: countries[i].DefaultName,
			Iso2:        countries[i].Iso2,
			Flag:        countries[i].Flag,
			PhonePrefix: PhonePrefix{
				Prefix:  countries[i].PhonePrefix.Prefix,
				Pattern: countries[i].PhonePrefix.Pattern,
			},
			Hidden: countries[i].Hidden,
		})
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    parsedCountries,
	}, http.StatusOK)

}

func (controller Controller) GetSupportedCountries(w http.ResponseWriter, r *http.Request) {
	// type Request struct{}

	countries, err := controller.interactor.GetSupportedCountries()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, http.StatusBadRequest)
		return
	}

	// Parse response
	var parsedCountries []Country = make([]Country, 0)

	for i := 0; i < len(countries); i++ {
		parsedCountries = append(parsedCountries, Country{
			Id:          countries[i].Id.String(),
			Name:        countries[i].Name,
			DefaultName: countries[i].DefaultName,
			Iso2:        countries[i].Iso2,
			Flag:        countries[i].Flag,
			PhonePrefix: PhonePrefix{
				Prefix:  countries[i].PhonePrefix.Prefix,
				Pattern: countries[i].PhonePrefix.Pattern,
			},
			Hidden: countries[i].Hidden,
		})
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    parsedCountries,
	}, http.StatusOK)

}

func (controller Controller) GetSearchAndFilterCountries(w http.ResponseWriter, r *http.Request) {
	// type Request struct{}

	var query = r.URL.Query().Get("search_query")
	var filterKey = r.URL.Query().Get("fk")

	countries, err := controller.interactor.SearchAndFilterCountries(query, filterKey)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, http.StatusBadRequest)
		return
	}

	// Parse response
	var parsedCountries []Country = make([]Country, 0)

	for i := 0; i < len(countries); i++ {
		parsedCountries = append(parsedCountries, Country{
			Id:          countries[i].Id.String(),
			Name:        countries[i].Name,
			DefaultName: countries[i].DefaultName,
			Iso2:        countries[i].Iso2,
			Flag:        countries[i].Flag,
			PhonePrefix: PhonePrefix{
				Prefix:  countries[i].PhonePrefix.Prefix,
				Pattern: countries[i].PhonePrefix.Pattern,
			},
			Hidden: countries[i].Hidden,
		})
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    parsedCountries,
	}, http.StatusOK)

}
