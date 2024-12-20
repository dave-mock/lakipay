package entity

type StoredAccount struct {
	Balance float64
}

func (acc StoredAccount) Credit(from Account, amount float64) {
	// return nil
}

func (acc StoredAccount) Debit(to Account, amount float64) {
	// return nil
}
