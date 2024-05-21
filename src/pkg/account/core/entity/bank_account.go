package entity

type BankAccount struct {
	Bank   Bank
	Number string
	Holder struct {
		Name  string
		Phone string
	}
}

func (acc BankAccount) Credit(from Account, amount float64) {
	// return nil
}

func (acc BankAccount) Debit(to Account, amount float64) {
	// return nil
}
