package rest

import (
	"auth/src/pkg/org/core/entity"
	"auth/src/pkg/org/usecase"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type LegalCondition struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	CountryWhitelist []string  `json:"country_whitelist"`
	CountryBlacklist []string  `json:"country_blacklist"`
}

func NewLegalConditionFromEntity(i entity.LegalCondition) LegalCondition {
	return LegalCondition{
		Id:               i.Id,
		Name:             i.Name,
		CountryWhitelist: i.CountryWhitelist,
		CountryBlacklist: i.CountryBlacklist,
	}
}

func (controller Controller) GetAddLegalCondition(w http.ResponseWriter, r *http.Request) {
	// Authenticate

	// Authorize

	// Parse request
	type Request struct {
		Name             string   `json:"name"`
		CountryWhitelist []string `json:"country_whitelist"`
		CountryBlacklist []string `json:"country_blacklist"`
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

	// Usecase
	legalCondition, err := controller.interactor.CreateLegalCondition(req.Name, req.CountryWhitelist, req.CountryBlacklist)
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

	// Send Response
	SendJSONResponse(w, Response{
		Success: true,
		Data:    NewLegalConditionFromEntity(*legalCondition),
	}, http.StatusOK)
}

func (controller Controller) GetLegalConditions(w http.ResponseWriter, r *http.Request) {
	legalconditions, err := controller.interactor.GetLegalConditions()
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

	var serLegalConditions []LegalCondition = make([]LegalCondition, 0)

	for i := 0; i < len(legalconditions); i++ {
		serLegalConditions = append(serLegalConditions, NewLegalConditionFromEntity(legalconditions[i]))
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    serLegalConditions,
	}, http.StatusOK)
}
