package usecase

import (
	"auth/src/pkg/account/core/entity"

	"github.com/google/uuid"
)

type Interactor interface {
	// Bank
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

	CreateTransaction(userId uuid.UUID, from uuid.UUID, to uuid.UUID, amount float64, txnType string, token string) (*entity.Transaction, error)
	VerifyTransaction(userId, txnId uuid.UUID, code string, amount float64) (*entity.Transaction, error)
	GetUserTransactions(id uuid.UUID) ([]entity.Transaction, error)
	GetAllTransactions() ([]entity.Transaction, error)
	TransactionsDashboardUsecase(year int) (interface{}, error)
	// GetUserTransactions(userId uuid.UUID) ([]entity.Transaction, error)

	// Verify Bank Account
	VerifyAccount(userId, accountId uuid.UUID, method string, details interface{}, code string) (string, error)
	DeleteAccount(userId, accId uuid.UUID) error
}
