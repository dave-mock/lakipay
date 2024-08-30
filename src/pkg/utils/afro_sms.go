package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SMSRequestBody struct {
	Sender  string      `json:"sender"`
	To      string      `json:"to"`
	Message interface{} `json:"message"`
}

type SMSResponse struct {
	Acknowledge string `json:"acknowledge"`
}

// func getEnv(key string) string {
// 	return os.Getenv(key)
// }

// func encryptWithPrivateKey(message []byte, privKey *rsa.PrivateKey) ([]byte, error) {
// 	// Encrypt using RSA private key
// 	ciphertext, err := rsa.SignPKCS1v15(rand.Reader, privKey, 0, message)
// 	if err != nil {
// 		return nil, fmt.Errorf("encryption error: %v", err)
// 	}
// 	return ciphertext, nil
// }

// func decryptWithPublicKey(ciphertext []byte, pubKey *rsa.PublicKey) ([]byte, error) {
//     // Decrypt using RSA public key
//     plaintext, err := rsa.VerifyPKCS1v15(pubKey, 0, ciphertext, nil)
//     if err != nil {
//         return nil, fmt.Errorf("decryption error: %v", err)
//     }
//     return plaintext, nil
// }

func SendSMS(message interface{}, phone string) {
	baseURL := "https://api.afromessage.com/api/send"
	token := "eyJhbGciOiJIUzI1NiJ9.eyJpZGVudGlmaWVyIjoiV2dKVzRBUmtDM3ZNeE5xM3VwWTNkS2Y4V3hSRjhiam4iLCJleHAiOjE4NjM4MDM2MjIsImlhdCI6MTcwNTk1MDgyMiwianRpIjoiODNkYWQxNWUtODU4Ny00NTU0LTkxZWItZDYxYTU0NDNjYTkzIn0.xTzHDwpU9qkRurb0iCnixPFs3qanu3pk86L3hJTXZEI"

	client := &http.Client{}

	body := SMSRequestBody{
		To:      phone,
		Message: message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}

		var jsonResponse SMSResponse
		if err := json.Unmarshal(bodyBytes, &jsonResponse); err != nil {
			fmt.Printf("Error unmarshalling JSON: %v\n", err)
			return
		}

		fmt.Println(string(bodyBytes))
		if jsonResponse.Acknowledge == "success" {
			fmt.Println("API success")
		} else {
			fmt.Println("API error")
		}
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("HTTP error ... code: %d , msg: %s\n", resp.StatusCode, string(bodyBytes))
	}
}

// func main() {
// 	message := "Your message here"
// 	phone := "1234567890"
// 	sendSMS(message, phone)
// }
