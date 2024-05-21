package rest

import (
	"auth/src/pkg/org/core/entity"
	"auth/src/pkg/org/usecase"
	"encoding/json"
	"net/http"
)

type Tax struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Rate             float64  `json:"rate"`
	From             string   `json:"from"`
	CountryWhitelist []string `json:"countryWhitelist"`
	CountryBlacklist []string `json:"countryBlacklist"`
}

func TaxFromentity(v entity.Tax) Tax {
	return Tax{
		Id:               v.Id.String(),
		Name:             v.Name,
		Description:      v.Description,
		Rate:             v.Rate,
		From:             string(v.From),
		CountryWhitelist: v.CountryWhitelist,
		CountryBlacklist: v.CountryBlacklist,
	}
}

func (controller Controller) GetAddTax(w http.ResponseWriter, r *http.Request) {
	// AuthN

	// AuthZ

	// Parse Req
	type Request struct {
		Name             string   `json:"name"`
		Description      string   `json:"description"`
		Rate             float64  `json:"rate"`
		From             string   `json:"from"`
		CountryWhitelist []string `json:"country_whitelist"`
		CountryBlacklist []string `json:"country_blacklist"`
		Hidden           bool     `json:"hidden"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	defer r.Body.Close()

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

	// Usecase
	tax, err := controller.interactor.CreateTax(
		req.Name,
		req.Description,
		req.Rate,
		entity.TaxableEntity(req.From),
		req.CountryWhitelist,
		req.CountryBlacklist,
		req.Hidden,
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

	SendJSONResponse(w, Response{
		Success: true,
		Data:    TaxFromentity(*tax),
	}, http.StatusOK)

}

// Get Stored Taxes
func (controller Controller) GetTaxes(w http.ResponseWriter, r *http.Request) {
	// AuthN

	// AuthZ

	// Parse

	// Usecase
	taxes, err := controller.interactor.GetTaxes()
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

	var serTaxes []Tax = make([]Tax, 0)

	for i := 0; i < len(taxes); i++ {
		serTaxes = append(serTaxes, TaxFromentity(taxes[i]))
	}

	// Return
	SendJSONResponse(w, Response{
		Success: true,
		Data:    serTaxes,
	}, http.StatusOK)
}
