package usecase

import (
	"auth/src/pkg/account/core/entity"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Create transaction
func (uc Usecase) CreateTransaction(userId uuid.UUID, from uuid.UUID, to uuid.UUID, amount float64, txnType string, token string) (*entity.Transaction, error) {
	var txn entity.Transaction

	var sender *entity.Account
	var receipient *entity.Account
	var err error

	// get accounts
	if from != uuid.Nil {
		sender, err = uc.repo.FindAccountById(from)
		if err != nil {
			return nil, Error{
				Type:    "SENDER_ACCOUNT_NOT_FOUND",
				Message: err.Error(),
			}
		}

		if userId != sender.User.Id {
			return nil, Error{
				Type:    "SENDER_ACCOUNT_MISMATCH",
				Message: "Please use your own account",
			}
		}

		// If it's a stored value check balance
		if sender.Type == entity.STORED && sender.Detail.(entity.StoredAccount).Balance < amount {
			return nil, Error{
				Type:    "NOT_ENOUGH_FUND",
				Message: "amount is larger than your balance",
			}
		}
	}

	//
	if to != uuid.Nil {
		receipient, err = uc.repo.FindAccountById(to)
		if err != nil {
			return nil, Error{
				Type:    "RECEPIENT_ACCOUNT_NOT_FOUND",
				Message: err.Error(),
			}
		}
	}

	uc.log.Println(entity.TransactionType(txnType))

	// TXN Type
	switch entity.TransactionType(txnType) {
	case entity.REPLENISHMENT:
		{
			uc.log.Println("A2A")
			// Make A2A transaction
			// if sender.Type == entity.STORED && receipient.Type == entity.STORED {
			// 	// Make transaction no further operations
			// }

			uc.log.Println(sender.Type)

			uc.log.Println(token)

			if sender.Type == entity.BANK && receipient.Type == entity.STORED {
				// Its replenishment
				// Create transaction
				uc.log.Println("Replenishment")
				txn = entity.Transaction{
					Id:        uuid.New(),
					From:      *sender,
					To:        *receipient,
					Type:      entity.TransactionType(txnType),
					Verified:  false,
					CreatedAt: time.Now(),
					Reference: strings.Split(uuid.New().String(), "-")[4],
					Details: entity.Replenishment{
						Amount: amount,
					},
				}
				fmt.Println("||||||||||||||||| [transaction] ", sender.Detail.(entity.BankAccount).Bank.SwiftCode)
				// Validate transaction
				switch sender.Detail.(entity.BankAccount).Bank.SwiftCode {
				case "AWINETAA":
					{
						uc.log.Println("Switching Amhara Bank")
						var netTransport = &http.Transport{
							Dial: (&net.Dialer{
								Timeout: 1 * time.Minute,
							}).Dial,
							TLSHandshakeTimeout: 1 * time.Minute,
						}

						var client = &http.Client{
							Timeout:   time.Minute * 1,
							Transport: netTransport,
						}

						// Authorize client
						serBody, _ := json.Marshal(&struct {
							Username string `json:"username"`
							Password string `json:"password"`
						}{
							Username: "qetxjgflmn",
							Password: "w9'MwO9F$n",
						})

						req, err := http.NewRequest(http.MethodPost, "http://10.10.101.144:8080/b2b/awash/api/v1/auth/getToken", bytes.NewBuffer(serBody))
						if err != nil {
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						req.Header.Set("Content-Type", "application/json")

						res, err := client.Do(req)
						if err != nil {
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						// body, _ := io.ReadAll(res.Body)
						uc.log.Println("Link")
						// uc.log.Println(string(body))
						uc.log.Println("Link")

						if res.StatusCode != http.StatusOK {
							// Unsuccessful request
							uc.log.Println("Send Error")
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: "Failed to link account",
							}
						}

						type AuthRes struct {
							Status bool   `json:"status"`
							Token  string `json:"token"`
						}

						uc.log.Println("Link 1")
						var authRes AuthRes
						decoder := json.NewDecoder(res.Body)
						uc.log.Println("Link 2")
						// uc.log.Println(res.b)
						err = decoder.Decode(&authRes)
						uc.log.Println("Link 3")
						if err != nil {
							uc.log.Println(authRes.Token)
							uc.log.Println(err)
							uc.log.Println("Link 4")
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						defer res.Body.Close()

						// Make transaction
						/*

													{
							    "bankCode": "awsbnk",
							    "amount": 1,
							    "reference": "ASWQERDFTGHY",
							    "narration": "string",
							    "awashAccount": "01320209107500",
							    "creditAccount": "01320206449500",
							    "commisionAmount": 0,
							    "awashAccountName": "string",
							    "creditAccountName": "string"
							}

						*/

						uc.log.Println(txn.Reference)

						serBody, err = json.Marshal(&struct {
							BankCode          string  `json:"bankCode"`
							Amount            float64 `json:"amount"`
							Reference         string  `json:"reference"`
							Narration         string  `json:"narration"`
							AwashAccount      string  `json:"awashAccount"`
							CreditAccount     string  `json:"creditAccount"`
							CommissionAmount  float64 `json:"commisionAmount"`
							AwashAccountName  string  `json:"awashAccountName"`
							CreditAccountName string  `json:"creditAccountName"`
						}{
							BankCode:          "awsbnk",
							Amount:            amount,
							Reference:         txn.Reference,
							Narration:         "",
							AwashAccount:      sender.Detail.(entity.BankAccount).Number,
							CreditAccount:     "01320209107500",
							CommissionAmount:  0,
							AwashAccountName:  sender.Detail.(entity.BankAccount).Holder.Name,
							CreditAccountName: "LakiPay",
						})

						uc.log.Println(txn.Reference)

						uc.log.Println("Amhara Bank 15")
						if err != nil {
							uc.log.Println("Amhara Bank 16")
							uc.log.Println(err)
							return nil, Error{
								Type:    "FAILED_TOVERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 17")
						req, err = http.NewRequest(http.MethodPost, "http://10.10.101.144:8080/b2b/awash/api/v1/monetize/post", bytes.NewBuffer(serBody))
						if err != nil {
							uc.log.Println("Amhara Bank 18")
							return nil, Error{
								Type:    "NO_CLIENT_AUTH_FOUND",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 19")
						req.Header.Set("Content-Type", "application/json")
						req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authRes.Token))

						// Authorize request

						res, err = client.Do(req)
						if err != nil {
							uc.log.Println("Amhara Bank 20")
							return nil, Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 21")

						if res.StatusCode != http.StatusOK {
							body, _ := io.ReadAll(res.Body)
							type TxnRes struct {
								ResponseCode int    `json:"response_code"`
								Status       string `json:"status"`
								Message      string `json:"message"`
							}
							var txnRes TxnRes
							json.Unmarshal(body, &txnRes)
							uc.log.Println(txnRes.Message)
							uc.log.Println("Amhara Bank 22")
							return nil, Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: txnRes.Message,
							}
						}

						type TxnRes struct {
							TransactionStatus  string `json:"TransactionStatus"`
							TransactionAmount  string `json:"TransactionAmount"`
							Status             int    `json:"status"`
							DateProcessed      string `json:"DateProcessed"`
							TransactionDetails string `json:"TransactionDetails"`
						}

						uc.log.Println("Amhara Bank 23")
						// Store transaction
						err = uc.repo.StoreTransaction(txn)
						if err != nil {
							return nil, Error{
								Type:    "FAILED_TO_STORE_TRANSACTION",
								Message: err.Error(),
							}
						}

					}
				case "AMHRETAA":
					{
						uc.log.Println("Amhara Bank")

						var netTransport = &http.Transport{
							Dial: (&net.Dialer{
								Timeout: 1 * time.Minute,
							}).Dial,
							TLSHandshakeTimeout: 1 * time.Minute,
						}

						uc.log.Println("Amhara Bank 1")
						var client = &http.Client{
							Timeout:   time.Minute * 1,
							Transport: netTransport,
						}

						uc.log.Println("Amhara Bank 2")
						// Authorize client
						serBody, _ := json.Marshal(&struct {
							Username string `json:"username"`
							Password string `json:"password"`
						}{
							Username: "LakiPay",
							Password: "e3i1OehzfV0Iz16asdTjZEbYG4F769Vx8Unuo5chkM9V",
						})

						uc.log.Println("Amhara Bank 3")
						req, err := http.NewRequest(http.MethodPost, "http://172.31.2.30:8600/abaApi/v1/lakiPay/authenticate", bytes.NewBuffer(serBody))
						if err != nil {
							uc.log.Println("Amhara Bank 4")
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 5")
						req.Header.Set("Content-Type", "application/json")

						res, err := client.Do(req)
						uc.log.Println("Amhara Bank 6")
						if err != nil {
							uc.log.Println("Amhara Bank 7")
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 8")
						// body, _ := io.ReadAll(res.Body)
						// uc.log.Println(string(body))

						if res.StatusCode != http.StatusOK {
							uc.log.Println("Amhara Bank 9")
							// Unsuccessful request
							uc.log.Println("Send Error")
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: "Failed to link account",
							}
						}

						uc.log.Println("Amhara Bank 10")
						type AuthRes struct {
							ResponseCode int    `json:"response_code"`
							Status       string `json:"status"`
							Message      string `json:"message"`
							Token        string `json:"token"`
						}

						uc.log.Println("Amhara Bank 11")
						var authRes AuthRes
						decoder := json.NewDecoder(res.Body)
						err = decoder.Decode(&authRes)
						if err != nil {
							uc.log.Println("Amhara Bank 12")
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 13")
						defer res.Body.Close()

						uc.log.Println("Amhara Bank 14")
						serBody, err = json.Marshal(&struct {
							Token         string `json:"token"`
							AccountNumber string `json:"account_number"`
							AccountHolder string `json:"account_holder"`
							Merchant      struct {
								AccountNumber string `json:"account_number"`
								AccountHolder string `json:"account_holder"`
							} `json:"merchant"`
							Order struct {
								Id     string  `json:"id"`
								Amount float64 `json:"amount"`
							} `json:"order"`
						}{
							Token:         token,
							AccountNumber: sender.Detail.(entity.BankAccount).Number,
							AccountHolder: sender.Detail.(entity.BankAccount).Holder.Phone,
							Merchant: struct {
								AccountNumber string "json:\"account_number\""
								AccountHolder string "json:\"account_holder\""
							}{
								AccountNumber: "9900000001655",
								AccountHolder: "251942816493",
							},
							Order: struct {
								Id     string  "json:\"id\""
								Amount float64 "json:\"amount\""
							}{
								Id:     txn.Reference,
								Amount: amount,
							},
						})

						uc.log.Println("Amhara Bank 15")
						if err != nil {
							uc.log.Println("Amhara Bank 16")
							uc.log.Println(err)
							return nil, Error{
								Type:    "FAILED_TOVERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 17")
						req, err = http.NewRequest(http.MethodPost, "http://172.31.2.30:8600/abaApi/v1/lakiPay/processPayment", bytes.NewBuffer(serBody))
						if err != nil {
							uc.log.Println("Amhara Bank 18")
							return nil, Error{
								Type:    "NO_CLIENT_AUTH_FOUND",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 19")
						req.Header.Set("Content-Type", "application/json")
						req.Header.Add("Authorization", authRes.Token)

						// Authorize request

						res, err = client.Do(req)
						if err != nil {
							uc.log.Println("Amhara Bank 20")
							return nil, Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 21")

						if res.StatusCode != http.StatusOK {
							body, _ := io.ReadAll(res.Body)
							type TxnRes struct {
								ResponseCode int    `json:"response_code"`
								Status       string `json:"status"`
								Message      string `json:"message"`
							}
							var txnRes TxnRes
							json.Unmarshal(body, &txnRes)
							uc.log.Println(txnRes.Message)
							uc.log.Println("Amhara Bank 22")
							return nil, Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: txnRes.Message,
							}
						}

						uc.log.Println("Amhara Bank 23")
						// Store transaction
						txn.Verified = true
						err = uc.repo.StoreTransaction(txn)
						if err != nil {
							return nil, Error{
								Type:    "FAILED_TO_STORE_TRANSACTION",
								Message: err.Error(),
							}
						}

						// Update transaction / not
						uc.repo.UpdateAccount(entity.Account{
							Id:                 receipient.Id,
							Title:              receipient.Title,
							Type:               receipient.Type,
							Default:            receipient.Default,
							User:               receipient.User,
							VerificationStatus: receipient.VerificationStatus,
							Detail: entity.StoredAccount{
								Balance: receipient.Detail.(entity.StoredAccount).Balance + amount,
							},
						})
					}
				case "ORIRETAA":
					{
						uc.log.Println("Oromia Bank")

						var netTransport = &http.Transport{
							Dial: (&net.Dialer{
								Timeout: 1 * time.Minute,
							}).Dial,
							TLSHandshakeTimeout: 1 * time.Minute,
						}

						var client = &http.Client{
							Timeout:   time.Minute * 1,
							Transport: netTransport,
						}

						serBody, err := json.Marshal(&struct {
							FromAccount     string  `json:"fromAccount"`
							Amount          float64 `json:"amount"`
							Remark          string  `json:"remark"`
							ExplanationCode string  `json:"explanationCode"`
						}{
							FromAccount:     sender.Detail.(entity.BankAccount).Number,
							Amount:          amount,
							ExplanationCode: "9904",
							Remark:          txn.Reference,
						})

						if err != nil {
							uc.log.Println("Oromia Bank 16")
							uc.log.Println(err)
							return nil, Error{
								Type:    "FAILED_TOVERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Oromia Bank 17")
						req, err := http.NewRequest(http.MethodPost, "http://10.10.20.47/fund-transfer/customer-to-settlement", bytes.NewBuffer(serBody))
						if err != nil {
							uc.log.Println("Oromia Bank 18")
							return nil, Error{
								Type:    "NO_CLIENT_AUTH_FOUND",
								Message: err.Error(),
							}
						}

						uc.log.Println("Oromia Bank 19")
						req.Header.Set("Content-Type", "application/json")
						// Authorize request
						req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJPQiIsImp0aSI6ImM1OTg0YTc2YTAyMjA1MDIwNzQ1MTliYThhNWU0OWMzMDk3NTJmMTAyYThhNzhkYjNmNThiM2QxMzAxMzhiMjEiLCJzdWIiOiJsYWtpcGF5IiwiaWF0IjoxNzAxODUwNzc5fQ.sD_C4nwadpgClQADGOPjWjKembyxqCit2tmD_rLsOg7NsFVDv2xbzvnvDnAjD0OKZSfEfhfuKKHsOZfx1crbAA"))

						res, err := client.Do(req)
						if err != nil {
							uc.log.Println("Oromia Bank 20")
							return nil, Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Oromia Bank 21")

						if res.StatusCode != http.StatusOK {
							body, _ := io.ReadAll(res.Body)
							type TxnRes struct {
								ResponseCode int    `json:"response_code"`
								Status       string `json:"status"`
								Message      string `json:"message"`
							}
							var txnRes TxnRes
							json.Unmarshal(body, &txnRes)
							uc.log.Println(txnRes.Message)
							uc.log.Println("Oromia Bank 22")
							return nil, Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: txnRes.Message,
							}
						}

						uc.log.Println("Oromia Bank 23")
						// Store transaction
						txn.Verified = true
						err = uc.repo.StoreTransaction(txn)
						if err != nil {
							return nil, Error{
								Type:    "FAILED_TO_STORE_TRANSACTION",
								Message: err.Error(),
							}
						}

						// Update transaction / not
						uc.repo.UpdateAccount(entity.Account{
							Id:                 receipient.Id,
							Title:              receipient.Title,
							Type:               receipient.Type,
							Default:            receipient.Default,
							User:               receipient.User,
							VerificationStatus: receipient.VerificationStatus,
							Detail: entity.StoredAccount{
								Balance: receipient.Detail.(entity.StoredAccount).Balance + amount,
							},
						})
					}
				case "BUNAETAA":
					{
						uc.log.Println("Switching Bunna Bank")
						var netTransport = &http.Transport{
							Dial: (&net.Dialer{
								Timeout: 1 * time.Minute}).Dial,
							TLSHandshakeTimeout: 1 * time.Minute,
						}

						var client = &http.Client{
							Timeout:   time.Minute * 1,
							Transport: netTransport,
						}

						// Authorize client
						serBody, _ := json.Marshal(&struct {
							Username string `json:"username"`
							Password string `json:"password"`
						}{
							Username: "lakipay@bunnabanksc.com",
							Password: "Laki@1234",
						})

						req, err := http.NewRequest(http.MethodPost, "http://10.1.13.12/auth/login", bytes.NewBuffer(serBody))
						if err != nil {
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						req.Header.Set("Content-Type", "application/json")

						res, err := client.Do(req)
						if err != nil {
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						// body, _ := io.ReadAll(res.Body)
						uc.log.Println("Link")
						// uc.log.Println(string(body))
						uc.log.Println("Link")

						if res.StatusCode != http.StatusOK {
							// Unsuccessful request
							uc.log.Println("Send Error")
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: "Failed to link account",
							}
						}

						type AuthRes struct {
							Token string `json:"token"`
						}

						uc.log.Println("Link 1")
						var authRes AuthRes
						decoder := json.NewDecoder(res.Body)
						uc.log.Println("Link 2")
						// uc.log.Println(res.b)
						err = decoder.Decode(&authRes)
						uc.log.Println("Link 3")
						if err != nil {
							uc.log.Println(authRes.Token)
							uc.log.Println(err)
							uc.log.Println("Link 4")
							return nil, Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						defer res.Body.Close()

						uc.log.Println(txn.Reference)

						serBody, err = json.Marshal(&struct {
							CreditAccount string                   `json:"credit_account"`
							DebitAccount  string                   `json:"debit_account"`
							Date          time.Time                `json:"date"`
							Amount        float64                  `json:"amount"`
							Payloads      []map[string]interface{} `json:"payloads"`
						}{
							CreditAccount: "01320209107500",
							DebitAccount:  sender.Detail.(entity.BankAccount).Number,
							Amount:        amount,
							Date:          txn.CreatedAt,
							Payloads: []map[string]interface{}{
								{
									"txn_ref": txn.Reference,
								},
							},
						})

						uc.log.Println(txn.Reference)

						uc.log.Println("Bunna Bank 15")
						if err != nil {
							uc.log.Println("Bunna Bank 16")
							uc.log.Println(err)
							return nil, Error{
								Type:    "FAILED_TOVERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Bunna Bank 17")
						req, err = http.NewRequest(http.MethodPost, "http://10.1.13.12/api/core/transaction/open_c2c/initiate", bytes.NewBuffer(serBody))
						if err != nil {
							uc.log.Println("Bunna Bank 18")
							return nil, Error{
								Type:    "NO_CLIENT_AUTH_FOUND",
								Message: err.Error(),
							}
						}

						uc.log.Println("Bunna Bank 19")
						req.Header.Set("Content-Type", "application/json")
						req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authRes.Token))

						// Authorize request

						res, err = client.Do(req)
						if err != nil {
							uc.log.Println("Bunna Bank 20")
							return nil, Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Bunna Bank 21")
						body, err := io.ReadAll(res.Body)
						if err != nil {
							return nil, Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}
						defer res.Body.Close()

						if res.StatusCode != http.StatusOK {
							type TxnRes struct {
								Message string `json:"message"`
							}
							var txnRes TxnRes
							json.Unmarshal(body, &txnRes)
							uc.log.Println(txnRes.Message)
							uc.log.Println("Bunna Bank 22")
							return nil, Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: txnRes.Message,
							}
						}

						type TxnRes struct {
							Status         string `json:"status"`
							ResponseStatus string `json:"response_status"`
							ReferenceId    string `json:"reference_id"`
						}

						var txnRes TxnRes
						json.Unmarshal(body, &txnRes)

						txn.Reference = txnRes.ReferenceId

						uc.log.Println("Bunna Bank 23")
						// Store transaction
						err = uc.repo.StoreTransaction(txn)
						if err != nil {
							return nil, Error{
								Type:    "FAILED_TO_STORE_TRANSACTION",
								Message: err.Error(),
							}
						}
					}
				}

			}
		}
	case entity.P2P:
		fmt.Println("||||||||||||||||||||| p2p ||||||||||||||||||||||||")
		id := uuid.New()
		txn = entity.Transaction{
			Id:        id,
			From:      *sender,
			To:        *receipient,
			Type:      entity.TransactionType(txnType),
			Verified:  false,
			CreatedAt: time.Now(),
			Reference: strings.Split(uuid.New().String(), "-")[4],
			Details: entity.P2p{
				Amount: amount,
			},
		}

		err = uc.repo.StoreTransaction(txn)
		if err != nil {
			return nil, Error{
				Type:    "FAILED_TO_STORE_TRANSACTION",
				Message: err.Error(),
			}
		}

	}

	return &txn, nil
}

func (uc Usecase) VerifyTransaction(userId, txnId uuid.UUID, code string, amount float64) (*entity.Transaction, error) {
	// Find transaction
	txn, err := uc.repo.FindTransactionById(txnId)
	if err != nil {
		return nil, Error{
			Type:    "COULD_NOT_FIND_TRANSACTION",
			Message: err.Error(),
		}
	}

	sender, err := uc.repo.FindAccountById(txn.From.Id)
	if err != nil {
		return nil, Error{
			Type:    "COULD_NOT_FIND_TRANSACTION",
			Message: err.Error(),
		}
	}

	switch sender.Type {
	case entity.BANK:
		{
			switch sender.Detail.(entity.BankAccount).Bank.SwiftCode {
			case "AWINETAA":
				{
					uc.log.Println("Switching Amhara Bank")
					var netTransport = &http.Transport{
						Dial: (&net.Dialer{
							Timeout: 1 * time.Minute,
						}).Dial,
						TLSHandshakeTimeout: 1 * time.Minute,
					}

					var client = &http.Client{
						Timeout:   time.Minute * 1,
						Transport: netTransport,
					}

					// Authorize client
					serBody, _ := json.Marshal(&struct {
						Username string `json:"username"`
						Password string `json:"password"`
					}{
						Username: "qetxjgflmn",
						Password: "w9'MwO9F$n",
					})

					req, err := http.NewRequest(http.MethodPost, "http://10.10.101.144:8080/b2b/awash/api/v1/auth/getToken", bytes.NewBuffer(serBody))
					if err != nil {
						return nil, Error{
							Type:    "NO_RESPONSE",
							Message: err.Error(),
						}
					}

					req.Header.Set("Content-Type", "application/json")

					res, err := client.Do(req)
					if err != nil {
						return nil, Error{
							Type:    "NO_RESPONSE",
							Message: err.Error(),
						}
					}

					// body, _ := io.ReadAll(res.Body)
					uc.log.Println("Link")
					// uc.log.Println(string(body))
					uc.log.Println("Link")

					if res.StatusCode != http.StatusOK {
						// Unsuccessful request
						uc.log.Println("Send Error")
						return nil, Error{
							Type:    "NO_RESPONSE",
							Message: "Failed to link account",
						}
					}

					type AuthRes struct {
						Status string `json:"status"`
						Token  string `json:"token"`
					}

					uc.log.Println("Link 1")
					var authRes AuthRes
					decoder := json.NewDecoder(res.Body)
					uc.log.Println("Link 2")
					// uc.log.Println(res.b)
					err = decoder.Decode(&authRes)
					uc.log.Println("Link 3")
					if err != nil {
						uc.log.Println(authRes.Token)
						uc.log.Println(err)
						uc.log.Println("Link 4")
						return nil, Error{
							Type:    "NO_RESPONSE",
							Message: err.Error(),
						}
					}

					defer res.Body.Close()

					// Validate transaction
					/*

												{
						    "bankCode": "awsbnk",
						    "amount": 1,
						    "reference": "ASWQERDFTGHY",
						    "narration": "string",
						    "awashAccount": "01320209107500",
						    "creditAccount": "01320206449500",
						    "commisionAmount": 0,
						    "awashAccountName": "string",
						    "creditAccountName": "string"
						}

					*/

					serBody, err = json.Marshal(&struct {
						Phone string `json:"phone"`
						OTP   string `json:"otp"`
					}{
						Phone: sender.Detail.(entity.BankAccount).Holder.Phone,
						OTP:   code,
					})

					uc.log.Println("Amhara Bank 15")
					if err != nil {
						uc.log.Println("Amhara Bank 16")
						uc.log.Println(err)
						return nil, Error{
							Type:    "FAILED_TOVERIFY_TRANSACTION",
							Message: err.Error(),
						}
					}

					uc.log.Println("Amhara Bank 17")
					req, err = http.NewRequest(http.MethodPost, "http://10.10.101.144:8080/b2b/awash/api/v1/monetize/validate", bytes.NewBuffer(serBody))
					if err != nil {
						uc.log.Println("Amhara Bank 18")
						return nil, Error{
							Type:    "NO_CLIENT_AUTH_FOUND",
							Message: err.Error(),
						}
					}

					uc.log.Println("Amhara Bank 19")
					req.Header.Set("Content-Type", "application/json")
					req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authRes.Token))

					// Authorize request

					res, err = client.Do(req)
					if err != nil {
						uc.log.Println("Amhara Bank 20")
						return nil, Error{
							Type:    "FAILED_TO_VERIFY_TRANSACTION",
							Message: err.Error(),
						}
					}

					uc.log.Println("Amhara Bank 21")

					if res.StatusCode != http.StatusOK {
						body, _ := io.ReadAll(res.Body)
						type TxnRes struct {
							ResponseCode int    `json:"response_code"`
							Status       string `json:"status"`
							Message      string `json:"message"`
						}
						var txnRes TxnRes
						json.Unmarshal(body, &txnRes)
						uc.log.Println(txnRes.Message)
						uc.log.Println("Amhara Bank 22")
						return nil, Error{
							Type:    "FAILED_TO_VERIFY_TRANSACTION",
							Message: txnRes.Message,
						}
					}

					uc.log.Println("Amhara Bank 23")
					// Store transaction
					receipient, _ := uc.repo.FindAccountById(txn.To.Id)
					txn.Verified = true
					uc.repo.UpdateAccount(entity.Account{
						Id:                 receipient.Id,
						Title:              receipient.Title,
						Type:               receipient.Type,
						Default:            receipient.Default,
						User:               receipient.User,
						VerificationStatus: receipient.VerificationStatus,
						Detail: entity.StoredAccount{
							Balance: receipient.Detail.(entity.StoredAccount).Balance + amount,
						},
					})
					if err != nil {
						return nil, Error{
							Type:    "FAILED_TO_STORE_TRANSACTION",
							Message: err.Error(),
						}
					}
				}
			}
		}
	}

	return txn, nil

}

func (uc Usecase) GetUserTransactions(id uuid.UUID) ([]entity.Transaction, error) {

	// Check policy

	txs, err := uc.repo.FindTransactionsByUserId(id)
	if err != nil {
		return nil, Error{
			Type:    "COULD_NOT_FIND_TRANSACTIONS",
			Message: err.Error(),
		}
	}

	uc.log.Println("TXNS")
	uc.log.Println(len(txs))
	uc.log.Println(txs)

	return txs, nil
}

func (uc Usecase) GetAllTransactions() ([]entity.Transaction, error) {

	// Check policy

	txs, err := uc.repo.FindAllTransactions()
	if err != nil {
		return nil, Error{
			Type:    "COULD_NOT_FIND_TRANSACTIONS",
			Message: err.Error(),
		}
	}

	uc.log.Println("TXNS")
	uc.log.Println(len(txs))
	uc.log.Println(txs)

	return txs, nil
}

func (uc Usecase) TransactionsDashboardUsecase(year int) (interface{}, error) {

	// Check policy

	txs, err := uc.repo.TransactionsDashboardRepo(year)
	if err != nil {
		return nil, Error{
			Type:    "COULD_NOT_FIND_TRANSACTIONS",
			Message: err.Error(),
		}
	}

	return txs, nil
}

// func (uc Usecase) GetHotelTransactions(id uuid.UUID) ([]entity.Transaction, error) {

// 	// Check policy

// 	txs, err := uc.repo.FindTransactionsByUserId(id)
// 	if err != nil {
// 		return nil, Error{
// 			Type:    "COULD_NOT_FIND_TRANSACTIONS",
// 			Message: err.Error(),
// 		}
// 	}

// 	uc.log.Println("TXNS")
// 	uc.log.Println(len(txs))
// 	uc.log.Println(txs)

// 	return txs, nil
// }
