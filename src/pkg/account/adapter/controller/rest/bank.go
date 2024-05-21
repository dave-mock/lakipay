package rest

import (
	"auth/src/pkg/account/core/entity"
	"auth/src/pkg/account/usecase"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Bank struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ShortName string    `json:"short_name"`
	BIN       string    `json:"bin"`
	SwiftCode string    `json:"swift_code"`
	Logo      string    `json:"logo"`
}

func NewBankFromEntity(i entity.Bank) Bank {
	return Bank{
		Id:        i.Id,
		Name:      i.Name,
		ShortName: i.ShortName,
		BIN:       i.BIN,
		SwiftCode: i.SwiftCode,
		Logo:      i.Logo,
	}
}

func (controller Controller) GetAddBank(w http.ResponseWriter, r *http.Request) {

	// Req
	type Request struct {
		Name      string `json:"name"`
		ShortName string `json:"short_name"`
		BIN       string `json:"bin"`
		SwiftCode string `json:"swift_code"`
		Logo      string `json:"logo"`
	}

	// Parse request
	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// [UC]
	bank, err := controller.interactor.AddBank(req.Name, req.ShortName, req.BIN, req.SwiftCode, req.Logo)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Send success response
	SendJSONResponse(w, Response{
		Success: true,
		Data: Bank{
			Id:        bank.Id,
			Name:      bank.Name,
			ShortName: bank.ShortName,
			BIN:       bank.BIN,
			SwiftCode: bank.SwiftCode,
			Logo:      bank.Logo,
		},
	}, http.StatusOK)
}

func (controller Controller) GetBanks(w http.ResponseWriter, r *http.Request) {

	banks, err := controller.interactor.GetBanks()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Println(len(banks))

	var _banks []Bank = make([]Bank, 0)

	for i := 0; i < len(banks); i++ {
		_banks = append(_banks, NewBankFromEntity(banks[i]))
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    _banks,
	}, http.StatusOK)
}
