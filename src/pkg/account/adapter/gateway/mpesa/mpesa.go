package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type TransactionType string

const (
	C2B TransactionType = "CustomerToBusiness"
	B2C TransactionType = "BusinessToCustomer"
	B2B TransactionType = "BusinessToBusiness"
)

// USSDPushRequest defines the structure of the USSD Push request.
type USSDPushRequest struct {
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
	MerchantName      string  `json:"MerchantName"`
}

// HandleSTKPushRequest sends an STK Push request to the M-Pesa API.
func HandleSTKPushRequest(req USSDPushRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("Probl marshalling STK push request: %v", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Sending USSD Request: %s", jsonData)

	// Create the HTTP POST request
	request, err := http.NewRequest("POST", "https://apisandbox.safaricom.et/mpesa/stkpush/v1/processrequest", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Prob creating new request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Obtain the access token
	accessToken, err := getAccessToken()
	if err != nil {
		log.Printf("Error getting access token: %v", err)
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Set the request headers
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending STK push request: %v", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Log response status and body
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to send STK push request: %s, Status Code: %d", body, resp.StatusCode)
		return fmt.Errorf("STK push request failed: Status Code: %d, response: %s", resp.StatusCode, body)
	}

	log.Printf("USSD push request successful: Status Code: %d", resp.StatusCode)
	return nil
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
func HandleTransactionStatusRequest(req TransactionStatusRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error marshalling STK push request: %v", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Sending Transaction Status Request: %s", jsonData)

	request, err := http.NewRequest("POST", "https://apisandbox.safaricom.et/mpesa/transactionstatus/v1/query", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Prob creating new request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Obtain the access token
	accessToken, err := getAccessToken()
	if err != nil {
		log.Printf("Error getting access token: %v", err)
		return fmt.Errorf("failed to get access token: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending STK push request: %v", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Log response status and body
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to send STK push request: %s, Status Code: %d", body, resp.StatusCode)
		return fmt.Errorf("STK push request failed: Status Code: %d, response: %s", resp.StatusCode, body)
	}

	log.Printf("USSD push request successful: Status Code: %d", resp.StatusCode)
	return nil
}

// GetAccessToken obtains an access token from the Safaricom API.
func getAccessToken() (string, error) {
	username := os.Getenv("SAFARICOM_USERNAME")
	password := os.Getenv("SAFARICOM_PASSWORD")
	grantType := "client_credentials"

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	url := "https://apisandbox.safaricom.et/v1/token/generate?grant_type=" + grantType

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Response body: %s", body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s, response: %s", resp.Status, body)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   string `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return tokenResp.AccessToken, nil
}
