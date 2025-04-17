package usecase

import (
	"auth/src/pkg/account/core/entity"

	"github.com/google/uuid"
)

type Interactor interface {
	// Bank
	UpdateUserUsecase(users entity.User2) (entity.User2, error)
	AddBank(name, shortName, bin, swiftCode, logo string) (*entity.Bank, error)
	GetBanks() ([]entity.Bank, error)
	/// Accounts
	GetUserAccounts(id uuid.UUID) ([]entity.Account, error)
	// Stored Account
	CreateStoredAccount(userId uuid.UUID, title string, isDefault bool) (*entity.Account, error)
	// Bank Account
	CreateBankAccount(userId uuid.UUID, bankId uuid.UUID, accountNumber string, accountHolderName string, accountHolderPhone string, title string, makeDefault bool) (*entity.Account, error)

	// Transaction
	// InitTransaction(
	// 	from uuid.UUID,
	// 	to []struct {
	// 		Account uuid.UUID
	// 		Ratio   float64
	// 		Amount  float64
	// 	},
	// 	medium entity.TransactionMedium,
	// 	txType entity.TransactionType,
	// 	amount float64,
	// )

	CreateRegisterKeys(id uuid.UUID, username string, password string) (string, error)
	CreateHostedTransactionInitiate(amount float64, currency string, callback_url string, secretKey string, stringData string, token string) (interface{}, error)

	CreateTransactionInitiate(userId uuid.UUID, from uuid.UUID, to uuid.UUID, amount float64, medium entity.TransactionMedium, txnType string, token string, detail string) (interface{}, error)
	CreateTransaction(userId uuid.UUID, from uuid.UUID, to uuid.UUID, amount float64, txnType string, token string, challenge_type string, challenge entity.TransactionChallange) (*entity.Transaction, error)

	// VerifyTransaction(userId, txnId uuid.UUID, code string, amount float64) (*entity.Transaction, error)
	GetUserTransactions(id uuid.UUID) ([]entity.Transaction, error)
	GetAllTransactions() ([]entity.Transaction, error)
	TransactionsDashboardUsecase(year int) (interface{}, error)

	// GetUserTransactions(userId uuid.UUID) ([]entity.Transaction, error)

	// Verify Bank Account
	VerifyAccount(userId, accountId uuid.UUID, method string, details interface{}, code string) (string, error)
	DeleteAccount(userId, accId uuid.UUID) error

	SendOtpUsecase(userId uuid.UUID) (string, error)
	SendSetFIngerPrintUsecase(userId uuid.UUID, data interface{}) (string, error)
	SendGenerateChallenge(id uuid.UUID, deviceId string) (string, error)
	GetverifySignature(id uuid.UUID, challenge string, sign string) (string, error)
	GetstorePublicKeyHandler(key string, id uuid.UUID, device string) (string, error)
	GetfetchPublicKeyHandler(id uuid.UUID) (string, error)

	InitPreSession(txtId uuid.UUID) (entity.TransactionSession, error)
	VerifyTransaction(UserId uuid.UUID, token string, transactionChallenges entity.TransactionChallange, challengeType string) (string, error)

	GetApiKeysUsecase(id uuid.UUID) (string, error)
	ApplyForTokenUsecase(username string, password string) (string, error)
	CheckBalance(from uuid.UUID) (float64, error)
}
