package rest

import (
	"auth/src/pkg/account/core/entity"
	"auth/src/pkg/account/usecase"
	auth "auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Account struct {
	Id                 string `json:"id"`
	Title              string `json:"title"`
	Type               string `json:"type"`
	Default            bool   `json:"is_default"`
	VerificationStatus struct {
		Verified   bool `json:"verified"`
		VerifiedBy *struct {
			Method  string      `json:"type"`
			Details interface{} `json:"details"`
		} `json:"verified_by"`
	} `json:"verification_status"`
	Detail interface{} `json:"detail"`
}

func NewAccountFromEntity(i entity.Account) Account {
	var accVerificationStat interface{}
	var acc Account
	if i.VerificationStatus.VerifiedBy != nil {
		switch i.VerificationStatus.VerifiedBy.Method {
		case "SMS":
			{
				accVerificationStat = struct {
					Verified   bool "json:\"verified\""
					VerifiedBy *struct {
						Method  string      "json:\"type\""
						Details interface{} "json:\"details\""
					} "json:\"verified_by\""
				}{
					Verified: i.VerificationStatus.Verified,
					VerifiedBy: &struct {
						Method  string      "json:\"type\""
						Details interface{} "json:\"details\""
					}{
						Method: i.VerificationStatus.VerifiedBy.Method,
						Details: struct {
							Lenth   int `json:"length"`
							Timeout int `json:"timeout"`
						}{
							Lenth: i.VerificationStatus.VerifiedBy.Details.(struct {
								Length  int
								Timeout int
							}).Length,
							Timeout: i.VerificationStatus.VerifiedBy.Details.(struct {
								Length  int
								Timeout int
							}).Timeout,
						},
					},
				}
			}
		}
	}

	if i.VerificationStatus.Verified {
		accVerificationStat = struct {
			Verified   bool "json:\"verified\""
			VerifiedBy *struct {
				Method  string      "json:\"type\""
				Details interface{} "json:\"details\""
			} "json:\"verified_by\""
		}{
			Verified: true,
		}
	}

	acc = Account{
		Id:      i.Id.String(),
		Title:   i.Title,
		Type:    string(i.Type),
		Default: i.Default,
	}
	if accVerificationStat != nil {
		acc = Account{
			Id:      i.Id.String(),
			Title:   i.Title,
			Type:    string(i.Type),
			Default: i.Default,
			VerificationStatus: accVerificationStat.(struct {
				Verified   bool "json:\"verified\""
				VerifiedBy *struct {
					Method  string      "json:\"type\""
					Details interface{} "json:\"details\""
				} "json:\"verified_by\""
			}),
		}
	}

	switch i.Type {
	case entity.STORED:
		{
			acc.Detail = NewStoredAccountFromEntity(i.Detail.(entity.StoredAccount))
		}
	case entity.BANK:
		{
			acc.Detail = NewBankAccountFromEntity(i.Detail.(entity.BankAccount))
		}
	}

	return acc
}

type BankAccount struct {
	Bank          Bank   `json:"bank"`
	AccountNumber string `json:"account_number"`
	Holder        struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	} `json:"holder"`
	Verified bool `json:"verified"`
}

type StoredAccount struct {
	Balance float64 `json:"balance"`
}

func NewBankAccountFromEntity(i entity.BankAccount) BankAccount {
	return BankAccount{
		Bank:          NewBankFromEntity(i.Bank),
		AccountNumber: i.Number,
		Holder: struct {
			Name  string "json:\"name\""
			Phone string "json:\"phone\""
		}{
			Name:  i.Holder.Name,
			Phone: i.Holder.Phone,
		},
	}
}

func NewStoredAccountFromEntity(i entity.StoredAccount) StoredAccount {
	return StoredAccount{
		Balance: i.Balance,
	}
}

// Get all accounts of a user
func (controller Controller) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	// Parse req
	type Request struct {
		Token string
	}

	var req Request

	token := strings.Split(r.Header.Get("Authorization"), " ")

	if len(token) == 2 {
		req.Token = token[1]
	}

	// Authenticate user
	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)
	controller.log.Println("PASSED 0")
	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	// [USECASE]
	accs, err := controller.interactor.GetUserAccounts(session.User.Id)
	if err != nil {
		var status int
		if err.(usecase.Error).Type == "UNAUTHORIZED" {
			status = http.StatusUnauthorized
		} else {
			status = http.StatusBadRequest
		}
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, status)
		return
	}

	var accsRes []Account = make([]Account, 0)

	for i := 0; i < len(accs); i++ {
		accsRes = append(accsRes, NewAccountFromEntity(accs[i]))
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    accsRes,
	}, http.StatusOK)
}

func (controller Controller) GetAddBankAccount(w http.ResponseWriter, r *http.Request) {

	// Authenticate request
	var token string

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide a valid header token",
			},
		}, http.StatusUnauthorized)
		return
	}

	token = strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}

	// [TODO] Authorize request

	// Parse request

	type Request struct {
		Title         string    `json:"title"`
		MakeDefault   bool      `json:"make_default"`
		BankId        uuid.UUID `json:"bank_id"`
		AccountNumber string    `json:"account_number"`
		Holder        struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
		} `json:"holder"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "The request contains one or more invalid parameters please refer to the spec",
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	/*

		[Add Bank Account = Steps = ]
		1. Validate account
		- Check if the provided account and the holder information corresponds
		2. Create account
		- Create a bank account and store
		3. Verify account

	*/

	// Usecase
	acc, err := controller.interactor.CreateBankAccount(session.User.Id, req.BankId, req.AccountNumber, req.Holder.Name, req.Holder.Phone, req.Title, req.MakeDefault)
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

	controller.log.Println(acc.VerificationStatus.Verified)

	SendJSONResponse(w, Response{
		Success: true,
		Data:    NewAccountFromEntity(*acc),
	}, http.StatusOK)

}

func (controller Controller) GetVerifyAccount(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	var token string

	controller.log.Println("verify request")

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide a valid header token",
			},
		}, http.StatusUnauthorized)
		return
	}

	token = strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}

	// [TODO] Authorize request

	// Parse request

	type Request struct {
		Method  string `json:"method"`
		Details struct {
			Code string `json:"code"`
		} `json:"details"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "The request contains one or more invalid parameters please refer to the spec",
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	accountId, _ := uuid.Parse(r.URL.Query().Get("id"))

	// Usecase
	token, err = controller.interactor.VerifyAccount(session.User.Id, accountId, req.Method, req.Details, req.Details.Code)
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

	SendJSONResponse(w, Response{
		Success: true,
		Data:    token,
	}, http.StatusOK)
}

func (controller Controller) GetDeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	var token string

	controller.log.Println("verify request")

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide a valid header token",
			},
		}, http.StatusUnauthorized)
		return
	}

	token = strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}

	accId, _ := uuid.Parse(r.URL.Query().Get("id"))

	// Usecase
	err = controller.interactor.DeleteAccount(session.User.Id, accId)
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

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "Account deleted successfully",
	}, http.StatusOK)

}
