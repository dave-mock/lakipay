package rest

import (
	"auth/src/pkg/help/core/entity"
	"auth/src/pkg/help/usecase"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CurrencyCountry struct {
	Id          uuid.UUID `json:"id"`
	Default     bool      `json:"default"`
	Name        string    `json:"name"`
	DefaultName string    `json:"default_name"`
	Iso2        string    `json:"iso2"`
	Flag        string    `json:"flag"`
}

func NewCurrencyCountryFromEntity(country entity.CurrencyCountry) CurrencyCountry {
	return CurrencyCountry{
		Id:          country.Id,
		Default:     country.Default,
		Name:        country.Name,
		DefaultName: country.DefaultName,
		Iso2:        country.Iso2,
		Flag:        country.Flag,
	}
}

type Currency struct {
	Id uuid.UUID `json:"id"`
	// Name - dollar
	Name string `json:"name"`
	// Currency - USD
	Currency string `json:"currency"`
	// Symbol - $
	Symbol    string            `json:"symbol"`
	Rate      float64           `json:"rate"`
	BaseId    uuid.UUID         `json:"base_id"`
	Hidden    bool              `json:"hidden"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Countries []CurrencyCountry `json:"countries"`
}

func NewCurrencyFromEntity(currency entity.Currency) Currency {

	currCountries := make([]CurrencyCountry, 0)
	for i := 0; i < len(currency.Countries); i++ {
		currCountries = append(currCountries, NewCurrencyCountryFromEntity(currency.Countries[i]))
	}

	return Currency{
		Id:        currency.Id,
		Name:      currency.Name,
		Currency:  currency.Currency,
		Symbol:    currency.Symbol,
		Rate:      currency.Rate,
		BaseId:    currency.BaseId,
		Hidden:    currency.Hidden,
		CreatedAt: currency.CreatedAt,
		UpdatedAt: currency.UpdatedAt,
		Countries: currCountries,
	}
}

func (controller Controller) GetAddCurrency(w http.ResponseWriter, r *http.Request) {

	// Authorize

	// Parse Request
	type Request struct {
		Name      string    `json:"name"`
		Currency  string    `json:"currency"`
		Symbol    string    `json:"symbol"`
		Rate      float64   `json:"rate"`
		BaseId    uuid.UUID `json:"base_id"`
		Countries []struct {
			Id      uuid.UUID `json:"id"`
			Default bool      `json:"default"`
		} `json:"countries"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALId_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return

	}

	defer r.Body.Close()

	var countries []entity.CurrencyCountry = make([]entity.CurrencyCountry, 0)

	if req.Countries != nil && len(req.Countries) > 0 {
		for i := 0; i < len(req.Countries); i++ {
			countries = append(countries, entity.CurrencyCountry{
				Id:      req.Countries[i].Id,
				Default: req.Countries[i].Default,
			})
		}
	}

	// Usecase
	curr, err := controller.interactor.CreateCurrency(req.Name, req.Currency, req.Symbol, req.Rate, req.BaseId, countries)
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

	// Response
	SendJSONResponse(w, Response{
		Success: true,
		Data:    NewCurrencyFromEntity(*curr),
	}, http.StatusOK)
}

// Get Currencies
func (controller Controller) GetCurrencies(w http.ResponseWriter, r *http.Request) {

	// Authorize

	// Parse request

	// Usecase Operation
	currs, err := controller.interactor.GetCurrencies()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(*usecase.Error).Type,
				Message: err.(*usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	var serCurrs []Currency = make([]Currency, 0)
	for i := 0; i < len(currs); i++ {
		serCurrs = append(serCurrs, NewCurrencyFromEntity(currs[i]))
	}

	// Response
	SendJSONResponse(w, Response{
		Success: true,
		Data:    serCurrs,
	}, http.StatusOK)
}
