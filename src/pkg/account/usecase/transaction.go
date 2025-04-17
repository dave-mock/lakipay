package usecase

import (
	"auth/src/pkg/account/core/entity"
	auth_repo "auth/src/pkg/auth/adapter/gateway/repo/psql"
	auth_entity "auth/src/pkg/auth/core/entity"
	db_psql "auth/src/pkg/auth/infra/storage/psql"
	"auth/src/pkg/jwt"
	"auth/src/pkg/utils"
	"bytes"
	"crypto"
	rand2 "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Create transaction

func (uc Usecase) UpdateUserUsecase(users entity.User2) (entity.User2, error) {

	users, err := uc.repo.UpdateUserRepo(users)

	if err != nil {
		return entity.User2{}, err
	}

	return users, err
}

func (uc Usecase) InitPreSession(txtId uuid.UUID) (entity.TransactionSession, error) {

	uc.log.SetPrefix("[AUTH] [USECASE] [InitPreSession] ")

	// Errors
	var ErrFailedToInitiateAuth string = "FAILED_TO_CREATE_PRE_SESSION"

	// id := uuid.New()

	token := jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + 5000,
		Public: txtId,
	}, "pre_session_secret")

	preSession := entity.TransactionSession{
		Id:        txtId,
		Token:     token,
		CreatedAt: time.Now(),
	}

	uc.log.Println("created pre session")

	// Store pre session record
	err := uc.repo.StoreTransactionSession(preSession)
	uc.log.Println("storing pre session")
	if err != nil {
		uc.log.Printf("failed to store pre session : %s\n", err)
		return preSession, Error{
			Type:    ErrFailedToInitiateAuth,
			Message: err.Error(),
		}
	}

	uc.log.Println("initiated authentication")
	// Return presession
	return preSession, nil
}

func (uc Usecase) GetApiKeysUsecase(id uuid.UUID) (string, error) {

	merchantKeys, err := uc.repo.GetApiKeysRepo(id)
	if err != nil {
		return "", err
	}

	return merchantKeys.PrivateKey, nil
}

func (uc Usecase) CreateRegisterKeys(id uuid.UUID, username string, password string) (string, error) {
	fmt.Print("|||||| private: ")

	private, public, err := generateRSAKeyPair()

	if err != nil {
		return "", err
	}

	err, isUser_exist := uc.repo.CheckMerchantsKeysByUsername(username)

	if err != nil {
		return "", err
	}

	if isUser_exist {

		return "", Error{
			Type:    "Error",
			Message: "Username is already used",
		}
	}

	hasher := sha256.New()
	_, err = hasher.Write([]byte(password))
	if err != nil {
		return "", err
	}
	pwd := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	fmt.Println("****************************************************************************************")
	err = uc.repo.StoreKeys(id, public, private, pwd, username)
	if err != nil {
		return "", err
	}

	return private, nil
}

func (uc Usecase) ApplyForTokenUsecase(username string, password string) (string, error) {

	err, isUser_exist := uc.repo.CheckMerchantsKeysByUsername(username)

	if err != nil {
		return "", err
	}

	if !isUser_exist {

		return "", Error{
			Type:    "Error",
			Message: "Username is not found",
		}
	}

	merchant, err := uc.repo.GetMerchantsKeys(username)

	if err != nil {

		return "", Error{
			Type:    "Error",
			Message: "Merchant not found",
		}
	}

	hasher := sha256.New()
	_, err = hasher.Write([]byte(password))
	if err != nil {
		return "", err
	}
	pwd := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// compare passwords
	if pwd != merchant.Password {
		return "", Error{
			Type:    "INCORRECT_OASSWORD",
			Message: "Pasword is incorect",
		}
	}

	token := jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + 5000,
		Public: merchant.MerchantId,
	}, "initaite_secret")

	return token, nil
}

func (uc Usecase) CreateHostedTransactionInitiate(amount float64, currency string, callback_url string, secretKey string, stringData string, token string) (interface{}, error) {
	// var receipient entity.Account
	// var txn entity.Transaction

	val, err := jwt.Decode(token, "initaite_secret")
	if err != nil {
		uc.log.Printf("failed checking presession : %s\n", err.Error())
		return "", Error{
			Type:    "ErrInvalidPreSessionToken",
			Message: err.Error(),
		}
	}

	uuidStr, ok := val.Public.(string)
	if !ok {
		fmt.Println("Value is not a string")
		return "", err
	}

	// parsedUUID, err := uuid.Parse(uuidStr)
	// if err != nil {
	//     fmt.Println("Error parsing UUID:", err)
	//     return "", err
	// }

	// merchantKeys,err:= uc.repo.GetApiKeysRepo(parsedUUID)
	// 		if err != nil{
	// 			return "",err
	// 		}

	// var isVerify bool = false

	// if verify(merchantKeys.PublickKey, stringData, secretKey) {

	// 	isVerify = true
	// }

	// fmt.Println("====== ",merchantKeys.Id)

	// if !isVerify  {
	// 	return "", Error{
	// 		Type:    "FAILED_VERIFICATION",
	// 		Message: "Verification error",
	// 	}
	// }

	// receipient.Id = merchant_kyes.MerchantId

	// type Detail struct{
	// 	Amount float64
	// 	CallbackUrl string
	// }

	// id := uuid.New()
	// 		txn = entity.Transaction{
	// 			Id:        id,
	// 			To:        receipient,
	// 			Verified:  false,
	// 			CreatedAt: time.Now(),
	// 			Amount:    amount,
	// 			Currency: "birr",
	// 			Reference: strings.Split(uuid.New().String(), "-")[4],
	// 			Details:Detail{
	// 				Amount: amount,
	// 				CallbackUrl:callback_url,
	// 			},
	// 		}

	// fmt.Println("---------------------------------------------------------------------------------------------------------",txn.HasChallenge)

	// 		err = uc.repo.StoreTransaction(txn)
	// 		if err != nil {
	// 			return nil, Error{
	// 				Type:    "FAILED_TO_STORE_TRANSACTION",
	// 				Message: err.Error(),
	// 			}
	// 		}

	amountStr := fmt.Sprintf("%.2f", amount)
	fmt.Println("UUID user_id", uuidStr)

	encString := string(uuidStr + "," + amountStr + "," + currency + "," + callback_url)

	encriptedData, err := utils.AesEncrption(encString)
	if err != nil {
		return nil, Error{
			Type:    "FAILED_TO_ENCRPT",
			Message: err.Error(),
		}
	}

	fmt.Println("+++++++++++++", encString)
	type res struct {
		HostedURL string `json:"hosted_url"`
	}
	var response res
	hostedURL := fmt.Sprintf("http://196.189.126.183:3005/%s", encriptedData)
	response.HostedURL = hostedURL
	return response, nil
}

func (uc Usecase) CheckBalance(from uuid.UUID) (float64, error) {
	// x:="s"
	txns, err := uc.repo.FindTransactionsByUserId(from)

	if err != nil {
		return 0, err
	}
	blance := 0.0

	for _, txn := range txns {
		if txn.From.Id == from {
			blance = blance - txn.TotalAmount
		} else if txn.To.Id == from {
			blance = blance + txn.Amount
		}
	}

	fmt.Println("BLANCE __ : ", blance)

	return blance, nil
}

func (uc Usecase) CreateTransactionInitiate(userId uuid.UUID, from uuid.UUID, to uuid.UUID, amount float64, mediums entity.TransactionMedium, txnType string, token string, detail string) (interface{}, error) {

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
		fmt.Println(mediums)
		if mediums == "LAKIPAY" {

			balance, err := uc.CheckBalance(from)
			if err != nil {
				return nil, Error{
					Type:    "NOT_ENOUGH_FUND",
					Message: "amount is larger than your balance",
				}
			}

			if balance < amount {
				return nil, Error{
					Type:    "NOT_ENOUGH_FUND",
					Message: "amount is larger than your balance",
				}
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
	} else {
		receipient = &entity.Account{}
	}

	type Res struct {
		Token        string `json:"token"`
		HasChallenge bool   `json:"has_challenge"`
	}

	has_challenge := true

	if amount < 0 {
		has_challenge = false
	}

	id := uuid.New()
	txn = entity.Transaction{
		Id:           id,
		From:         *sender,
		To:           *receipient,
		Type:         entity.TransactionType(txnType),
		Medium:       mediums,
		Verified:     false,
		CreatedAt:    time.Now(),
		Amount:       amount,
		HasChallenge: has_challenge,
		Reference:    strings.Split(uuid.New().String(), "-")[4],
		Details: entity.P2p{
			Amount: amount,
		},
		Phone: detail,
	}

	err = uc.repo.StoreTransaction(txn)
	if err != nil {
		return nil, Error{
			Type:    "FAILED_TO_STORE_TRANSACTION",
			Message: err.Error(),
		}
	}

	transactionSession, err := uc.InitPreSession(txn.Id)

	if err != nil {
		return nil, Error{
			Type:    "Sommting is wrong",
			Message: err.Error(),
		}
	}
	var response Res
	response.Token = transactionSession.Token
	response.HasChallenge = has_challenge

	return response, nil
}

func (uc Usecase) VerifyTransaction(UserId uuid.UUID, token string, challenge entity.TransactionChallange, challengeType string) (string, error) {
	var ErrInvalidPreSessionToken string = "INVALID_PRE_SESSION_TOKEN"
	var sender *entity.Account
	var receipient *entity.Account
	var err error

	val, err := jwt.Decode(token, "pre_session_secret")
	if err != nil {
		uc.log.Printf("failed checking presession : %s\n", err.Error())
		return "", Error{
			Type:    ErrInvalidPreSessionToken,
			Message: err.Error(),
		}
	}

	uuidStr, ok := val.Public.(string)
	if !ok {
		fmt.Println("Value is not a string")
		return "", err
	}

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return "", err
	}

	// get accounts

	fmt.Print("{{{{{{{{{{{{{{{{{{onr teo three}}}}}}}}}}}}}}}}}}")

	txt, err := uc.repo.FindTransactionById(parsedUUID)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return "", err
	}
	fmt.Println("--------------------------------------------------------------------------------------------------------- two")

	// 1aaaf438-3d67-41d1-9503-9ffce899bc25
	// jsonData, err := json.Marshal(txt)

	// if err != nil {
	//     fmt.Println("Error parsing UUID:", err)
	//     return "", err
	// }
	// fmt.Println("|||||||||| txt: ",string(jsonData))

	if txt.From.Id != uuid.Nil {
		sender, err = uc.repo.FindAccountById(txt.From.Id)
		if err != nil {
			return "", Error{
				Type:    "SENDER_ACCOUNT_NOT_FOUND",
				Message: err.Error(),
			}
		}

		if UserId != sender.User.Id {
			return "", Error{
				Type:    "SENDER_ACCOUNT_MISMATCH",
				Message: "Please use your own account",
			}
		}

		// // If it's a stored value check balance
		// if sender.Type == entity.STORED && sender.Detail.(entity.StoredAccount).Balance < txt.Amount {
		// 	return "", Error{
		// 		Type:    "NOT_ENOUGH_FUND",
		// 		Message: "amount is larger than your balance",
		// 	}
		// }
	}

	//
	if txt.Tag.Id != uuid.Nil {
		receipient, err = uc.repo.FindAccountById(txt.To.Id)
		if err != nil {
			return "", Error{
				Type:    "RECEPIENT_ACCOUNT_NOT_FOUND",
				Message: err.Error(),
			}
		}
	} else {
		receipient = &entity.Account{}
	}

	fmt.Println("-------------------------==========================", txt.HasChallenge)
	fmt.Println("-------------------------==========================", !(challengeType == "2fa" || challengeType == "otp" || challengeType == "finger_print"))

	if txt.HasChallenge {
		if !(challengeType == "2fa" || challengeType == "otp" || challengeType == "finger_print") {
			return "", Error{
				Type:    "Challenge_NOT_PASTED",
				Message: "you have to pass the challenge",
			}
		}
	}

	switch challengeType {

	case "2fa":
		{
			log1 := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

			// [Ouput Adapters]
			// [DB] Postgres
			db, err := db_psql.New(log1)
			if err != nil {
				log.Fatal(err.Error())
			}
			r, err := auth_repo.NewPsqlRepo(uc.log, db)
			if err != nil {
				return "", Error{
					Type:    "FAILED_TO_CHEKE_2FA",
					Message: err.Error(),
				}
			}
			pass, err := r.FindPasswordIdentityByUser(UserId)
			if err != nil {
				return "", Error{
					Type:    "FAILED_TO_CHEKE_2FA11",
					Message: err.Error(),
				}
			}

			if pass == nil {
				return "", Error{
					Type:    "FAILED_TO_CHEKE_2FA22",
					Message: err.Error(),
				}
			}

			hasher := sha256.New()
			_, err = hasher.Write([]byte(challenge.TwoFA))
			if err != nil {
				return "", Error{
					Type:    "FAILED_TO_CHEKE_2FA33",
					Message: err.Error(),
				}
			}

			// Compare passwords

			if base64.URLEncoding.EncodeToString(hasher.Sum(nil)) != pass.Password {
				return "", Error{
					Type:    "FAILED_TO_CHEKE_2FA",
					Message: "2fa faild",
				}
			}
		}
	case "otp":
		{
			log1 := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

			// [Ouput Adapters]
			// [DB] Postgres
			db, err := db_psql.New(log1)
			if err != nil {
				log.Fatal(err.Error())
			}
			r3, err := auth_repo.NewPsqlRepo(uc.log, db)
			if err != nil {
				return "", Error{
					Type:    "FAILED_TO_CHEKE_2FA",
					Message: err.Error(),
				}
			}
			phoneAuth, err := r3.FindPhoneAuthWithoutPhone(UserId.String())
			if err != nil {
				log.Println("1")
				return "", Error{
					Type:    "FAILED_TO_AUTH_PHONE",
					Message: err.Error(),
				}
			}
			// // Validate phone
			// if phoneAuth.Phone.Prefix != prefix || phoneAuth.Phone.Number != number {
			// 	log.Println("2")
			// 	return "", Error{
			// 		Type:    "FAILED_TO_AUTH_PHONE",
			// 		Message: "Phone is not valid",
			// 	}
			// }

			// Check the status
			if phoneAuth.Status {
				return "", Error{
					Type:    "",
					Message: "already used",
				}
			}

			_, err = jwt.Decode(phoneAuth.Code, challenge.OTP)
			if err != nil {
				if err.Error() == "invalid token" {
					return "", Error{
						Type:    "",
						Message: "Incorrect OTP",
					}
				}
				return "", Error{
					Type:    "",
					Message: err.Error(),
				}
			}
			err = r3.UpdatePhoneAuthStatus(phoneAuth.Id, true)
			if err != nil {
				log.Println("5")
				return "", Error{
					Type:    "",
					Message: err.Error(),
				}
			}

		}
	case "finger_print":
		{
			public_keys, err := uc.repo.GetPuplicKey(challenge.Challenge, UserId)

			fmt.Println("|||||||||||| ====== public key", public_keys)
			if err != nil {

				return "", Error{
					Type:    "FAILED_TO_GET_PUBLIC_KEY",
					Message: err.Error(),
				}
			}

			var isVerify bool = false
			for _, pk := range public_keys {
				if verify(pk.PublicKey, challenge.Challenge, challenge.Signature) {
					err = uc.repo.UpdatePublicKeysUsed(UserId)

					if err != nil {
						return "", Error{
							Type:    "FAILED_VERIFICATION",
							Message: err.Error(),
						}
					}

					isVerify = true
				}
			}

			if !isVerify {
				return "", Error{
					Type:    "FAILED_VERIFICATION",
					Message: "Verification error",
				}
			}
		}

	}

	fmt.Println("--------------------------------------------------------------------------------------------------------- three")

	switch entity.TransactionType(txt.Type) {
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
							return "", Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						req.Header.Set("Content-Type", "application/json")

						res, err := client.Do(req)
						if err != nil {
							return "", Error{
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
							return "", Error{
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
							return "", Error{
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

						// uc.log.Println(txn.Reference)

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
							Amount:            txt.Amount,
							Reference:         txt.Reference,
							Narration:         "",
							AwashAccount:      sender.Detail.(entity.BankAccount).Number,
							CreditAccount:     "01320209107500",
							CommissionAmount:  0,
							AwashAccountName:  sender.Detail.(entity.BankAccount).Holder.Name,
							CreditAccountName: "LakiPay",
						})

						uc.log.Println(txt.Reference)

						uc.log.Println("Amhara Bank 15")
						if err != nil {
							uc.log.Println("Amhara Bank 16")
							uc.log.Println(err)
							return "", Error{
								Type:    "FAILED_TOVERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 17")
						req, err = http.NewRequest(http.MethodPost, "http://10.10.101.144:8080/b2b/awash/api/v1/monetize/post", bytes.NewBuffer(serBody))
						if err != nil {
							uc.log.Println("Amhara Bank 18")
							return "", Error{
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
							return "", Error{
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
							return "", Error{
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
						err = uc.repo.UpdateTransaction(parsedUUID)
						if err != nil {
							return "", Error{
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
							return "", Error{
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
							return "", Error{
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
							return "", Error{
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
							return "", Error{
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
								Id:     txt.Reference,
								Amount: txt.Amount,
							},
						})

						uc.log.Println("Amhara Bank 15")
						if err != nil {
							uc.log.Println("Amhara Bank 16")
							uc.log.Println(err)
							return "", Error{
								Type:    "FAILED_TOVERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Amhara Bank 17")
						req, err = http.NewRequest(http.MethodPost, "http://172.31.2.30:8600/abaApi/v1/lakiPay/processPayment", bytes.NewBuffer(serBody))
						if err != nil {
							uc.log.Println("Amhara Bank 18")
							return "", Error{
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
							return "", Error{
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
							return "", Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: txnRes.Message,
							}
						}

						uc.log.Println("Amhara Bank 23")
						// Store transaction
						txt.Verified = true
						err = uc.repo.UpdateTransaction(parsedUUID)
						if err != nil {
							return "", Error{
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
								Balance: receipient.Detail.(entity.StoredAccount).Balance + txt.Amount,
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
							Amount:          txt.Amount,
							ExplanationCode: "9904",
							Remark:          txt.Reference,
						})

						if err != nil {
							uc.log.Println("Oromia Bank 16")
							uc.log.Println(err)
							return "", Error{
								Type:    "FAILED_TOVERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Oromia Bank 17")
						req, err := http.NewRequest(http.MethodPost, "http://10.10.20.47/fund-transfer/customer-to-settlement", bytes.NewBuffer(serBody))
						if err != nil {
							uc.log.Println("Oromia Bank 18")
							return "", Error{
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
							return "", Error{
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
							return "", Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: txnRes.Message,
							}
						}

						uc.log.Println("Oromia Bank 23")
						// Store transaction
						txt.Verified = true
						err = uc.repo.UpdateTransaction(parsedUUID)
						if err != nil {
							return "", Error{
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
								Balance: receipient.Detail.(entity.StoredAccount).Balance + txt.Amount,
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
							return "", Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						req.Header.Set("Content-Type", "application/json")

						res, err := client.Do(req)
						if err != nil {
							return "", Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						if res.StatusCode != http.StatusOK {
							// Unsuccessful request
							uc.log.Println("Send Error")
							return "", Error{
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
							return "", Error{
								Type:    "NO_RESPONSE",
								Message: err.Error(),
							}
						}

						defer res.Body.Close()

						uc.log.Println(txt.Reference)

						serBody, err = json.Marshal(&struct {
							CreditAccount string                   `json:"credit_account"`
							DebitAccount  string                   `json:"debit_account"`
							Date          time.Time                `json:"date"`
							Amount        float64                  `json:"amount"`
							Payloads      []map[string]interface{} `json:"payloads"`
						}{
							CreditAccount: "01320209107500",
							DebitAccount:  sender.Detail.(entity.BankAccount).Number,
							Amount:        txt.Amount,
							Date:          txt.CreatedAt,
							Payloads: []map[string]interface{}{
								{
									"txn_ref": txt.Reference,
								},
							},
						})

						uc.log.Println(txt.Reference)

						uc.log.Println("Bunna Bank 15")
						if err != nil {
							uc.log.Println("Bunna Bank 16")
							uc.log.Println(err)
							return "", Error{
								Type:    "FAILED_TOVERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Bunna Bank 17")
						req, err = http.NewRequest(http.MethodPost, "http://10.1.13.12/api/core/transaction/open_c2c/initiate", bytes.NewBuffer(serBody))
						if err != nil {
							uc.log.Println("Bunna Bank 18")
							return "", Error{
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
							return "", Error{
								Type:    "FAILED_TO_VERIFY_TRANSACTION",
								Message: err.Error(),
							}
						}

						uc.log.Println("Bunna Bank 21")
						body, err := io.ReadAll(res.Body)
						if err != nil {
							return "", Error{
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
							return "", Error{
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

						txt.Reference = txnRes.ReferenceId

						uc.log.Println("Bunna Bank 23")
						// Store transaction
						err = uc.repo.UpdateTransaction(parsedUUID)
						if err != nil {
							return "", Error{
								Type:    "FAILED_TO_STORE_TRANSACTION",
								Message: err.Error(),
							}
						}
					}
				}

			}
		}

	case entity.P2P:
		{
			err = uc.repo.UpdateTransaction(parsedUUID)
			if err != nil {
				return "", Error{
					Type:    "FAILED_TO_STORE_TRANSACTION",
					Message: err.Error(),
				}
			}

		}
	case entity.SALE:
		{
			err = uc.repo.UpdateTransaction(parsedUUID)
			if err != nil {
				return "", Error{
					Type:    "FAILED_TO_STORE_TRANSACTION",
					Message: err.Error(),
				}
			}
		}
	case entity.SETTLEMENT:
		{
			err = uc.repo.UpdateTransaction(parsedUUID)
			if err != nil {
				return "", Error{
					Type:    "FAILED_TO_STORE_TRANSACTION",
					Message: err.Error(),
				}
			}
		}
	case entity.BILL:
		{
			type RefreshResponse struct {
				AccessToken string `json:"access_token"`
				ExpiresIn   int    `json:"expires_in"`
				// Add other fields if needed
			}
			baseURL := "https://auth.teleport.et/auth/realms/tp_merchant/protocol/openid-connect/token"
			// token := "eyJhbGciOiJIUzI1NiJ9.eyJpZGVudGlmaWVyIjoiV2dKVzRBUmtDM3ZNeE5xM3VwWTNkS2Y4V3hSRjhiam4iLCJleHAiOjE4NjM4MDM2MjIsImlhdCI6MTcwNTk1MDgyMiwianRpIjoiODNkYWQxNWUtODU4Ny00NTU0LTkxZWItZDYxYTU0NDNjYTkzIn0.xTzHDwpU9qkRurb0iCnixPFs3qanu3pk86L3hJTXZEI"

			client := &http.Client{}

			form := url.Values{}
			form.Add("grant_type", "client_credentials")

			req, err := http.NewRequest("POST", baseURL, bytes.NewBufferString(form.Encode()))
			if err != nil {
				return "", Error{
					Type:    "FAILED",
					Message: err.Error(),
				}
			}
			req.Header.Set("Authorization", "Basic VHJ1c3RlZDpHb2dUYkJwSTNBTldEZWFteEVXdmFocWdKVGlTYTlsMw==")
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req)
			if err != nil {
				return "", Error{
					Type:    "FAILED",
					Message: err.Error(),
				}
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", Error{
					Type:    "FAILED",
					Message: err.Error(),
				}
			}

			var refreshResponse RefreshResponse
			err = json.Unmarshal(body, &refreshResponse)
			if err != nil {
				return "", Error{
					Type:    "FAILED_TO_STORE_TRANSACTION",
					Message: err.Error(),
				}
			}
			// Process the response
			fmt.Print("{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{{}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}")
			// fmt.Println("Response:", string(body))

			// fmt.Println("|||||| puplics: ",refreshResponse.AccessToken,)

			baseURL2 := "https://api.teleport.et/api/airtime-topup"

			client2 := &http.Client{}

			req2, err := http.NewRequest("POST", baseURL2, bytes.NewBuffer([]byte(fmt.Sprintf(`{"msisdn": %s,"topupType": "PREPAID","amount": %v}`, txt.Phone, 2))))
			if err != nil {
				return "", Error{
					Type:    "FAILED",
					Message: err.Error(),
				}
			}

			req2.Header.Set("Authorization", "Bearer "+refreshResponse.AccessToken)
			req2.Header.Set("Content-Type", "application/json")

			resp2, err := client2.Do(req2)
			if err != nil {
				return "", Error{
					Type:    "FAILED REQ 2",
					Message: err.Error(),
				}
			}

			defer resp2.Body.Close()
			body3, err := ioutil.ReadAll(resp2.Body)
			if err != nil {
				return "", Error{
					Type:    "FAILED",
					Message: err.Error(),
				}
			}
			fmt.Println(string(body3))
			// type Res struct{
			// 	Code :
			// }
			fmt.Println("=============")
			fmt.Println(resp2.StatusCode)
			if resp2.StatusCode == 200 {
				err = uc.repo.UpdateTransaction(parsedUUID)
				if err != nil {
					return "", Error{
						Type:    "FAILED_TO_STORE_TRANSACTION",
						Message: err.Error(),
					}
				}
			} else {

				return "", Error{
					Type:    "FAILED",
					Message: "Something is wrong",
				}

			}

		}
	}

	return "", nil

}

func (uc Usecase) CreateTransaction(userId uuid.UUID, from uuid.UUID, to uuid.UUID, amount float64, txnType string, token string, challenge_type string, challenge entity.TransactionChallange) (*entity.Transaction, error) {
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
	} else {
		receipient = &entity.Account{}
	}

	uc.log.Println(entity.TransactionType(txnType))

	switch challenge_type {

	case "2fa":
		{
			log1 := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

			// [Ouput Adapters]
			// [DB] Postgres
			db, err := db_psql.New(log1)
			if err != nil {
				log.Fatal(err.Error())
			}
			r, err := auth_repo.NewPsqlRepo(uc.log, db)
			if err != nil {
				return nil, Error{
					Type:    "FAILED_TO_CHEKE_2FA",
					Message: err.Error(),
				}
			}
			pass, err := r.FindPasswordIdentityByUser(userId)
			if err != nil {
				return nil, Error{
					Type:    "FAILED_TO_CHEKE_2FA11",
					Message: err.Error(),
				}
			}

			if pass == nil {
				return nil, Error{
					Type:    "FAILED_TO_CHEKE_2FA22",
					Message: err.Error(),
				}
			}

			hasher := sha256.New()
			_, err = hasher.Write([]byte(challenge.TwoFA))
			if err != nil {
				return nil, Error{
					Type:    "FAILED_TO_CHEKE_2FA33",
					Message: err.Error(),
				}
			}

			// Compare passwords

			if base64.URLEncoding.EncodeToString(hasher.Sum(nil)) != pass.Password {
				return nil, Error{
					Type:    "FAILED_TO_CHEKE_2FA",
					Message: "2fa faild",
				}
			}
		}
	case "otp":
		{
			log1 := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

			// [Ouput Adapters]
			// [DB] Postgres
			db, err := db_psql.New(log1)
			if err != nil {
				log.Fatal(err.Error())
			}
			r3, err := auth_repo.NewPsqlRepo(uc.log, db)
			if err != nil {
				return nil, Error{
					Type:    "FAILED_TO_CHEKE_2FA",
					Message: err.Error(),
				}
			}
			phoneAuth, err := r3.FindPhoneAuthWithoutPhone(userId.String())
			if err != nil {
				log.Println("1")
				return nil, Error{
					Type:    "FAILED_TO_AUTH_PHONE",
					Message: err.Error(),
				}
			}
			// // Validate phone
			// if phoneAuth.Phone.Prefix != prefix || phoneAuth.Phone.Number != number {
			// 	log.Println("2")
			// 	return nil, Error{
			// 		Type:    "FAILED_TO_AUTH_PHONE",
			// 		Message: "Phone is not valid",
			// 	}
			// }

			// Check the status
			if phoneAuth.Status {
				return nil, Error{
					Type:    "",
					Message: "already used",
				}
			}

			_, err = jwt.Decode(phoneAuth.Code, challenge.OTP)
			if err != nil {
				if err.Error() == "invalid token" {
					return nil, Error{
						Type:    "",
						Message: "Incorrect OTP",
					}
				}
				return nil, Error{
					Type:    "",
					Message: err.Error(),
				}
			}
			err = r3.UpdatePhoneAuthStatus(phoneAuth.Id, true)
			if err != nil {
				log.Println("5")
				return nil, Error{
					Type:    "",
					Message: err.Error(),
				}
			}

		}
	case "finger_print":
		{
			public_keys, err := uc.repo.GetPuplicKey(challenge.Challenge, userId)

			fmt.Println("|||||||||||| ====== public key", public_keys)
			if err != nil {

				return nil, Error{
					Type:    "FAILED_TO_GET_PUBLIC_KEY",
					Message: err.Error(),
				}
			}

			var isVerify bool = false
			for _, pk := range public_keys {
				if verify(pk.PublicKey, challenge.Challenge, challenge.Signature) {
					err = uc.repo.UpdatePublicKeysUsed(userId)

					if err != nil {
						return nil, Error{
							Type:    "FAILED_VERIFICATION",
							Message: err.Error(),
						}
					}

					isVerify = true
				}
			}

			if !isVerify {
				return nil, Error{
					Type:    "FAILED_VERIFICATION",
					Message: "Verification error",
				}
			}
		}

	}

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
		{
			fmt.Println("||||||||||||||||||||| p2p ||||||||||||||||||||||||")
			id := uuid.New()
			txn = entity.Transaction{
				Id:        id,
				From:      *sender,
				To:        *receipient,
				Type:      entity.TransactionType(txnType),
				Verified:  false,
				Amount:    amount,
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
	case entity.SALE:
		{
			fmt.Println("||||||||||||||||||||| sale ||||||||||||||||||||||||")
			id := uuid.New()
			txn = entity.Transaction{
				Id:        id,
				From:      *sender,
				To:        *receipient,
				Type:      entity.TransactionType(txnType),
				Verified:  false,
				CreatedAt: time.Now(),
				Amount:    amount,
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
	case entity.SETTLEMENT:
		{
			fmt.Println("||||||||||||||||||||| Settlement ||||||||||||||||||||||||", receipient)
			id := uuid.New()
			txn = entity.Transaction{
				Id:        id,
				From:      *sender,
				To:        *receipient,
				Type:      entity.TransactionType(txnType),
				Verified:  false,
				CreatedAt: time.Now(),
				Amount:    amount,
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
	}

	return &txn, nil
}

func (uc Usecase) VerifyTransaction2(userId, txnId uuid.UUID, code string, amount float64) (*entity.Transaction, error) {
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

func (uc Usecase) SendOtpUsecase(userId uuid.UUID) (string, error) {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	random_number := 100000 + r.Intn(900000)

	var phoneAuth auth_entity.PhoneAuth
	id := uuid.New()

	phoneAuth = auth_entity.PhoneAuth{
		Id:      id,
		Token:   userId.String(),
		Phone:   auth_entity.Phone{},
		Method:  "sms",
		Length:  int64(6),
		Timeout: int64(120),
		Code: jwt.Encode(jwt.Payload{
			Exp: time.Now().Unix() + 30*60,
		}, fmt.Sprint(random_number)),
	}
	log1 := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	// [Ouput Adapters]
	// [DB] Postgres
	db, err := db_psql.New(log1)
	if err != nil {
		log.Fatal(err.Error())
	}
	r2, err := auth_repo.NewPsqlRepo(uc.log, db)
	if err != nil {
		return "", Error{
			Type:    "FAILED_TO_CHEKE_2FA",
			Message: err.Error(),
		}
	}
	err = r2.StorePhoneAuth(phoneAuth)
	if err != nil {
		return "", nil
	}
	go func() {
		utils.SendSMS(random_number, "+251986680094")
	}()

	return "", nil
}

func (uc Usecase) SendSetFIngerPrintUsecase(userId uuid.UUID, data interface{}) (string, error) {

	dataBytes, err := json.Marshal(data)
	if err != nil {

		return "", err
	}

	fmt.Println("||||||||||||||||||||||||||||||||||| byte ", dataBytes)

	// Hash the JSON data using SHA-256
	hash := sha256.Sum256(dataBytes)
	hashString := fmt.Sprintf("%x", hash[:])

	fmt.Println("||||||||||||||||||||||||||||||||||| ", data)

	log1 := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	// [Ouput Adapters]
	// [DB] Postgres
	db, err := db_psql.New(log1)
	if err != nil {
		log.Fatal(err.Error())
	}
	r2, err := auth_repo.NewPsqlRepo(uc.log, db)
	if err != nil {
		return "", Error{
			Type:    "FAILED_TO_CHEKE_2FA",
			Message: err.Error(),
		}
	}
	err = r2.UpdatePasswordIdentity(hashString, userId)
	if err != nil {
		return "", err
	}
	// Store the hash in PostgreSQL
	//   err = storeHashInDB(hashString)
	//   if err != nil {
	// 	  return "" ,err
	//   }

	return "true", nil
}

func (uc Usecase) SendGenerateChallenge(userID uuid.UUID, deviceId string) (string, error) {
	challengeBytes := make([]byte, 16)
	_, err := rand.Read(challengeBytes)
	if err != nil {
		return "", err
	}
	challenge := base64.StdEncoding.EncodeToString(challengeBytes)

	err = uc.repo.UpdateGeneratedChallenge(challenge, userID, deviceId)
	// _, err = db.Exec("UPDATE public_keys SET challenge=$1, expires_at=$2, used=FALSE WHERE user_id=$3", challenge, expiresAt, userID)
	if err != nil {
		return "", err
	}

	return challenge, nil
}

func (uc Usecase) GetverifySignature(userID uuid.UUID, challenge string, sign string) (string, error) {

	public_keys, err := uc.repo.GetPuplicKey(challenge, userID)

	fmt.Println("|||||||||||| ====== ", public_keys)

	if err != nil {
		return "", err
	}

	for _, pk := range public_keys {
		if verify(pk.PublicKey, challenge, sign) {
			err = uc.repo.UpdatePublicKeysUsed(userID)

			if err != nil {
				return "", err
			}

			return "signature verification successed!!", nil
		}
	}

	// err=uc.repo.UpdateGeneratedChallenge(challenge,userID)
	// // _, err = db.Exec("UPDATE public_keys SET challenge=$1, expires_at=$2, used=FALSE WHERE user_id=$3", challenge, expiresAt, userID)
	// if err != nil {
	//     return "", err
	// }

	return "signature verification failed", nil
}

func verify(publicKeyPEM, challenge, signature string) bool {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	fmt.Println(block)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return false
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return false
	}

	hash := sha256.New()
	hash.Write([]byte(challenge))
	hashed := hash.Sum(nil)

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed, sigBytes) == nil
}

// func verify(challenge string, signature []byte, publicKeyPem string) bool {
//     block, _ := pem.Decode([]byte(publicKeyPem))
//     if block == nil || block.Type != "RSA PUBLIC KEY" {
//         return false
//     }

//     pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
//     if err != nil {
//         return false
//     }

//     hashed := sha256.Sum256([]byte(challenge))
//     return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature) == nil
// }

func (uc Usecase) GetstorePublicKeyHandler(key string, id uuid.UUID, device string) (string, error) {

	err := uc.repo.GetstorePublicKeyHandler(key, id, device)

	if err != nil {
		return "", err
	}

	return "Public key have stored!", nil
}

func (uc Usecase) GetfetchPublicKeyHandler(id uuid.UUID) (string, error) {

	res, err := uc.repo.GetPuplicKey("", id)

	if err != nil {
		return "", err
	}

	return res[0].PublicKey, nil
}

func generateRSAKeyPair() (string, string, error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand2.Reader, 1024)
	if err != nil {
		return "", "", err
	}

	// Encode the private key to PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Encode the public key to PEM format
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}
