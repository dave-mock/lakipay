package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	REPLENISHMENT TransactionType = "REPLENISHMENT"
	P2P           TransactionType = "P2P"
	SALE          TransactionType = "SALE"
	BILL_PAYMENT  TransactionType = "BILL_PAYMENT"
	SETTLEMENT    TransactionType = "SETTLEMENT"
	BILL          TransactionType = "BILL_AIRTIME"
)

type TransactionMedium string

const (
	LAKIPAY     TransactionMedium = "LAKIPAY"
	CYBERSOURCE TransactionMedium = "CYBERSOURCE"
	ETHSWITCH   TransactionMedium = "ETHSWITCH"
)

type Transaction struct {
	Id         uuid.UUID
	From       Account
	To         Account
	Type       TransactionType
	Medium     TransactionMedium
	Reference  string
	Comment    string
	Tag        Tag
	Verified   bool
	TTL        int64
	Commission float64
	Details    interface{}
	CreatedAt  time.Time
	UpdatedAt  time.Time
	// new
	ErrorMessage      string
	Confirm_Timestamp time.Time
	BankReference     string
	PaymentMethod     string
	Test              bool
	Status            string
	Description       string
	Token             string
	Amount            float64
	HasChallenge      bool
	TotalAmount       float64
	Currency          string
	Phone             string
}

type TransactionChallange struct {
	TwoFA     string `json:"2fa"`
	Challenge string `json:"challenge"`
	OTP       string `json:"otp"`
	Signature string `json:"signature"`
}

type Replenishment struct {
	Amount float64
}

type MerchantKeys struct {
	Id         uuid.UUID
	PublickKey string
	PrivateKey string
	MerchantId uuid.UUID
	Username   string
	Password   string
}

type P2p struct {
	Amount float64
}

type BatchTransaction struct {
	Id   uuid.UUID
	From []struct {
		Account Account
		Amount  float64
	}
	To []struct {
		Account Account
		Amount  float64
	}
	Amount       float64
	Transactions []Transaction
	Verified     bool
	TTL          int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PublicKey struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	PublicKey string    `db:"public_key"`
	DeviceID  string    `db:"device_id"`
	Challenge string    `db:"challenge"`
	ExpiresAt time.Time `db:"expires_at"`
	Used      bool      `db:"used"`
	CreatedAt time.Time `db:"created_at"`
}

type TransactionSession struct {
	Id        uuid.UUID
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
