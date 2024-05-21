package psql

import (
	"auth/src/pkg/account/core/entity"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func (repo PsqlRepo) FindAccountsByUserId(userId uuid.UUID) ([]entity.Account, error) {
	var accs []entity.Account = make([]entity.Account, 0)

	rows, err := repo.db.Query(`
	SELECT id, title, "type", "default"
	FROM accounts.accounts
	WHERE "user_id" = $1::UUID AND verified = $2;
	`, userId, true)

	repo.log.Println(rows)

	if err != nil {
		repo.log.Println("[ACC PSQL] FOUND ACCS - 1")
		repo.log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		repo.log.Println("[ACC PSQL] FOUND ACCS")
		var nullTitle sql.NullString
		var acc entity.Account
		err := rows.Scan(&acc.Id, &nullTitle, &acc.Type, &acc.Default)
		if nullTitle.Valid {
			acc.Title = nullTitle.String
		}
		if err == nil {
			repo.log.Println("[ACC PSQL] NO ERR SCANNING")
			// Switch account types
			switch acc.Type {
			case entity.STORED:
				{
					repo.log.Println("[ACC PSQL] STORED ACC")
					// Get stored account
					var storedAcc entity.StoredAccount
					err = repo.db.QueryRow(`
					SELECT balance
					FROM accounts.stored_accounts
					WHERE account_id = $1::UUID;
					`, acc.Id).Scan(&storedAcc.Balance)

					acc.Detail = storedAcc

					if err == nil {
						repo.log.Println("[ACC PSQL] STORED ACC ERROR")
						accs = append(accs, acc)
					}
				}
			case entity.BANK:
				{
					repo.log.Println("[ACC PSQL] BANK ACC")
					// Get stored account
					var bankAcc entity.BankAccount
					err = repo.db.QueryRow(`
				SELECT account_number, holder_name, holder_phone,
				banks.id, banks.name, banks.short_name, banks.swift_code, banks.logo
				FROM accounts.bank_accounts
				INNER JOIN accounts.banks ON accounts.banks.id = bank_id
				WHERE account_id = $1::UUID;
				`, acc.Id).Scan(
						&bankAcc.Number,
						&bankAcc.Holder.Name,
						&bankAcc.Holder.Phone,
						&bankAcc.Bank.Id,
						&bankAcc.Bank.Name,
						&bankAcc.Bank.ShortName,
						&bankAcc.Bank.SwiftCode,
						&bankAcc.Bank.Logo,
					)

					acc.Detail = bankAcc

					repo.log.Println(acc.Detail)

					repo.log.Println(err)

					if err == nil {
						repo.log.Println("[ACC PSQL] BANK ACC ERROR")
						accs = append(accs, acc)
					}
				}

			}
		}
	}

	repo.log.Println(len(accs))

	return accs, nil
}

func (repo PsqlRepo) StoreAccount(acc entity.Account) error {

	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	INSERT INTO accounts.accounts (id, title, "type", "default", "user_id", "verified", "created_at")
	VALUES ($1::UUID, $2, $3, $4, $5::UUID, $6, $7);
	`, acc.Id, sql.NullString{Valid: acc.Title != "", String: acc.Title}, acc.Type, acc.Default, acc.User.Id, true, acc.CreatedAt)

	if err != nil {
		tx.Rollback()
		return err
	}

	switch acc.Type {
	case entity.STORED:
		{
			_, err = tx.Exec(`
			INSERT INTO accounts.stored_accounts (account_id, balance)
			VALUES ($1::UUID, $2)
			`, acc.Id, acc.Detail.(entity.StoredAccount).Balance)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	case entity.BANK:
		{
			_, err = tx.Exec(`
				INSERT INTO accounts.bank_accounts (account_id, bank_id, account_number, holder_name, holder_phone)
				VALUES ($1::UUID, $2::UUID, $3, $4, $5)
				`, acc.Id, acc.Detail.(entity.BankAccount).Bank.Id, acc.Detail.(entity.BankAccount).Number, acc.Detail.(entity.BankAccount).Holder.Name, acc.Detail.(entity.BankAccount).Holder.Phone)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}

	return err
}

func (repo PsqlRepo) UpdateAccount(acc entity.Account) error {

	repo.log.Println("Update Account")

	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	repo.log.Println(acc.VerificationStatus.Verified)

	_, err = tx.Exec(`
	UPDATE accounts.accounts 
	SET title = $2, type = $3, "default" = $4, "user_id" = $5, verified = $6
	WHERE id = $1;
	`, acc.Id, sql.NullString{Valid: acc.Title != "", String: acc.Title}, acc.Type, acc.Default, acc.User.Id, acc.VerificationStatus.Verified)

	if err != nil {
		repo.log.Println(err)
		tx.Rollback()
		return err
	}

	switch acc.Type {
	case entity.STORED:
		{
			_, err = tx.Exec(`
			UPDATE accounts.stored_accounts 
			SET balance = $2
			WHERE account = $1;
			`, acc.Id, acc.Detail.(entity.StoredAccount).Balance)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	case entity.BANK:
		{
			// _, err = tx.Exec(`
			// 	UPDATE accounts.bank_accounts (account_id, bank_id, account_number, holder_name, holder_phone)
			// 	VALUES ($1::UUID, $2::UUID, $3, $4, $5)
			// 	`, acc.Id, acc.Detail.(entity.BankAccount).Bank.Id, acc.Detail.(entity.BankAccount).Number, acc.Detail.(entity.BankAccount).Holder.Name, acc.Detail.(entity.BankAccount).Holder.Phone)
			// if err != nil {
			// 	tx.Rollback()
			// 	return err
			// }
		}
	}

	err = tx.Commit()
	if err != nil {
		repo.log.Println(err)
		tx.Rollback()
	}

	return err
}

func (repo PsqlRepo) FindAccountById(accId uuid.UUID) (*entity.Account, error) {
	var acc entity.Account

	var title sql.NullString

	err := repo.db.QueryRow(`
	SELECT id, title, "type", "default", "user_id", "verified"
	FROM accounts.accounts
	WHERE "id" = $1::UUID;
	`, accId).Scan(
		&acc.Id, &title, &acc.Type, &acc.Default, &acc.User.Id, &acc.VerificationStatus.Verified,
	)

	if err != nil {
		return nil, err
	}

	if title.Valid {
		acc.Title = title.String
	}

	// Switch details
	switch acc.Type {
	case entity.STORED:
		{
			repo.log.Println("[ACC PSQL] STORED ACC")
			// Get stored account
			var storedAcc entity.StoredAccount
			err = repo.db.QueryRow(`
			SELECT balance
			FROM accounts.stored_accounts
			WHERE account_id = $1::UUID;
			`, acc.Id).Scan(&storedAcc.Balance)

			if err != nil {
				return nil, err
			}

			acc.Detail = storedAcc
		}
	case entity.BANK:
		{
			repo.log.Println("[ACC PSQL] STORED ACC")
			// Get stored account
			var bankAcc entity.BankAccount
			err = repo.db.QueryRow(`
				SELECT account_number, holder_name, holder_phone,
				banks.id, banks.name, banks.short_name, banks.swift_code, banks.logo
				FROM accounts.bank_accounts
				INNER JOIN accounts.banks ON accounts.banks.id = bank_id
				WHERE account_id = $1::UUID;
				`, acc.Id).Scan(
				&bankAcc.Number,
				&bankAcc.Holder.Name,
				&bankAcc.Holder.Phone,
				&bankAcc.Bank.Id,
				&bankAcc.Bank.Name,
				&bankAcc.Bank.ShortName,
				// &bankAcc.Bank.BIN,
				&bankAcc.Bank.SwiftCode,
				&bankAcc.Bank.Logo,
			)

			if err != nil {
				return nil, err
			}

			acc.Detail = bankAcc
		}

	}

	return &acc, nil
}

func (repo PsqlRepo) DeleteAccount(accId uuid.UUID) error {

	_, err := repo.db.Exec(`
	DELETE FROM accounts.accounts
	WHERE id = $1::UUID
	`, accId)

	return err
}
