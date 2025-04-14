package psql

import (
	"auth/src/pkg/checkout/core/entity"
	"encoding/json"
)

func (repo CheckoutPSQLRepo) StoreTransaction(v entity.Transaction) error {
	details, _ := json.Marshal(v.Details)
	pricing, _ := json.Marshal(v.Pricing)
	status, _ := json.Marshal(v.Status)

	_, err := repo.db.Exec(`
	INSERT INTO checkout.transactions
	("id", "to", "for", "ttl", "pricing", "status", "gateway", "type", "details", "created_at")
	VALUES
	($1,$2,$3,$4,$5,$6,$7,$8, $9, $10);
	`, v.Id, v.To, v.For, v.Ttl, string(pricing), string(status), v.GateWay, v.Type, string(details), v.CreatedAt)

	return err
}

func (repo CheckoutPSQLRepo) FindTransaction(v string) (*entity.Transaction, error) {
	var txn entity.Transaction

	var details string
	var pricing string
	var status string

	err := repo.db.QueryRow(`
	SELECT 
	"id", "to", "for", "ttl", "pricing", "status", "gateway", "type", "details", "created_at", "updated_at"
	FROM checkout.transactions
	WHERE "id" = $1;
	`, v).Scan(
		&txn.Id,
		&txn.To,
		&txn.For,
		&txn.Ttl,
		&pricing,
		&status,
		&txn.GateWay,
		&txn.Type,
		&details,
		&txn.CreatedAt,
		&txn.UpdatedAt,
	)

	json.Unmarshal([]byte(details), &txn.Details)
	json.Unmarshal([]byte(pricing), &txn.Pricing)
	json.Unmarshal([]byte(status), &txn.Status)

	repo.log.Println("txn.GateWay")
	repo.log.Println(txn.GateWay)

	return &txn, err
}

func (repo CheckoutPSQLRepo) UpdateTransaction(v entity.Transaction) error {

	status, _ := json.Marshal(v.Status)

	_, err := repo.db.Exec(`
	UPDATE checkout.transactions
	SET status = $1
	WHERE id = $2;
	`, status, v.Id)

	return err
}
