package usecase

import (
	"auth/src/pkg/account/core/entity"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (uc Usecase) GetUserAccounts(id uuid.UUID) ([]entity.Account, error) {
	var accs []entity.Account

	// Find Accounts
	accs, err := uc.repo.FindAccountsByUserId(id)

	if err != nil {
		print(err)
		return nil, Error{
			Type:    "COULD_NOT_FIND_ACCOUNTS",
			Message: err.Error(),
		}
	}

	print(len(accs))
	// Check if there is a stored account
	exists := false
	for i := 0; i < len(accs); i++ {
		if accs[i].Type == entity.STORED {
			exists = true
			break
		}
	}

	if !exists {
		acc, err := uc.CreateStoredAccount(id, "Stored Value", true)
		if err != nil {
			return nil, Error{
				Type:    "COULD_NOT_FIND_ACCOUNTS",
				Message: err.Error(),
			}
		}

		accs = append(accs, *acc)
	}

	return accs, nil
}

// Create a stored account
func (uc Usecase) CreateStoredAccount(user uuid.UUID, title string, isDefault bool) (*entity.Account, error) {

	// Create a stored account
	var acc entity.Account
	id := uuid.New()
	acc = entity.Account{
		Id:      id,
		Title:   title,
		Type:    entity.STORED,
		Default: isDefault,
		VerificationStatus: struct {
			Verified   bool
			VerifiedBy *struct {
				Method  string
				Details interface{}
			}
		}{
			Verified: true,
		},
		Detail: entity.StoredAccount{
			Balance: 0.0,
		},
		User: entity.User{
			Id: user,
		},
	}

	// Store account
	err := uc.repo.StoreAccount(acc)
	if err != nil {
		return nil, Error{
			Type:    "FAILED_TO_STORE_ACCOUNT",
			Message: err.Error(),
		}
	}

	return &acc, nil
}

type AccountVerification struct {
	Status bool
}

// Create bank account
func (uc Usecase) CreateBankAccount(userId uuid.UUID, bankId uuid.UUID, accountNumber string, accountHolderName string, accountHolderPhone string, title string, makeDefault bool) (*entity.Account, error) {
	var acc entity.Account

	// Find Bank
	bank, err := uc.repo.FindBankById(bankId)
	if err != nil {
		return nil, Error{
			Type:    "FAILED_TO_ADD_ACCOUNT",
			Message: err.Error(),
		}
	}

	// Create account
	acc = entity.Account{
		Id:      uuid.New(),
		Title:   title,
		Type:    entity.BANK,
		Default: makeDefault,
		User: entity.User{
			Id: userId,
		},
		Detail: entity.BankAccount{
			Bank:   *bank,
			Number: accountNumber,
			Holder: struct {
				Name  string
				Phone string
			}{
				Name:  accountHolderName,
				Phone: accountHolderPhone,
			},
		},
		CreatedAt: time.Now(),
	}

	// Return response

	// Bank specific operation

	switch bank.SwiftCode {
	case "AMHRETAA":
		{
			uc.log.Println("Switching Amhara Bank")
			// var netTransport = &http.Transport{
			// 	Dial: (&net.Dialer{
			// 		Timeout: 1 * time.Minute,
			// 	}).Dial,
			// 	TLSHandshakeTimeout: 1 * time.Minute,
			// }

			// var client = &http.Client{
			// 	Timeout:   time.Minute * 1,
			// 	Transport: netTransport,
			// }

			// // Authorize client
			// serBody, _ := json.Marshal(&struct {
			// 	Username string `json:"username"`
			// 	Password string `json:"password"`
			// }{
			// 	Username: "LakiPay",
			// 	Password: "e3i1OehzfV0Iz16asdTjZEbYG4F769Vx8Unuo5chkM9V",
			// })

			// req, err := http.NewRequest(http.MethodPost, "http://172.31.2.30:8600/abaApi/v1/lakiPay/authenticate", bytes.NewBuffer(serBody))
			// if err != nil {
			// 	return nil, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// req.Header.Set("Content-Type", "application/json")

			// res, err := client.Do(req)
			// if err != nil {
			// 	return nil, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// // body, _ := io.ReadAll(res.Body)
			// uc.log.Println("Link")
			// // uc.log.Println(string(body))
			// uc.log.Println("Link")

			// if res.StatusCode != http.StatusOK {
			// 	// Unsuccessful request
			// 	uc.log.Println("Send Error")
			// 	return nil, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: "Failed to link account",
			// 	}
			// }

			// type AuthRes struct {
			// 	ResponseCode int    `json:"response_code"`
			// 	Status       string `json:"status"`
			// 	Message      string `json:"message"`
			// 	Token        string `json:"token"`
			// }

			// uc.log.Println("Link 1")
			// var authRes AuthRes
			// decoder := json.NewDecoder(res.Body)
			// uc.log.Println("Link 2")
			// // uc.log.Println(res.b)
			// err = decoder.Decode(&authRes)
			// uc.log.Println("Link 3")
			// if err != nil {
			// 	uc.log.Println(authRes.Token)
			// 	uc.log.Println(err)
			// 	uc.log.Println("Link 4")
			// 	return nil, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// defer res.Body.Close()

			// // Send link request

			// uc.log.Println("Link")

			// serBody, _ = json.Marshal(&struct {
			// 	AccountNumber string `json:"account_number"`
			// 	AccountHolder string `json:"account_holder"`
			// }{
			// 	AccountNumber: accountNumber,
			// 	AccountHolder: accountHolderPhone,
			// })

			// req, err = http.NewRequest(http.MethodPost, "http://172.31.2.30:8600/abaApi/v1/lakiPay/linkAccount", bytes.NewBuffer(serBody))
			// if err != nil {
			// 	return nil, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// req.Header.Set("Content-Type", "application/json")
			// req.Header.Add("Authorization", authRes.Token)

			// res, err = client.Do(req)
			// if err != nil {
			// 	return nil, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// type LinkRes struct {
			// 	ResponseCode int    `json:"response_code"`
			// 	Status       string `json:"status"`
			// 	Message      string `json:"message"`
			// 	Length       int    `json:"length"`
			// 	Timeout      int    `json:"timeout"`
			// }

			// var linkRes LinkRes
			// decoder = json.NewDecoder(res.Body)
			// err = decoder.Decode(&linkRes)
			// if err != nil {
			// 	return nil, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// defer res.Body.Close()

			// if res.StatusCode != http.StatusOK {
			// 	// Unsuccessful request
			// 	uc.log.Println("Send Error")
			// 	return nil, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: linkRes.Message,
			// 	}
			// }

			// Send response
			time.Sleep(6 * time.Second)

			// Store account
			acc.VerificationStatus = struct {
				Verified   bool
				VerifiedBy *struct {
					Method  string
					Details interface{}
				}
			}{
				Verified: false,
				VerifiedBy: &struct {
					Method  string
					Details interface{}
				}{
					Method: "SMS",
					Details: struct {
						Length  int
						Timeout int
					}{
						Length:  5,
						Timeout: 120,
					},
				},
			}
			err = uc.repo.StoreAccount(acc)
			if err != nil {
				return nil, Error{
					Type:    "FAILED_TO_ADD_ACCOUNT",
					Message: err.Error(),
				}
			}
		}
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
				Token       string `json:"token"`
				Description string `json:"errorDescription"`
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

			uc.log.Println(accountNumber)

			// Verify Acc
			req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("http://10.10.101.144:8080/b2b/awash/api/v1/monetize/getAccount?bankCode=awsbnk&account=%s", accountNumber), nil)
			if err != nil {
				return nil, Error{
					Type:    "NO_RESPONSE",
					Message: err.Error(),
				}
			}

			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authRes.Token))

			res, err = client.Do(req)
			if err != nil {
				return nil, Error{
					Type:    "NO_RESPONSE",
					Message: err.Error(),
				}
			}

			if res.StatusCode != http.StatusOK {
				// Unsuccessful request
				uc.log.Println("Send Error")
				return nil, Error{
					Type:    "NO_RESPONSE",
					Message: "Failed to link account",
				}
			}

			type VerifyAccRes struct {
				Status      int    `json:"status"`
				Description string `json:"errorDescription"`
			}

			uc.log.Println("Link 1")
			var verifyRes VerifyAccRes
			decoder = json.NewDecoder(res.Body)
			uc.log.Println("Link 2")
			// uc.log.Println(res.b)
			err = decoder.Decode(&verifyRes)
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

			if verifyRes.Status != 1 {
				// Unsuccessful request
				uc.log.Println("Send Error status")
				uc.log.Println(verifyRes.Status)
				uc.log.Println(verifyRes.Description)
				return nil, Error{
					Type:    "NO_RESPONSE",
					Message: authRes.Description,
				}
			}

			// Send Response
			// Store account
			acc.VerificationStatus = struct {
				Verified   bool
				VerifiedBy *struct {
					Method  string
					Details interface{}
				}
			}{
				Verified: true,
			}
			err = uc.repo.StoreAccount(acc)
			if err != nil {
				return nil, Error{
					Type:    "FAILED_TO_ADD_ACCOUNT",
					Message: err.Error(),
				}
			}
		}
	case "ORIRETAA":
		{
			uc.log.Println("Switching Oromia Bank")
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

			uc.log.Println(accountNumber)

			serBody, _ := json.Marshal(&struct {
				AccountNumber string `json:"accountNumber"`
			}{
				AccountNumber: accountNumber,
			})

			// Verify Acc
			req, err := http.NewRequest(http.MethodPost, "http://10.10.20.47/customer/query-account-holder", bytes.NewBuffer(serBody))
			if err != nil {
				return nil, Error{
					Type:    "NO_RESPONSE",
					Message: err.Error(),
				}
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJPQiIsImp0aSI6ImM1OTg0YTc2YTAyMjA1MDIwNzQ1MTliYThhNWU0OWMzMDk3NTJmMTAyYThhNzhkYjNmNThiM2QxMzAxMzhiMjEiLCJzdWIiOiJsYWtpcGF5IiwiaWF0IjoxNzAxODUwNzc5fQ.sD_C4nwadpgClQADGOPjWjKembyxqCit2tmD_rLsOg7NsFVDv2xbzvnvDnAjD0OKZSfEfhfuKKHsOZfx1crbAA"))

			res, err := client.Do(req)
			if err != nil {
				return nil, Error{
					Type:    "NO_RESPONSE",
					Message: err.Error(),
				}
			}

			if res.StatusCode != http.StatusOK {
				// Unsuccessful request
				uc.log.Println("Send Error")
				uc.log.Println(res.StatusCode)
				return nil, Error{
					Type:    "NO_RESPONSE",
					Message: "Failed to link account",
				}
			}

			// Send Response
			// Store account
			acc.VerificationStatus = struct {
				Verified   bool
				VerifiedBy *struct {
					Method  string
					Details interface{}
				}
			}{
				Verified: true,
			}
			err = uc.repo.StoreAccount(acc)
			if err != nil {
				return nil, Error{
					Type:    "FAILED_TO_ADD_ACCOUNT",
					Message: err.Error(),
				}
			}
		}
	case "BUNAETAA":
		{
			uc.log.Println("Switching Bunna Bank")

			// Send Response
			// Store account
			acc.VerificationStatus = struct {
				Verified   bool
				VerifiedBy *struct {
					Method  string
					Details interface{}
				}
			}{
				Verified: true,
			}
			err = uc.repo.StoreAccount(acc)
			if err != nil {
				return nil, Error{
					Type:    "FAILED_TO_ADD_ACCOUNT",
					Message: err.Error(),
				}
			}
		}
	}

	uc.log.Println(acc.VerificationStatus.Verified)

	return &acc, nil
}

func (uc Usecase) VerifyAccount(userId, accountId uuid.UUID, method string, details interface{}, code string) (string, error) {
	var token string
	acc, err := uc.repo.FindAccountById(accountId)
	if err != nil {
		return token, Error{
			Type:    "ACCOUNT_NOT_FOUND",
			Message: err.Error(),
		}
	}

	// Bank specific verifications
	switch acc.Detail.(entity.BankAccount).Bank.SwiftCode {
	case "AMHRETAA":
		{
			// var netTransport = &http.Transport{
			// 	Dial: (&net.Dialer{
			// 		Timeout: 1 * time.Minute,
			// 	}).Dial,
			// 	TLSHandshakeTimeout: 1 * time.Minute,
			// }

			// var client = &http.Client{
			// 	Timeout:   time.Minute * 1,
			// 	Transport: netTransport,
			// }

			// // Authorize client
			// serBody, _ := json.Marshal(&struct {
			// 	Username string `json:"username"`
			// 	Password string `json:"password"`
			// }{
			// 	Username: "LakiPay",
			// 	Password: "e3i1OehzfV0Iz16asdTjZEbYG4F769Vx8Unuo5chkM9V",
			// })

			// req, err := http.NewRequest(http.MethodPost, "http://172.31.2.30:8600/abaApi/v1/lakiPay/authenticate", bytes.NewBuffer(serBody))
			// if err != nil {
			// 	return token, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// req.Header.Set("Content-Type", "application/json")

			// res, err := client.Do(req)
			// if err != nil {
			// 	return token, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// // body, _ := io.ReadAll(res.Body)
			// // uc.log.Println(string(body))

			// if res.StatusCode != http.StatusOK {
			// 	// Unsuccessful request
			// 	uc.log.Println("Send Error")
			// 	return token, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: "Failed to link account",
			// 	}
			// }

			// type AuthRes struct {
			// 	ResponseCode int    `json:"response_code"`
			// 	Status       string `json:"status"`
			// 	Message      string `json:"message"`
			// 	Token        string `json:"token"`
			// }

			// var authRes AuthRes
			// decoder := json.NewDecoder(res.Body)
			// err = decoder.Decode(&authRes)
			// if err != nil {
			// 	return token, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// defer res.Body.Close()

			// Verify Account

			// serBody, err = json.Marshal(struct {
			// 	AccountNumber string `json:"account_number"`
			// 	AccountHolder string `json:"account_holder"`
			// 	Code          string `json:"code"`
			// }{
			// 	AccountNumber: acc.Detail.(entity.BankAccount).Number,
			// 	AccountHolder: acc.Detail.(entity.BankAccount).Holder.Phone,
			// 	Code:          code,
			// })

			// if err != nil {
			// 	return token, Error{
			// 		Type:    "FAILED_TO_VERIFY",
			// 		Message: err.Error(),
			// 	}
			// }

			// req, err = http.NewRequest(http.MethodPost, "http://172.31.2.30:8600/abaApi/v1/lakiPay/linkAccount/verifyOTP", bytes.NewBuffer(serBody))
			// if err != nil {
			// 	return token, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// req.Header.Set("Content-Type", "application/json")
			// req.Header.Add("Authorization", authRes.Token)

			// res, err = client.Do(req)
			// if err != nil {
			// 	return token, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// type VerifyRes struct {
			// 	ResponseCode int    `json:"response_code"`
			// 	Status       string `json:"status"`
			// 	Message      string `json:"message"`
			// 	Token        string `json:"token"`
			// }

			// var verifyRes VerifyRes
			// decoder = json.NewDecoder(res.Body)
			// err = decoder.Decode(&verifyRes)
			// if err != nil {
			// 	return token, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// defer res.Body.Close()

			// if res.StatusCode != http.StatusOK {
			// 	// Unsuccessful request
			// 	body, _ := io.ReadAll(res.Body)
			// 	uc.log.Println(string(body))
			// 	uc.log.Println(res.StatusCode)
			// 	uc.log.Println("Send Error")
			// 	return token, Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: verifyRes.Message,
			// 	}
			// }

			// token = verifyRes.Token
			time.Sleep(4 * time.Second)
			token = "verifyRes.Token"
		}
	}

	// Update account

	acc.VerificationStatus = struct {
		Verified   bool
		VerifiedBy *struct {
			Method  string
			Details interface{}
		}
	}{
		Verified: true,
	}
	uc.repo.UpdateAccount(*acc)

	return token, nil
}

func (uc Usecase) DeleteAccount(userId, accId uuid.UUID) error {
	acc, err := uc.repo.FindAccountById(accId)
	if err != nil {
		return Error{
			Type:    "ACCOUNT_NOT_FOUND",
			Message: err.Error(),
		}
	}

	// Bank specific verifications
	switch acc.Detail.(entity.BankAccount).Bank.SwiftCode {
	case "AMHRETAA":
		{
			// var netTransport = &http.Transport{
			// 	Dial: (&net.Dialer{
			// 		Timeout: 1 * time.Minute,
			// 	}).Dial,
			// 	TLSHandshakeTimeout: 1 * time.Minute,
			// }

			// var client = &http.Client{
			// 	Timeout:   time.Minute * 1,
			// 	Transport: netTransport,
			// }

			// // Authorize client
			// serBody, _ := json.Marshal(&struct {
			// 	Username string `json:"username"`
			// 	Password string `json:"password"`
			// }{
			// 	Username: "LakiPay",
			// 	Password: "e3i1OehzfV0Iz16asdTjZEbYG4F769Vx8Unuo5chkM9V",
			// })

			// req, err := http.NewRequest(http.MethodPost, "http://172.31.2.30:8600/abaApi/v1/lakiPay/authenticate", bytes.NewBuffer(serBody))
			// if err != nil {
			// 	return Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// req.Header.Set("Content-Type", "application/json")

			// res, err := client.Do(req)
			// if err != nil {
			// 	return Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// // body, _ := io.ReadAll(res.Body)
			// // uc.log.Println(string(body))

			// if res.StatusCode != http.StatusOK {
			// 	// Unsuccessful request
			// 	uc.log.Println("Send Error")
			// 	return Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: "Failed to link account",
			// 	}
			// }

			// type AuthRes struct {
			// 	ResponseCode int    `json:"response_code"`
			// 	Status       string `json:"status"`
			// 	Message      string `json:"message"`
			// 	Token        string `json:"token"`
			// }

			// var authRes AuthRes
			// decoder := json.NewDecoder(res.Body)
			// err = decoder.Decode(&authRes)
			// if err != nil {
			// 	return Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// defer res.Body.Close()

			// // Verify Account

			// serBody, err = json.Marshal(struct {
			// 	AccountNumber string `json:"account_number"`
			// 	AccountHolder string `json:"account_holder"`
			// }{
			// 	AccountNumber: acc.Detail.(entity.BankAccount).Number,
			// 	AccountHolder: acc.Detail.(entity.BankAccount).Holder.Phone,
			// })

			// if err != nil {
			// 	return Error{
			// 		Type:    "FAILED_TO_VERIFY",
			// 		Message: err.Error(),
			// 	}
			// }

			// req, err = http.NewRequest(http.MethodPost, "http://172.31.2.30:8600/abaApi/v1/lakiPay/unlinkAccount", bytes.NewBuffer(serBody))
			// if err != nil {
			// 	return Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// req.Header.Set("Content-Type", "application/json")
			// req.Header.Add("Authorization", authRes.Token)

			// res, err = client.Do(req)
			// if err != nil {
			// 	return Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: err.Error(),
			// 	}
			// }

			// // body, _ = io.ReadAll(res.Body)
			// // uc.log.Println(string(body))

			// defer res.Body.Close()

			// if res.StatusCode != http.StatusOK {
			// 	// Unsuccessful request
			// 	uc.log.Println("Send Error")
			// 	return Error{
			// 		Type:    "NO_RESPONSE",
			// 		Message: "Failed to link account",
			// 	}
			// }
		}
	}

	// Update account

	uc.repo.DeleteAccount(acc.Id)

	return nil
}
