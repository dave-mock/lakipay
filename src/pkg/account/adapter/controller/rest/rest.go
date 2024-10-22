package rest

import (
	"auth/src/pkg/account/usecase"
	auth "auth/src/pkg/auth/adapter/controller/procedure"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Controller struct
type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	auth       auth.Controller
	sm         *http.ServeMux
}

// New function initializes the Controller and sets up routes
func New(log *log.Logger, interactor usecase.Interactor, sm *http.ServeMux, auth auth.Controller) Controller {
	controller := Controller{log: log, interactor: interactor, auth: auth}

	// Handle routing
	// User Update
	sm.HandleFunc("/user-update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.UpdateUser(w, r)
		}
	})

	// Send OTP
	sm.HandleFunc("/send-otp", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetSendOtp(w, r)
		}
	})

	// Finger Print
	sm.HandleFunc("/finger-print", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			controller.GetSetFingerPrint(w, r)
		}
	})

	// Generate Challenge
	sm.HandleFunc("/generate-challenge", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetGenerateChallenge(w, r)
		}
	})

	// Verify Signature
	sm.HandleFunc("/verify-signature", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetVerifySignatureHandler(w, r)
		}
	})

	// Store Public Key
	sm.HandleFunc("/store-public-key", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetstorePublicKeyHandler(w, r)
		}
	})

	// Accounts
	sm.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetUserAccounts(w, r)
		case http.MethodPost:
			switch r.URL.Query().Get("type") {
			case "bank":
				controller.GetAddBankAccount(w, r)
			}
		case http.MethodDelete:
			controller.GetDeleteAccount(w, r)
		}
	})

	sm.HandleFunc("/accounts/verify", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			controller.GetVerifyAccount(w, r)
		}
	})

	// Banks
	sm.HandleFunc("/accounts/banks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetBanks(w, r)
		case http.MethodPost:
			controller.GetAddBank(w, r)
		}
	})

	// Transactions
	sm.HandleFunc("/accounts/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetUserTransactions(w, r)
		case http.MethodPost:
			controller.GetRequestTransaction(w, r)
		}
	})

	sm.HandleFunc("/accounts/balance", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetCheckBalanceApi(w, r)
		}
	})

	sm.HandleFunc("/accounts/apply-for-token", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetApplyForToken(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-keys", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetRegisterKeys(w, r)
		case http.MethodGet:
			controller.GetApiKeys(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-hosted-intiate", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetRequestHostedTransactionInitiate(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-hosted", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetVerifyTransactionHosted(w, r)
		case http.MethodGet:
			controller.GetTransactionDetails(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-intiate", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetRequestTransactionInitiate(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-intiate-hosted", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetRequestTransactionInitiateForHosted(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-varify", func(w http.ResponseWriter, r *http.Request) {
		controller.GetVerifyTransaction(w, r)
	})

	sm.HandleFunc("/accounts/transactions-dashboard", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.TransactionsDashboard(w, r)
		}
	})

	// Test endpoint
	sm.HandleFunc("/accounts/epg", func(w http.ResponseWriter, r *http.Request) {
		controller.log.Println(r.Method)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		controller.log.Println(string(body))
		w.Write([]byte("EPG"))
	})

	// Endpoint for USSD Push requests
	sm.HandleFunc("/mpesa/ussd-push", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var requestBody struct {
				Amount           float64 `json:"Amount"`
				PartyA           string  `json:"PartyA"`
				PartyB           string  `json:"PartyB"`
				PhoneNumber      string  `json:"PhoneNumber"`
				CallBackURL      string  `json:"CallBackURL"`
				AccountReference string  `json:"AccountReference"`
				TransactionDesc  string  `json:"TransactionDesc"`
				OrderID          string  `json:"OrderID"`
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				log.Printf("Error reading request body: %v", err)
				return
			}
			if err := json.Unmarshal(body, &requestBody); err != nil {
				http.Error(w, "Error parsing JSON", http.StatusBadRequest)
				log.Printf("Error parsing JSON: %v", err)
				return
			}
			stkPushRequest := StkPushRequest{
				BusinessShortCode: "20017",                                                            // Replace with your actual business short code
				Password:          "7fc641d11b7e5037f486741802a66db71de96fab6feba56a4ddb9cf41b912b9f", // Replace with actual password
				Timestamp:         time.Now().Format("20060102150405"),
				TransactionType:   "CustomerPayBillOnline",
				Amount:            requestBody.Amount,
				PartyA:            requestBody.PartyA,
				PartyB:            requestBody.PartyB,
				PhoneNumber:       requestBody.PhoneNumber,
				CallBackURL:       requestBody.CallBackURL,
				AccountReference:  requestBody.AccountReference,
				TransactionDesc:   requestBody.TransactionDesc,
				OrderID:           requestBody.OrderID,
			}
			go handleSTKPushRequest(stkPushRequest)
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, "Payment request accepted")
		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})
	// Endpoint for LakiPay Transaction Status Query
	sm.HandleFunc("/mpesa/transaction-status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var requestBody struct {
				TransactionID   string `json:"TransactionID"`
				PartyA          string `json:"PartyA"`
				ResultURL       string `json:"ResultURL"`
				QueueTimeOutURL string `json:"QueueTimeOutURL"`
				Remarks         string `json:"Remarks"`
				Occasion        string `json:"Occasion"`
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				log.Printf("Error reading request body: %v", err)
				return
			}
			if err := json.Unmarshal(body, &requestBody); err != nil {
				http.Error(w, "Error parsing JSON", http.StatusBadRequest)
				log.Printf("Error parsing JSON: %v", err)
				return
			}

			statusRequest := TransactionStatusRequest{
				Initiator:          "apitest",
				SecurityCredential: "K2tH8KzkmIjCSEOAwgiI6T4ThpQT2DQa/rZfRc8+6iYA62vIsCUhhV0yxM84p0O/70aAiJ6EZwr/bE17Ww1VPMpzPwBadgf+dTwdz8LueZi4kUyZleoIYiJ3jNjkaXMT2g29JHJKRePbd4fsk+y38avf9zl5HG1N9UrpjaT2LhZOjmmPj4U1P/l3NP9C+AlcbAtHmtrkyIhvPrH4XwtDZasRcxUpPMbVcvajaBrcixIga8I4bvfBJfsnspSmpSDoTZ1f9glMYy1qRD83NX6R4/ToD0K0n7z0tKo35rJIn6a1bpVqKXuPmrkcK9ck0nAtmPQy8om6pwCmzDu+sG6slg==",
				CommandID:          "TransactionStatusQuery",
				TransactionID:      requestBody.TransactionID,
				PartyA:             requestBody.PartyA,
				IdentifierType:     "4",
				ResultURL:          requestBody.ResultURL,
				QueueTimeOutURL:    requestBody.QueueTimeOutURL,
				Remarks:            requestBody.Remarks,
				Occasion:           requestBody.Occasion,
			}

			go handleTransactionStatusRequest(statusRequest, w)
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, "Transaction status request accepted")
		} else {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	return controller
}

// USSDPushRequest struct
type StkPushRequest struct {
	BusinessShortCode string  `json:"BusinessShortCode"`
	Password          string  `json:"Password"`
	Timestamp         string  `json:"Timestamp"`
	TransactionType   string  `json:"TransactionType"`
	Amount            float64 `json:"Amount"`
	PartyA            string  `json:"PartyA"`
	PartyB            string  `json:"PartyB"`
	PhoneNumber       string  `json:"PhoneNumber"`
	CallBackURL       string  `json:"CallBackURL"`
	AccountReference  string  `json:"AccountReference"`
	TransactionDesc   string  `json:"TransactionDesc"`
	OrderID           string  `json:"OrderID"`
}

// handleSTKPushRequest function to send USSD push request
func handleSTKPushRequest(req StkPushRequest) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error marshalling STK push request: %v", err)
		return
	}

	// Log the JSON payload for debugging
	log.Printf("Sending USSD Request: %s", jsonData)
	request, err := http.NewRequest("POST", "https://apisandbox.safaricom.et/mpesa/stkpush/v1/processrequest", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		return
	}
	accessToken, err := getAccessToken()
	if err != nil {
		log.Printf("Error getting access token: %v", err)
		return
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending STK push request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Log response status and body
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to send STK push request: %s, Status Code: %d", body, resp.StatusCode)
	} else {
		log.Printf("USSD push request successful: Status Code: %d", resp.StatusCode)
	}
}

// Structure for Transaction Status Request
type TransactionStatusRequest struct {
	Initiator          string `json:"Initiator"`
	SecurityCredential string `json:"SecurityCredential"`
	CommandID          string `json:"CommandID"`
	TransactionID      string `json:"TransactionID"`
	PartyA             string `json:"PartyA"`
	IdentifierType     string `json:"IdentifierType"`
	ResultURL          string `json:"ResultURL"`
	QueueTimeOutURL    string `json:"QueueTimeOutURL"`
	Remarks            string `json:"Remarks"`
	Occasion           string `json:"Occasion"`
}

// Function to send status  request and handle the response
func handleTransactionStatusRequest(req TransactionStatusRequest, w http.ResponseWriter) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error marshalling transaction status request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Sending Transaction Status Request: %s", jsonData)

	request, err := http.NewRequest("POST", "https://apisandbox.safaricom.et/mpesa/transactionstatus/v1/query", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Problem creating new request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	accessToken, err := getAccessToken()
	if err != nil {
		log.Printf("Problem getting access token: %v", err)
		http.Error(w, "Failed to retrieve access token", http.StatusInternalServerError)
		return
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending transaction status request: %v", err)
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send transaction status request: %s, Status Code: %d", responseBody, resp.StatusCode)
		http.Error(w, fmt.Sprintf("Error: %s", responseBody), resp.StatusCode)
	} else {
		log.Printf("Transaction status request successful: %s", responseBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
	}
}

// This obtains an access token from the Safaricom API
func getAccessToken() (string, error) {
	username := os.Getenv("SAFARICOM_USERNAME")
	password := os.Getenv("SAFARICOM_PASSWORD")
	grantType := "client_credentials"

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	url := "https://apisandbox.safaricom.et/v1/token/generate?grant_type=" + grantType

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Log the response body for debugging
	log.Printf("Response body: %s", body)

	// Check for a non-200 status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s, response: %s", resp.Status, body)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   string `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return tokenResp.AccessToken, nil
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
