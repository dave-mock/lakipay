package rest

import (
	"auth/src/pkg/auth/infra/storage/psql"
	"auth/src/pkg/utils"
	"database/sql"
	"log"
	"os"

	"auth/src/pkg/account/core/entity"
	"auth/src/pkg/account/usecase"
	"auth/src/pkg/auth/adapter/controller/procedure"
	auth "auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id     uuid.UUID `json:"id"`
	From   Account   `json:"from"`
	To     Account   `json:"to"`
	Amount float64   `json:"amount"`
	Type   string    `json:"type"`
	Date   time.Time `json:"date"`
}

func NewTransactionFromEntity(i entity.Transaction) Transaction {
	return Transaction{
		Id:     i.Id,
		From:   NewAccountFromEntity(i.From),
		To:     NewAccountFromEntity(i.To),
		Type:   string(i.Type),
		Date:   i.CreatedAt,
		Amount: i.Amount,
	}
}

func (controller Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user entity.User2
	decoder := json.NewDecoder((r.Body))
	err := decoder.Decode(&user)

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

	users, err := controller.interactor.UpdateUserUsecase(user)
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
	SendJSONResponse(w, Response{
		Success: true,
		Data:    users,
	}, http.StatusOK)

}

func (controller Controller) GetVerifyTransactionHosted(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	log := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	db, err := psql.New(log)
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
	// Parse request
	type Request struct {
		Phone         string `json:"phone"`
		Token         string `json:"token"`
		TwoFA         string `json:"2fa"`
		Challenge     string `json:"challenge"`
		Signature     string `json:"signature"`
		OTP           string `json:"otp"`
		ChallengeType string `json:"challenge_type"`
		// Amount      float64   `json:"amount"`
	}

	var req Request
	var transactionChallenge entity.TransactionChallange

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
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

	defer r.Body.Close()

	transactionChallenge.Signature = req.Signature
	transactionChallenge.TwoFA = req.TwoFA
	transactionChallenge.Challenge = req.Challenge
	transactionChallenge.OTP = req.OTP

	var sender_id uuid.UUID
	sqlStmt := `select pi.user_id from auth.phones as p
join auth.phone_identities as pi on p.id = pi.phone_id
WHERE p.number = $1
;`
	err = db.QueryRow(sqlStmt, req.Phone).Scan(&sender_id)

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
	_, err = controller.interactor.VerifyTransaction(sender_id, req.Token, transactionChallenge, req.ChallengeType)

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
		Data:    "success",
	}, http.StatusOK)

	// Usecase
}

func (controller Controller) GetVerifyTransaction(w http.ResponseWriter, r *http.Request) {
	// Authenticate request

	controller.log.Println("Adding Transaction")

	// Check header
	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	// Parse request
	type Request struct {
		Token         string `json:"token"`
		TwoFA         string `json:"2fa"`
		Challenge     string `json:"challenge"`
		Signature     string `json:"signature"`
		OTP           string `json:"otp"`
		ChallengeType string `json:"challenge_type"`
		// Amount      float64   `json:"amount"`
	}

	var req Request
	var transactionChallenge entity.TransactionChallange

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
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

	defer r.Body.Close()

	transactionChallenge.Signature = req.Signature
	transactionChallenge.TwoFA = req.TwoFA
	transactionChallenge.Challenge = req.Challenge
	transactionChallenge.OTP = req.OTP

	_, err = controller.interactor.VerifyTransaction(session.User.Id, req.Token, transactionChallenge, req.ChallengeType)

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
		Data:    "success",
	}, http.StatusOK)

	// Usecase
}

func (controller Controller) GetApiKeys(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	fmt.Println("||| start auth")
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	private_key, err := controller.interactor.GetApiKeysUsecase(session.User.Id)
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

	type Response2 struct {
		PrivateKey string `json:"private_key"`
	}

	var res Response2
	res.PrivateKey = private_key
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) GetApplyForToken(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

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

	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Username and password are required",
			},
		}, http.StatusBadRequest)
		return
	}

	private_key, err := controller.interactor.ApplyForTokenUsecase(req.Username, req.Password)
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

	type Response2 struct {
		Token  string `json:"token"`
		Detail struct {
			Type string `json:"type"`
			Info string `json:"info"`
		} `json:"detail"`
	}

	var res Response2
	res.Token = private_key
	res.Detail.Type = "Bearer Token"
	res.Detail.Info = "Add the Bearer token to the header to authorize."
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)
}

func (controller Controller) GetCheckBalanceApi(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	fmt.Println("||| sesstion id", session.User.Id, session.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}
	type Request struct {
		From uuid.UUID `json:"id"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
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

	defer r.Body.Close()

	// if req.Username == "" || req.Password == "" {
	//     SendJSONResponse(w, Response{
	//         Success: false,
	//         Error: &Error{
	//             Type:    "INVALID_REQUEST",
	//             Message: "Username and password are required",
	//         },
	//     }, http.StatusBadRequest)
	//     return
	// }

	Balance, err := controller.interactor.CheckBalance(req.From)
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

	type Response2 struct {
		Balance float64 `json:"balance"`
	}

	var res Response2
	res.Balance = Balance
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) GetRegisterKeys(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	fmt.Println("||| start auth")
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
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

	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Username and password are required",
			},
		}, http.StatusBadRequest)
		return
	}

	private_key, err := controller.interactor.CreateRegisterKeys(session.User.Id, req.Username, req.Password)
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

	type Response2 struct {
		PrivateKey string `json:"private_key"`
	}

	var res Response2
	res.PrivateKey = private_key
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) CompleteTransaction(w http.ResponseWriter, r *http.Request) {

	transactionID := r.FormValue("transaction_id")
	log := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	// [Ouput Adapters]
	// [DB] Postgres
	db, err := psql.New(log)
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
	type Request struct {
		Phone string `json:"phone"`
	}
	var req Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST222",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var sender_id uuid.UUID
	sqlStmt := `select pi.user_id from auth.phones as p
join auth.phone_identities as pi on p.id = pi.phone_id
WHERE p.number = $1
;`

	err = db.QueryRow(sqlStmt, req.Phone).Scan(&sender_id)
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

	// chek the balance here

	_, err = db.Exec(`UPDATE accounts.transactions SET verified=true, "from"=$1 WHERE id=$2`, sender_id, transactionID)

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

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "Transaction completed successfully",
	}, http.StatusOK)
}
func (controller Controller) GetTransactionDetails(w http.ResponseWriter, r *http.Request) {

	transactionID := r.URL.Query().Get("transaction_id")
	log := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	db, err := psql.New(log)
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

	var transaction entity.Transaction
	var amount sql.NullString

	err = db.QueryRow("SELECT id, amount, currency FROM accounts.transactions WHERE id=$1", transactionID).Scan(&transaction.Id, &amount, &transaction.Currency)
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

	if amount.Valid {
		amount2, err := utils.AesDecription(amount.String)
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
		amount_float, err := strconv.ParseFloat(amount2, 64)
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
		transaction.Amount = amount_float

	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    transaction,
	}, http.StatusOK)
}

// **************************** hosted checkout initate

func (controller Controller) GetRequestHostedTransactionInitiate(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	type Request struct {
		Amount       float64 `json:"amount"`
		Currency     string  `json:"currency"`
		Signature    string  `json:"signature"`
		Callback_url string  `json:"callback_url"`
	}

	type Payment struct {
		Amount       float64 `json:"amount"`
		Currency     string  `json:"currency"`
		Callback_url string  `json:"callback_url"`
	}

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

	defer r.Body.Close()

	payment := Payment{
		Amount:       req.Amount,
		Currency:     req.Currency,
		Callback_url: req.Callback_url,
	}

	jsonOutput, err := json.Marshal(payment)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Convert the JSON byte slice to a string
	jsonString := string(jsonOutput)

	txn, err := controller.interactor.CreateHostedTransactionInitiate(req.Amount, payment.Currency, payment.Callback_url, req.Signature, jsonString, token)

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
		Data:    txn,
	}, http.StatusOK)

}

func (controller Controller) GetRequestTransactionInitiateForHosted(w http.ResponseWriter, r *http.Request) {
	log := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	// [Ouput Adapters]
	// [DB] Postgres
	db, err := psql.New(log)
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

	type Request struct {
		Phone string `json:"phone"`
		Data  string `json:"data"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
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

	defer r.Body.Close()

	var sender_id uuid.UUID
	sqlStmt := `select pi.user_id from auth.phones as p
join auth.phone_identities as pi on p.id = pi.phone_id
WHERE p.number = $1
;`

	err = db.QueryRow(sqlStmt, req.Phone).Scan(&sender_id)

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

	var sender_acount uuid.UUID
	sqlStmt2 := `select a.id from accounts.accounts as a 
where a.user_id= $1`
	err = db.QueryRow(sqlStmt2, sender_id).Scan(&sender_acount)

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

	fmt.Println("GetRequestTransactionInitiateForHosted:there ||||||||||||||||||||||||||||||||||||||||", sender_id)

	decryptedData, err := utils.AesDecription(req.Data)

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

	parts := strings.Split(decryptedData, ",")

	amount, err := strconv.ParseFloat(parts[1], 64)

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

	currency := parts[2]
	url := parts[3]
	merchantIdString := parts[0]

	parsedUUID, err := uuid.Parse(merchantIdString)
	fmt.Println("GetRequestTransactionInitiateForHosted:new ------ ||||||||||||||||||||||||||||||||||||||||", parsedUUID)

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

	var reciver_acount uuid.UUID
	sqlStmt3 := `select a.id from accounts.accounts as a 
	where a.user_id= $1`
	err = db.QueryRow(sqlStmt3, parsedUUID).Scan(&reciver_acount)

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

	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$wwwww", sender_acount)

	txn, err := controller.interactor.CreateTransactionInitiate(sender_id, sender_acount, reciver_acount, amount, "LAKIPAY", "", "", "")
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

	type Response2 struct {
		Result       interface{} `json:"result"`
		Currency     string      `json:"currency"`
		Callback_url string      `json:"callback_url"`
	}

	var res Response2

	res.Callback_url = url
	res.Currency = currency
	res.Result = txn

	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) GetRequestTransactionInitiate(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	fmt.Println("||| start auth")
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}
	type Request struct {
		From   uuid.UUID                `json:"from"`
		To     uuid.UUID                `json:"to"`
		Type   string                   `json:"type"`
		Amount float64                  `json:"amount"`
		Medium entity.TransactionMedium `json:"medium"`
		Detail string                   `json:"details"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
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

	defer r.Body.Close()

	txn, err := controller.interactor.CreateTransactionInitiate(session.User.Id, req.From, req.To, req.Amount, req.Medium, req.Type, token, req.Detail)
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
		Data:    txn,
	}, http.StatusOK)

}

func (controller Controller) GetRequestTransaction(w http.ResponseWriter, r *http.Request) {

	fmt.Println("||||||| GetRequestTransaction")
	controller.log.Println("Adding Transaction")
	// Authenticate (AuthN)

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	fmt.Println("||| start auth")
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	fmt.Println("||| pass auth")

	// Check access (AuthZ)

	// Parse
	type Request struct {
		From          uuid.UUID `json:"from"`
		To            uuid.UUID `json:"to"`
		Type          string    `json:"type"`
		Amount        float64   `json:"amount"`
		Medium        string    `json:"medium"`
		TwoFA         string    `json:"2fa"`
		Challenge     string    `json:"challenge"`
		Signature     string    `json:"signature"`
		OTP           string    `json:"otp"`
		ChallengeType string    `json:"challenge_type"`
	}

	var req Request
	var transactionChallenge entity.TransactionChallange
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
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

	transactionChallenge.Signature = req.Signature
	transactionChallenge.TwoFA = req.TwoFA
	transactionChallenge.Challenge = req.Challenge

	transactionChallenge.OTP = req.OTP

	// decode the transaction challenge
	// decoder2 := json.NewDecoder(r.Body)
	// err = decoder2.Decode(&transactionChallenge)
	// if err != nil {
	// 	SendJSONResponse(w, Response{
	// 		Success: false,
	// 		Error: &Error{
	// 			Type:    "INVALID_REQUEST22",
	// 			Message: err.Error(),
	// 		},
	// 	}, http.StatusBadRequest)
	// 	return
	// }

	defer r.Body.Close()

	// Usecase

	// Response

	// Usecase [CREATE TRANSACTION]
	txn, err := controller.interactor.CreateTransaction(session.User.Id, req.From, req.To, req.Amount, req.Type, token, req.ChallengeType, transactionChallenge)
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
		Data:    NewTransactionFromEntity(*txn),
	}, http.StatusOK)

	// w.Write([]byte("Received Request"))
}

// func (controller Controller) GetTransactions(w http.ResponseWriter, r *http.Request) {
// 	// Request
// 	type Request struct {
// 		Token string
// 	}

// 	req := Request{}

// 	token := r.Header.Get("Authorization")

// 	if len(strings.Split(token, " ")) == 2 {
// 		req.Token = strings.Split(token, " ")[1]
// 	}

// 	// Authenticate user
// 	session, err := controller.auth.GetCheckAuth(req.Token)
// 	if err != nil {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    err.(auth.Error).Type,
// 				Message: err.(auth.Error).Message,
// 			},
// 		}, http.StatusUnauthorized)
// 		return
// 	}

// 	// Get user id
// 	userId := session.User.Id

// 	// Get Transactions
// 	txns, err := controller.interactor.GetUserTransactions(userId)
// 	if err != nil {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    err.(usecase.Error).Type,
// 				Message: err.(usecase.Error).Message,
// 			},
// 		}, http.StatusBadRequest)
// 		return
// 	}

// 	// Map transactions to present
// 	var _txns []Transaction = make([]Transaction, 0)
// 	for i := 0; i < len(txns); i++ {
// 		_txns = append(_txns, NewTransactionFromEntity(txns[i]))
// 	}

// 	SendJSONResponse(w, Response{
// 		Success: true,
// 		Data:    _txns,
// 	}, http.StatusOK)
// }

func (controller Controller) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse req
	type Request struct {
		Token string
	}
	var req Request

	fmt.Println("************************************** one")
	id_string := r.URL.Query().Get("id")
	fmt.Println("************************************** two")

	fmt.Println("************************************** five")

	token := strings.Split(r.Header.Get("Authorization"), " ")

	if len(token) == 2 {
		req.Token = token[1]
	}

	// Authenticate user
	controller.log.Println("PASSED -1")
	_, err := controller.auth.GetCheckAuth(req.Token)
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

	var trx interface{}
	if id_string == "" {
		fmt.Println("||||||||||||||||||| one")
		trx, err = controller.interactor.GetAllTransactions()

	} else {
		fmt.Println("||||||||||||||||||| two")

		id, err := uuid.Parse(id_string)

		if err != nil {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "BAD_REQUEST",
					Message: "bad request",
				},
			}, http.StatusBadRequest)
			return
		}
		trx, err = controller.interactor.GetUserTransactions(id)
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
	}

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

	SendJSONResponse(w, Response{
		Success: true,
		Data:    trx,
	}, http.StatusOK)

}

func (controller Controller) TransactionsDashboard(w http.ResponseWriter, r *http.Request) {
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
	_, err := controller.auth.GetCheckAuth(req.Token)
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

	req_type := r.URL.Query().Get("type")
	year_str := r.URL.Query().Get("year")
	year, err := strconv.Atoi(year_str)
	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "error",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	var data2 interface{}
	switch req_type {
	case "month":
		{
			trx, err := controller.interactor.TransactionsDashboardUsecase(year)
			data2 = trx
			if err != nil {

				SendJSONResponse(w, Response{
					Success: false,
					Error: &Error{
						Type:    "error",
						Message: err.Error(),
					},
				}, http.StatusBadRequest)
				return
			}
		}
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    data2,
	}, http.StatusOK)

}

func (controller Controller) GetSendOtp(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var req Request
	fmt.Println("|||||||||||||||||||||||| one")

	token := strings.Split(r.Header.Get("Authorization"), " ")

	if len(token) == 2 {
		req.Token = token[1]
	}

	fmt.Println("|||||||||||||||||||||||| two", req.Token)

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)
	fmt.Println("|||||||||||||||||||||||| 3")

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

	token2, err := controller.interactor.SendOtpUsecase(session.User.Id)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: false,
		Data:    token2,
	}, http.StatusOK)

}

func (controller Controller) GetSetFingerPrint(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

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

	var data interface{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&data)

	// fmt.Println("||||||||||||||||||||||||||||||||||| ", data)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error()},
		}, http.StatusBadRequest)
		return
	}

	token2, err := controller.interactor.SendSetFIngerPrintUsecase(session.User.Id, data)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: false,
		Data:    token2,
	}, http.StatusOK)

}

func (controller Controller) GetGenerateChallenge(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var requestOne struct {
		// Username  string `json:"username"`
		// Challenge string `json:"challenge"`
		DeviceID string `json:"device_id"`
	}
	var res struct {
		Challenge string `json:"challenge"`
	}
	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

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

	err = json.NewDecoder(r.Body).Decode(&requestOne)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.SendGenerateChallenge(session.User.Id, requestOne.DeviceID)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	res.Challenge = token2
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) GetVerifySignatureHandler(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var requestOne struct {
		// Username  string `json:"username"`
		Challenge string `json:"challenge"`
		Signature string `json:"signature"`
	}

	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

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

	err = json.NewDecoder(r.Body).Decode(&requestOne)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.GetverifySignature(session.User.Id, requestOne.Challenge, requestOne.Signature)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: false,
		Data:    token2,
	}, http.StatusOK)

}

func (controller Controller) GetstorePublicKeyHandler(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var requestOne struct {
		DeviceID  string `json:"device_id"`
		PublicKey string `json:"public_key"`
	}

	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

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

	err = json.NewDecoder(r.Body).Decode(&requestOne)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.GetstorePublicKeyHandler(requestOne.PublicKey, session.User.Id, requestOne.DeviceID)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    token2,
	}, http.StatusOK)

}
