package usecase

import (
	"auth/src/pkg/account/core/entity"

	"github.com/google/uuid"
)

type Repo interface {
	UpdateUserRepo(users entity.User2) (entity.User2, error)
	StoreTransactionSession(preSession entity.TransactionSession) error
	UpdateTransaction(id uuid.UUID) error
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
	UpdateGeneratedChallenge(challenge string, id uuid.UUID, deviceId string) error
	UpdatePublicKeysUsed(id uuid.UUID) error

	GetPuplicKey(challenge string, id uuid.UUID) ([]*entity.PublicKey, error)
	GetstorePublicKeyHandler(key string, id uuid.UUID, device string) error

	StoreKeys(id uuid.UUID, string, privateKey string, password string, username string) error
	GetMerchantsKeys(secret_key string) (entity.MerchantKeys, error)
	CheckMerchantsKeysByUsername(username string) (error, bool)

	GetApiKeysRepo(id uuid.UUID) (entity.MerchantKeys, error)
}
