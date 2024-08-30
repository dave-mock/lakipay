package psql

import (
	"auth/src/pkg/account/core/entity"
	"auth/src/pkg/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (repo PsqlRepo) UpdateUserRepo(users entity.User2) (entity.User2, error) {
	var res_users entity.User2
	_, err := repo.db.Exec(`
	update  auth.users set first_name = $1
	where id  = $2
	`, users.Name, users.Id)

	_ = repo.db.QueryRow(`select id ,first_name from auth.users where id =$1`, users.Id).Scan(&res_users.Id, &res_users.Name)

	return res_users, err
}

func (repo PsqlRepo) StoreTransactionSession(preSession entity.TransactionSession) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.transaction_sessions (id, token, created_at)
	VALUES ($1::UUID, $2, $3);
	`, "accounts"), preSession.Id, sql.NullString{Valid: preSession.Token != "", String: preSession.Token}, preSession.CreatedAt)

	return err
}
func (repo PsqlRepo) StoreKeys(id uuid.UUID, publicKey string, private_key string, password string, username string) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.merchant_keys (id, public_key, private_key,merchant_id,password,username)
	VALUES ($1::UUID, $2, $3,$4,$5,$6);
	`, "accounts"), uuid.New(), publicKey, private_key, id, password, username)

	return err
}

func (repo PsqlRepo) GetMerchantsKeys(username string) (entity.MerchantKeys, error) {
	var todo entity.MerchantKeys
	sqlStmt := `SELECT * FROM accounts.merchant_keys WHERE  username=$1`

	err := repo.db.QueryRow(sqlStmt, username).Scan(&todo.Id, &todo.PublickKey, &todo.MerchantId, &todo.PrivateKey, &todo.Username, &todo.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.MerchantKeys{}, err
	} else if err != nil {
		return entity.MerchantKeys{}, err
	}

	return todo, nil
}
func (repo PsqlRepo) GetApiKeysRepo(id uuid.UUID) (entity.MerchantKeys, error) {
	var todo entity.MerchantKeys
	sqlStmt := `SELECT * FROM accounts.merchant_keys WHERE  merchant_id=$1`

	err := repo.db.QueryRow(sqlStmt, id).Scan(&todo.Id, &todo.PublickKey, &todo.MerchantId, &todo.PrivateKey, &todo.Username, &todo.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.MerchantKeys{}, err
	} else if err != nil {
		return entity.MerchantKeys{}, err
	}

	return todo, nil
	// 	var exists string

	//     err := repo.db.QueryRow("SELECT private_key FROM accounts.merchant_keys WHERE merchant_id=$1", id).Scan(&exists)

	//     if err != nil {
	//        return  "",err
	// 	}

	// return  exists,err
}
func (repo PsqlRepo) CheckMerchantsKeysByUsername(username string) (error, bool) {
	var exists bool

	err := repo.db.QueryRow("SELECT EXISTS (SELECT 1 FROM accounts.merchant_keys WHERE username=$1)", username).Scan(&exists)

	if err != nil {
		return err, false
	}

	return err, exists
}

func (repo PsqlRepo) StoreTransaction(txn entity.Transaction) error {
	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	// data, err := utils.AesEncrption("marsal")
	// data, err := utils.AesDecription("2BvHg7Fq4TSwdC8diew2WA==")

	// var prev_txn entity.Transaction
	// err = repo.db.QueryRow(`select id,created_at from accounts.transactions order by created_at desc limit 1`).Scan(&prev_txn.Id, &prev_txn.CreatedAt)
	// if err != nil {
	// 	return err
	// }
	// fmt.Print(prev_txn)
	token_str := txn.Id.String() + txn.CreatedAt.String()
	token, err := utils.AesEncrption(token_str)
	if err != nil {
		return err
	}

	amount, err := utils.AesEncrption(strconv.FormatFloat(txn.Amount, 'f', -1, 64))
	if err != nil {
		return err
	}

	fmt.Printf("Type of a: |||||||||| %s\n", reflect.TypeOf(txn.To.Id))
	// if txn.To.Id == sql.Null {

	// }

	switch txn.Type {
	case entity.REPLENISHMENT:
		{
			txn.Commission = 1
		}
	case entity.P2P:
		{
			txn.Commission = 1.25

		}
	case entity.SALE:
		{
			txn.Commission = 2.75

		}

	case entity.SETTLEMENT:
		{
			txn.Commission = 3

		}

	}

	com := (txn.Amount * txn.Commission) / 100
	total_amount, err := utils.AesEncrption(strconv.FormatFloat(txn.Amount+com, 'f', -1, 64))
	if err != nil {
		return err
	}

	fmt.Println("|||||||||||||||||||||||||||||||| total_amount: ", com)

	_, err = tx.Exec(`
	INSERT INTO accounts.transactions (id, "from", "to", "type", "reference", verified, created_at,token,medium,amount,has_challenge,commission,total_amount,phone)
	VALUES ($1::UUID, $2::UUID, $3::UUID, $4, $5, $6, $7,$8,$9,$10,$11,$12,$13,$14)
	`, txn.Id, txn.From.Id, txn.To.Id, txn.Type, txn.Reference, false, txn.CreatedAt, token, txn.Medium, amount, txn.HasChallenge, txn.Commission, total_amount, txn.Phone)

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
	case entity.SALE:
		{
			// txnDetail := txn.Details.(entity.P2p)
			// amount:= utils.AesEncrption(txnDetail.Amount)
			_, err = tx.Exec(`
			INSERT INTO accounts.sales (transaction_id)
			VALUES ($1::UUID)
			`, txn.Id)
			if err != nil {
				tx.Rollback()
				return err
			}

		}
	case entity.SETTLEMENT:
		{
			txnDetail := txn.Details.(entity.P2p)
			details, _ := json.Marshal(txnDetail)
			// amount:= utils.AesEncrption(txnDetail.Amount)
			_, err = tx.Exec(`
			INSERT INTO accounts.settlements (transaction_id,details)
			VALUES ($1::UUID,$2)
			`, txn.Id, details)
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

func (repo PsqlRepo) UpdateGeneratedChallenge(challenge string, id uuid.UUID, deviceId string) error {
	sqlStmt := `
		UPDATE accounts.public_keys SET challenge=$1, expires_at=$2, used=FALSE WHERE user_id=$3 and device_id=$4 `
	expiresAt := time.Now().Add(5 * time.Minute)

	_, err := repo.db.Exec(sqlStmt, challenge, expiresAt, id, deviceId)
	if err != nil {
		return err
	}

	return nil
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
		"to".id, "to".title, "to".type, "to".default, "to".user_id,token,transactions.amount
	FROM accounts.transactions
	LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
	LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to
	order by created_at ASC
	;
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
		var token sql.NullString

		var toTitle sql.NullString
		var toType sql.NullString
		var toDefualt sql.NullBool
		var amount sql.NullString

		err := rows.Scan(&txn.Id, &txn.Type, &txn.CreatedAt, &txn.UpdatedAt,
			&medium, &comment, &txn.Verified, &txn.Reference, &ttl, &commistion, &txtDetiail, &errMsg, &time, &BankReference,
			&pymentMethod, &test, &description,
			&txn.From.Id, &txn.From.Title, &txn.From.Type, &txn.From.Default, &txn.From.User.Id,
			&txn.To.Id, &toTitle, &toType, &toDefualt, &txn.To.User.Id, &token, &amount,
		)

		if toDefualt.Valid {
			txn.To.Default = toDefualt.Bool
		}
		if toTitle.Valid {
			txn.To.Title = toTitle.String
		}
		if toType.Valid {
			txn.To.Type = entity.AccountType(toType.String)
		}

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
		if token.Valid {
			txn.Token = token.String
		}

		if amount.Valid {
			amount2, err := utils.AesDecription(amount.String)
			if err != nil {
				return nil, err

			}
			amount_float, err := strconv.ParseFloat(amount2, 64)
			if err != nil {
				return nil, err
			}
			txn.Amount = amount_float

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

			targettoken := txns[i].Id.String() + txns[i].CreatedAt.String()

			token, err := utils.AesDecription(element.Token)
			if err != nil {
				return nil, err
			}

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
		"to".id, "to".title, "to".type, "to".default, "to".user_id,transactions.amount,transactions.total_amount,transactions.currency
	FROM accounts.transactions
	LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
	LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to
	WHERE ("from".id = $1::UUID OR "to".id = $1::UUID ) and accounts.transactions.verified = true;
	`, id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return []entity.Transaction{}, nil
		}
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
		var amount sql.NullString
		var total_amount sql.NullString

		err := rows.Scan(&txn.Id, &txn.Type, &txn.CreatedAt, &txn.UpdatedAt,
			&medium, &comment, &txn.Verified, &txn.Reference, &ttl, &commistion, &txtDetiail, &errMsg, &time, &BankReference,
			&pymentMethod, &test, &description,
			&txn.From.Id, &txn.From.Title, &txn.From.Type, &txn.From.Default, &txn.From.User.Id,
			&txn.To.Id, &txn.To.Title, &txn.To.Type, &txn.To.Default, &txn.To.User.Id, &amount, &total_amount, &txn.Currency,
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
		if amount.Valid {
			amount2, err := utils.AesDecription(amount.String)
			if err != nil {
				return nil, err

			}
			amount_float, err := strconv.ParseFloat(amount2, 64)
			if err != nil {
				return nil, err
			}
			txn.Amount = amount_float

		}
		if total_amount.Valid {
			total_amount2, err := utils.AesDecription(total_amount.String)
			fmt.Print("++++++++++++++++++++++++++++++++++++++++++++++++++++", total_amount2)
			if err != nil {
				return nil, err

			}
			total_amount_float, err := strconv.ParseFloat(total_amount2, 64)
			if err != nil {
				return nil, err
			}
			txn.TotalAmount = total_amount_float

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
	// "tag".id, "tag".name, "tag".color,
	// LEFT JOIN accounts.tags as "tag" ON "tag".id = transactions.tag
	var amount sql.NullString
	var titleTo sql.NullString
	var defultTo sql.NullBool
	var typeTo sql.NullString
	var phone sql.NullString

	// var titleFrom sql.NullString

	err := repo.db.QueryRow(`
	SELECT 
		transactions.id, transactions.type,transactions.amount,transactions.has_challenge, transactions.created_at, transactions.updated_at,
		
		"from".id, "from".title, "from".type, "from".default, "from".user_id,
		"to".id, "to".title, "to".type, "to".default, "to".user_id,transactions.phone
	FROM accounts.transactions
	LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
	LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to
	
	WHERE transactions.id=$1;
	`, id).Scan(
		&txn.Id, &txn.Type, &amount, &txn.HasChallenge, &txn.CreatedAt, &txn.UpdatedAt,
		&txn.From.Id, &txn.From.Title, &txn.From.Type, &txn.From.Default, &txn.From.User.Id,
		&txn.To.Id, &titleTo, &typeTo, &defultTo, &txn.To.User.Id, &phone,
	)

	if phone.Valid {
		txn.Phone = phone.String
	}
	if defultTo.Valid {
		txn.To.Default = defultTo.Bool
	}
	if titleTo.Valid {
		txn.To.Title = titleTo.String
	}
	if typeTo.Valid {
		txn.To.Type = entity.AccountType(typeTo.String)
	}
	if amount.Valid {
		amount2, err := utils.AesDecription(amount.String)
		if err != nil {
			return nil, err

		}
		amount_float, err := strconv.ParseFloat(amount2, 64)
		if err != nil {
			return nil, err
		}
		txn.Amount = amount_float

	}

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

func (repo PsqlRepo) GetPuplicKey(challenge string, id uuid.UUID) ([]*entity.PublicKey, error) {
	var todos []*entity.PublicKey
	sqlStmt := `
	SELECT * FROM accounts.public_keys WHERE user_id=$1 AND challenge=$2   AND used=FALSE
	`

	rows, err := repo.db.Query(sqlStmt, id, challenge)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo entity.PublicKey
		err := rows.Scan(&todo.ID, &todo.UserID, &todo.PublicKey, &todo.DeviceID, &todo.Challenge, &todo.ExpiresAt, &todo.Used, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (repo PsqlRepo) UpdatePublicKeysUsed(id uuid.UUID) error {
	sqlStmt := `
		UPDATE accounts.public_keys SET  used=FALSE WHERE user_id=$1`

	_, err := repo.db.Exec(sqlStmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo PsqlRepo) UpdateTransaction(id uuid.UUID) error {
	print("||||||||| ================================================= UpdateTransaction")
	sqlStmt := `
		UPDATE accounts.transactions SET  verified=true WHERE id=$1`

	_, err := repo.db.Exec(sqlStmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo PsqlRepo) GetstorePublicKeyHandler(key string, id uuid.UUID, device string) error {
	var count int
	err := repo.db.QueryRow(`SELECT COUNT(*) FROM accounts.public_keys WHERE user_id = $1 AND device_id = $2`, id, device).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing public keys: %v", err)
	}
	if count != 0 {
		sqlStmt := `
		UPDATE accounts.public_keys SET  used=FALSE,public_key=$1 WHERE user_id=$2`

		_, err := repo.db.Exec(sqlStmt, key, id)
		if err != nil {
			return err
		}
	} else {
		sqlStmt := `
		INSERT INTO accounts.public_keys (id,user_id, public_key, device_id) VALUES ($1, $2, $3,$4)`

		_, err = repo.db.Exec(sqlStmt, uuid.New(), id, key, device)
		if err != nil {
			return err
		}
	}

	return nil
}
