package psql

import (
	"auth/src/pkg/account/core/entity"
	"auth/src/pkg/utils"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (repo PsqlRepo) StoreTransaction(txn entity.Transaction) error {
	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	// data, err := utils.AesEncrption("marsal")
	// data, err := utils.AesDecription("2BvHg7Fq4TSwdC8diew2WA==")

	var prev_txn entity.Transaction
	err = repo.db.QueryRow(`select id,created_at from accounts.transactions order by created_at desc limit 1`).Scan(&prev_txn.Id, &prev_txn.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Print(prev_txn)
	token_str := prev_txn.Id.String() + prev_txn.CreatedAt.String()
	token, err := utils.AesEncrption(token_str)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
	INSERT INTO accounts.transactions (id, "from", "to", "type", "reference", verified, created_at,token,medium)
	VALUES ($1::UUID, $2::UUID, $3::UUID, $4, $5, $6, $7,$8,$9)
	`, txn.Id, txn.From.Id, txn.To.Id, txn.Type, txn.Reference, txn.Verified, txn.CreatedAt, token, txn.Medium)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Store transaction details
	switch txn.Type {
	case entity.REPLENISHMENT:
		{
			txnDetail := txn.Details.(entity.Replenishment)
			_, err = tx.Exec(`
			INSERT INTO accounts.a2a_transactions (transaction_id, amount)
			VALUES ($1::UUID,$2)
			`, txn.Id, txnDetail.Amount)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	case entity.P2P:
		{
			txnDetail := txn.Details.(entity.P2p)
			// amount:= utils.AesEncrption(txnDetail.Amount)
			_, err = tx.Exec(`
			INSERT INTO accounts.p2p_transactions (transaction_id, amount)
			VALUES ($1::UUID,$2)
			`, txn.Id, txnDetail.Amount)
			if err != nil {
				tx.Rollback()
				return err
			}

		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return err
}

func (repo PsqlRepo) FindAllTransactions() ([]entity.Transaction, error) {
	var txns []entity.Transaction = make([]entity.Transaction, 0)
	var txnsResult []entity.Transaction = make([]entity.Transaction, 0)

	// rows, err := repo.db.Query(`
	// SELECT
	// 	transactions.id, transactions.type, transactions.created_at, transactions.updated_at,
	// 	"tag".id, "tag".name, "tag".color,
	// 	"from".id, "from".title, "from".type, "from".default, "from".user,
	// 	"to".id, "to".title, "to".type, "to".default, "to".user
	// FROM accounts.transactions
	// LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
	// LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to
	// LEFT JOIN accounts.tags as "tag" ON "tag".id = transactions.tag
	// WHERE "from".user = $1::UUID OR "to".user = $1::UUID;
	// `, id)

	var txtDetiail []uint8
	rows, err := repo.db.Query(`
	SELECT 
		transactions.id, transactions.type, transactions.created_at, transactions.updated_at,
		medium, comment, accounts.transactions.verified, reference,ttl, commission,details,error_message,confirm_timestamp,bank_reference,
		payment_method,test,description,
		"from".id, "from".title, "from".type, "from".default, "from".user_id,
		"to".id, "to".title, "to".type, "to".default, "to".user_id,transactions.id,
	FROM accounts.transactions
	LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
	LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to;
	`)

	if err != nil {
		repo.log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		fmt.Println("||||||||||||||||||||| transcation ")
		var txn entity.Transaction
		var medium sql.NullString
		var comment sql.NullString
		var commistion sql.NullFloat64
		var ttl sql.NullInt64
		var errMsg sql.NullString
		var time sql.NullTime
		var BankReference sql.NullString
		var pymentMethod sql.NullString
		var test sql.NullBool
		var description sql.NullString

		err := rows.Scan(&txn.Id, &txn.Type, &txn.CreatedAt, &txn.UpdatedAt,
			&medium, &comment, &txn.Verified, &txn.Reference, &ttl, &commistion, &txtDetiail, &errMsg, &time, &BankReference,
			&pymentMethod, &test, &description,
			&txn.From.Id, &txn.From.Title, &txn.From.Type, &txn.From.Default, &txn.From.User.Id,
			&txn.To.Id, &txn.To.Title, &txn.To.Type, &txn.To.Default, &txn.To.User.Id,
		)

		if medium.Valid {
			txn.Medium = entity.TransactionMedium(medium.String)
		}
		if comment.Valid {
			txn.Comment = comment.String
		}
		if commistion.Valid {
			txn.Commission = commistion.Float64
		}
		if ttl.Valid {
			txn.TTL = ttl.Int64
		}
		if errMsg.Valid {
			txn.ErrorMessage = errMsg.String
		}
		if time.Valid {
			txn.Confirm_Timestamp = time.Time
		}
		if BankReference.Valid {
			txn.BankReference = BankReference.String
		}
		if test.Valid {
			txn.Test = test.Bool
		}
		if description.Valid {
			txn.Description = description.String
		}
		if pymentMethod.Valid {
			txn.PaymentMethod = pymentMethod.String
		}

		if err != nil {
			// Fetch txn details
			return nil, err

		}
		json.Unmarshal(txtDetiail, &txn.Details)

		switch txn.Type {
		case entity.REPLENISHMENT:
			{
				txn.Details = nil
			}
		}
		fmt.Println("|||||||  ", txn)

		txns = append(txns, txn)
	}

	for index, element := range txns {
		if index == 0 {
			txnsResult = append(txnsResult, element)

		} else {
			i := index - 1
			fmt.Println("*********************************************************************************** ", i)

			targettoken := txns[i].Id.String() + txns[i].CreatedAt.String()
			fmt.Println("*********************************************************************************** ", i)

			// 		token_str := prev_txn.Id.String() + prev_txn.CreatedAt.String()
			token, err := utils.AesDecription(element.Token)
			if err != nil {
				// Fetch txn details
				return nil, err

			}
			fmt.Println("*********************************************************************************** check ")

			if token == targettoken {
				txnsResult = append(txnsResult, element)
			}
		}
	}

	return txnsResult, nil
}
func (repo PsqlRepo) FindTransactionsByUserId(id uuid.UUID) ([]entity.Transaction, error) {
	var txns []entity.Transaction = make([]entity.Transaction, 0)

	// rows, err := repo.db.Query(`
	// SELECT
	// 	transactions.id, transactions.type, transactions.created_at, transactions.updated_at,
	// 	"tag".id, "tag".name, "tag".color,
	// 	"from".id, "from".title, "from".type, "from".default, "from".user,
	// 	"to".id, "to".title, "to".type, "to".default, "to".user
	// FROM accounts.transactions
	// LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
	// LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to
	// LEFT JOIN accounts.tags as "tag" ON "tag".id = transactions.tag
	// WHERE "from".user = $1::UUID OR "to".user = $1::UUID;
	// `, id)

	var txtDetiail []uint8
	rows, err := repo.db.Query(`
	SELECT 
		transactions.id, transactions.type, transactions.created_at, transactions.updated_at,
		medium, comment, accounts.transactions.verified, reference,ttl, commission,details,error_message,confirm_timestamp,bank_reference,
		payment_method,test,description,
		"from".id, "from".title, "from".type, "from".default, "from".user_id,
		"to".id, "to".title, "to".type, "to".default, "to".user_id
	FROM accounts.transactions
	LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
	LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to
	WHERE "from".id = $1::UUID OR "to".id = $1::UUID;
	`, id)

	if err != nil {
		repo.log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var txn entity.Transaction
		var medium sql.NullString
		var comment sql.NullString
		var commistion sql.NullFloat64
		var ttl sql.NullInt64
		var errMsg sql.NullString
		var time sql.NullTime
		var BankReference sql.NullString
		var pymentMethod sql.NullString
		var test sql.NullBool
		var description sql.NullString

		err := rows.Scan(&txn.Id, &txn.Type, &txn.CreatedAt, &txn.UpdatedAt,
			&medium, &comment, &txn.Verified, &txn.Reference, &ttl, &commistion, &txtDetiail, &errMsg, &time, &BankReference,
			&pymentMethod, &test, &description,
			&txn.From.Id, &txn.From.Title, &txn.From.Type, &txn.From.Default, &txn.From.User.Id,
			&txn.To.Id, &txn.To.Title, &txn.To.Type, &txn.To.Default, &txn.To.User.Id,
		)

		if medium.Valid {
			txn.Medium = entity.TransactionMedium(medium.String)
		}
		if comment.Valid {
			txn.Comment = comment.String
		}
		if commistion.Valid {
			txn.Commission = commistion.Float64
		}
		if ttl.Valid {
			txn.TTL = ttl.Int64
		}
		if errMsg.Valid {
			txn.ErrorMessage = errMsg.String
		}
		if time.Valid {
			txn.Confirm_Timestamp = time.Time
		}
		if BankReference.Valid {
			txn.BankReference = BankReference.String
		}
		if test.Valid {
			txn.Test = test.Bool
		}
		if description.Valid {
			txn.Description = description.String
		}
		if pymentMethod.Valid {
			txn.PaymentMethod = pymentMethod.String
		}

		if err == nil {
			// Fetch txn details
			json.Unmarshal(txtDetiail, &txn.Details)
			switch txn.Type {
			case entity.REPLENISHMENT:
				{
					txn.Details = nil
				}
			}

			txns = append(txns, txn)
		}
	}

	return txns, nil
}

func (repo PsqlRepo) FindTransactionById(id uuid.UUID) (*entity.Transaction, error) {
	var txn entity.Transaction

	err := repo.db.QueryRow(`
	SELECT 
		transactions.id, transactions.type, transactions.created_at, transactions.updated_at,
		"tag".id, "tag".name, "tag".color,
		"from".id, "from".title, "from".type, "from".default, "from".user,
		"to".id, "to".title, "to".type, "to".default, "to".user
	FROM accounts.transactions
	LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
	LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to
	LEFT JOIN accounts.tags as "tag" ON "tag".id = transactions.tag
	WHERE "from".user = $1::UUID OR "to".user = $1::UUID;
	`).Scan(
		&txn.Id, &txn.Type, &txn.CreatedAt, &txn.UpdatedAt,
		&txn.From.Id, &txn.From.Title, &txn.From.Type, &txn.From.Default, &txn.From.User.Id,
		&txn.To.Id, &txn.To.Title, &txn.To.Type, &txn.To.Default, &txn.To.User.Id,
	)

	return &txn, err
}

// func (repo PsqlRepo) FindTransactionsByHotel() ([]entity.Transaction, error) {
// 	var txns []entity.Transaction = make([]entity.Transaction, 0)

// 	rows, err := repo.db.Query(`
// 	SELECT
// 		transactions.id, transactions.type, transactions.created_at, transactions.updated_at,
// 		"tag".id, "tag".name, "tag".color,
// 		"from".id, "from".title, "from".type, "from".default, "from".user,
// 		"to".id, "to".title, "to".type, "to".default, "to".user
// 	FROM accounts.transactions
// 	LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
// 	LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to
// 	LEFT JOIN accounts.tags as "tag" ON "tag".id = transactions.tag
// 	WHERE "from".user = $1::UUID OR "to".user = $1::UUID;
// 	`)

// 	if err != nil {
// 		repo.log.Println(err)
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var txn entity.Transaction
// 		err := rows.Scan(&txn.Id, &txn.Type, &txn.CreatedAt, &txn.UpdatedAt,
// 			&txn.From.Id, &txn.From.Title, &txn.From.Type, &txn.From.Default, &txn.From.User.Id,
// 			&txn.To.Id, &txn.To.Title, &txn.To.Type, &txn.To.Default, &txn.To.User.Id,
// 		)

// 		if err == nil {
// 			// Fetch txn details
// 			switch txn.Type {
// 			case entity.REPLENISHMENT:
// 				{
// 					txn.Details = nil
// 				}
// 			}

// 			txns = append(txns, txn)
// 		}
// 	}

// 	return txns, nil
// }

func (repo PsqlRepo) TransactionsDashboardRepo(year int) (interface{}, error) {

	rows, err := repo.db.Query(`
	SELECT
  EXTRACT(MONTH FROM created_at) AS month,
  COUNT(*) AS count
FROM accounts.transactions
WHERE EXTRACT(YEAR FROM created_at) = $1
GROUP BY EXTRACT(MONTH FROM created_at)
ORDER BY month;
	`, year)

	if err != nil {
		repo.log.Println(err)
		return nil, err
	}

	defer rows.Close()

	type data struct {
		Month int
		Count int
	}
	type Response struct {
		Jan int
		Feb int
		Mar int
		Apr int
		May int
		Jun int
		Jul int
		Aug int
		Sep int
		Oct int
		Nov int
		Dec int
	}

	var month_acount Response
	for rows.Next() {
		var res data
		err := rows.Scan(&res.Month, &res.Count)

		if err != nil {
			// Fetch txn details
			return nil, nil

		}
		switch res.Month {
		case 1:
			{
				month_acount.Jan = res.Count
				break
			}
		case 2:
			{
				month_acount.Feb = res.Count
				break
			}
		case 3:
			{
				month_acount.Mar = res.Count
				break
			}
		case 4:
			{
				month_acount.Apr = res.Count
				break
			}
		case 5:
			{
				month_acount.May = res.Count
				break
			}
		case 6:
			{
				month_acount.Jun = res.Count
				break
			}
		case 7:
			{
				month_acount.Jul = res.Count
				break
			}
		case 8:
			{
				month_acount.Aug = res.Count
				break
			}
		case 9:
			{
				month_acount.Sep = res.Count
				break
			}
		case 10:
			{
				month_acount.Oct = res.Count
				break
			}
		case 11:
			{
				month_acount.Nov = res.Count
				break
			}
		case 12:
			{
				month_acount.Dec = res.Count
				break
			}
		}

	}
	return month_acount, nil

}
