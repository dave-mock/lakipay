package rest

import (
	"auth/src/pkg/account/usecase"
	auth "auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	auth       auth.Controller
	sm         *http.ServeMux
}

func New(log *log.Logger, interactor usecase.Interactor, sm *http.ServeMux, auth auth.Controller) Controller {
	controller := Controller{log: log, interactor: interactor, auth: auth}

	// Handle routing

	// Accounts

	sm.HandleFunc("/user-update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.UpdateUser(w, r)
			}
		}
	})

	sm.HandleFunc("/send-otp", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetSendOtp(w, r)
			}
		}
	})

	sm.HandleFunc("/finger-print", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			{
				controller.GetSetFingerPrint(w, r)
			}
		}
	})
	sm.HandleFunc("/generate-challenge", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetGenerateChallenge(w, r)
			}

		}
	})

	sm.HandleFunc("/verify-signature", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			{
				controller.GetVerifySignatureHandler(w, r)
			}
		}
	})

	sm.HandleFunc("/store-public-key", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetstorePublicKeyHandler(w, r)
			}
		}
	})

	sm.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("||||||||||||||| /accounts")
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetUserAccounts(w, r)
			}
		case http.MethodPost:
			{
				switch r.URL.Query().Get("type") {
				case "bank":
					{
						controller.GetAddBankAccount(w, r)
					}
				}
			}
		case http.MethodDelete:
			{
				controller.GetDeleteAccount(w, r)
			}
		}
	})
	sm.HandleFunc("/accounts/verify", func(w http.ResponseWriter, r *http.Request) {
		controller.log.Println("verify rest")
		switch r.Method {
		case http.MethodPatch:
			{
				controller.GetVerifyAccount(w, r)
			}
		}
	})

	// Banks
	sm.HandleFunc("/accounts/banks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetBanks(w, r)
			}
		case http.MethodPost:
			{
				controller.GetAddBank(w, r)
			}

		}
	})

	// // Transactions
	sm.HandleFunc("/accounts/transactions", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("||||| tsransactions one")
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetUserTransactions(w, r)
			}
		case http.MethodPost:
			{
				controller.GetRequestTransaction(w, r)
			}
		}
	})

	sm.HandleFunc("/accounts/balance", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			{
				controller.GetCheckBalanceApi(w, r)
			}
		}
	})

	sm.HandleFunc("/accounts/apply-for-token", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			{
				controller.GetApplyForToken(w, r)
			}
		}
	})

	sm.HandleFunc("/accounts/transactions-keys", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("||||| transactions one")
		switch r.Method {

		case http.MethodPost:
			{
				controller.GetRegisterKeys(w, r)
			}
		case http.MethodGet:
			{
				controller.GetApiKeys(w, r)
			}
		}
	})

	sm.HandleFunc("/accounts/transactions-hosted-intiate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("||||| transactions one")
		switch r.Method {

		case http.MethodPost:
			{
				controller.GetRequestHostedTransactionInitiate(w, r)
			}

		}
	})

	sm.HandleFunc("/accounts/transactions-hosted", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("||||| transactions one")
		switch r.Method {

		case http.MethodPost:
			{
				// controller.CompleteTransaction(w, r)
				controller.GetVerifyTransactionHosted(w, r)
			}
		case http.MethodGet:
			{
				controller.GetTransactionDetails(w, r)
			}
		}
	})

	sm.HandleFunc("/accounts/transactions-intiate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("||||| transactions one")
		switch r.Method {

		case http.MethodPost:
			{
				controller.GetRequestTransactionInitiate(w, r)
			}
		}
	})
	sm.HandleFunc("/accounts/transactions-intiate-hosted", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("||||| transactions two")
		switch r.Method {

		case http.MethodPost:
			{
				controller.GetRequestTransactionInitiateForHosted(w, r)
			}
		}
	})

	sm.HandleFunc("/accounts/transactions-varify", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("|||||| trnsaaction varify")
		controller.GetVerifyTransaction(w, r)
	})
	sm.HandleFunc("/accounts/transactions-dashboard", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("||||| transactions one")
		switch r.Method {
		case http.MethodGet:
			{
				controller.TransactionsDashboard(w, r)
			}
			// case http.MethodPost:
			// 	{
			// 		controller.GetRequestTransaction(w, r)
			// 	}
		}
	})
	// // Verify Transactions
	// sm.HandleFunc("/accounts/transactions", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		{
	// 			// controller.GetVerifyTransaction(w, r)
	// 		}
	// 	}
	// })

	// TEst
	sm.HandleFunc("/accounts/epg", func(w http.ResponseWriter, r *http.Request) {
		controller.log.Println(r.Method)
		controller.log.Println(io.ReadAll(r.Body))

		w.Write([]byte("EPG"))
	})

	controller.sm = sm

	return controller
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

func SendJSONResponse(w http.ResponseWriter, data Response, status int) {
	serData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(serData)
}

/*


	// // Bank Accounts
	// sm.HandleFunc("/accounts/bank-accounts", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		{
	// 			controller.GetAddBankAccount(w, r)
	// 		}
	// 	}
	// })

	// // Verify account
	// sm.HandleFunc("/accounts/bank-accounts/verify", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPatch:
	// 		{
	// 			controller.GetVerifyBankAccount(w, r)
	// 		}
	// 	}
	// })



*/
