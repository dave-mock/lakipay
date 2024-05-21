package usecase

import (
	"auth/src/pkg/account/core/entity"

	"github.com/google/uuid"
)

type Repo interface {
	// Account
	StoreAccount(acc entity.Account) error
	UpdateAccount(acc entity.Account) error
	FindAccountById(accId uuid.UUID) (*entity.Account, error)
	FindAccountsByUserId(userId uuid.UUID) ([]entity.Account, error)
	// Bank
	StoreBank(bank entity.Bank) error
	FindBanks() ([]entity.Bank, error)
	FindBankById(bankId uuid.UUID) (*entity.Bank, error)
	DeleteAccount(accId uuid.UUID) error

	// Transaction
	StoreTransaction(entity.Transaction) error
	FindTransactionById(id uuid.UUID) (*entity.Transaction, error)
	FindTransactionsByUserId(id uuid.UUID) ([]entity.Transaction, error)
	FindAllTransactions() ([]entity.Transaction, error)

	TransactionsDashboardRepo(year int) (interface{}, error)
}
