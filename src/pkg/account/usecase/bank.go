package usecase

import (
	"auth/src/pkg/account/core/entity"

	"github.com/google/uuid"
)

func (uc Usecase) AddBank(name, shortName, bin, swiftCode, logo string) (*entity.Bank, error) {
	// Errors
	var (
		ErrFailedToStoreBank string = "FAILED_TO_STORE_BANK"
	)

	// [TODO] - Validate inputs

	var bank *entity.Bank

	// Create bank
	id := uuid.New()
	uc.log.Println("CREATING BANK")
	bank = &entity.Bank{
		Id:        id,
		Name:      name,
		ShortName: shortName,
		BIN:       bin,
		SwiftCode: swiftCode,
		Logo:      logo,
	}

	uc.log.Println("CREATED BANK")

	// Store bank
	err := uc.repo.StoreBank(*bank)
	uc.log.Println("STORE BANK")
	if err != nil {
		uc.log.Println("STORE BANK ERROR")
		return nil, Error{
			Type:    ErrFailedToStoreBank,
			Message: err.Error(),
		}
	}

	uc.log.Println("STORED BANK")
	return bank, nil
}

func (uc Usecase) GetBanks() ([]entity.Bank, error) {
	// Errors
	var (
		ErrCouldNotFindBanks string = "COULD_NOT_FIND_BANKS"
	)
	var banks []entity.Bank

	banks, err := uc.repo.FindBanks()
	if err != nil {
		return nil, Error{
			Type:    ErrCouldNotFindBanks,
			Message: err.Error(),
		}
	}

	return banks, nil
}

// func (uc Usecase) VerifyBankAccount(bankId uuid.UUID, accountNumber string, accountHolderName string, accountHolderPhone string) error {

// 	// Find bank
// 	_, err := uc.repo.FindBankById(bankId)
// 	if err != nil {
// 		return Error{
// 			Type:    "BANK_NOT_FOUND",
// 			Message: err.Error(),
// 		}
// 	}

// 	// Find Bank Account Verification Data
// 	bankAccountVerification, err := uc.repo.FindBankAccountVerification(bankId)
// 	if err != nil {
// 		return Error{
// 			Type:    "CURRENT_BANK_DOESNOT_SUPPORT_ACCOUNT_VERIFICATION",
// 			Message: err.Error(),
// 		}
// 	}

// 	// Authorize request

// 	// Send Verify Account Request
// 	var netTransport = &http.Transport{
// 		Dial: (&net.Dialer{
// 			Timeout: 1 * time.Minute,
// 		}).Dial,
// 		TLSHandshakeTimeout: 1 * time.Minute,
// 	}

// 	var client = &http.Client{
// 		Timeout:   time.Minute * 1,
// 		Transport: netTransport,
// 	}

// 	queryParams := url.Values{
// 		"bankCode": {"awsbnk"},
// 		"account":  {accountNumber},
// 	}

// 	req, err := http.NewRequest(http.MethodGet, bankAccountVerification.Url+"?"+queryParams.Encode(), nil)
// 	if err != nil {
// 		return Error{
// 			Type:    "NO_RESPONSE",
// 			Message: err.Error(),
// 		}
// 	}

// 	// Authorize request

// 	// Find Bank Client Auths
// 	clientAuth, err := uc.repo.FindBankClientAuth(bankId)
// 	if err != nil {
// 		return Error{
// 			Type:    "NO_CLIENT_AUTH_FOUND",
// 			Message: err.Error(),
// 		}
// 	}

// 	// Auth Req
// 	// uc.log.Println(clientAuth.Type)
// 	// uc.log.Println(clientAuth.Authorizer.(*oAuth2.OAuth2))
// 	clientAuth.Authorizer.AuthorizeHTTP(req)

// 	uc.log.Println(req.Header.Get("Authorization"))

// 	res, err := client.Do(req)
// 	if err != nil {
// 		return Error{
// 			Type:    "NO_RESPONSE",
// 			Message: err.Error(),
// 		}
// 	}

// 	type Res struct {
// 		Status           int    `json:"status"`
// 		ErrorDescription string `json:"errorDescription"`
// 	}

// 	var _res Res

// 	decoder := json.NewDecoder(res.Body)
// 	err = decoder.Decode(&_res)
// 	if err != nil {
// 		return Error{
// 			Type:    "UNDEFINED",
// 			Message: err.Error(),
// 		}
// 	}

// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK || _res.Status == 0 {
// 		// Unsuccessful request
// 		return Error{
// 			Type:    "NO_RESPONSE",
// 			Message: _res.ErrorDescription,
// 		}
// 	}

// 	return nil
// }
