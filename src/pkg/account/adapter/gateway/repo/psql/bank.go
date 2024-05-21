package psql

import (
	"auth/src/pkg/account/core/entity"
	"database/sql"

	"github.com/google/uuid"
)

func (repo PsqlRepo) StoreBank(bank entity.Bank) error {
	_, err := repo.db.Exec(`
	INSERT INTO accounts.banks(id, name, short_name, bin, swift_code, logo, created_at)
	VALUES ($1::UUID, $2, $3, $4, $5, $6, $7)
	`, bank.Id, bank.Name, sql.NullString{Valid: bank.ShortName != "", String: bank.ShortName}, bank.BIN, bank.SwiftCode, bank.Logo, bank.CreatedAt)

	return err
}

func (repo PsqlRepo) FindBanks() ([]entity.Bank, error) {
	var banks []entity.Bank

	rows, err := repo.db.Query(`
	SELECT id, name, short_name, bin, swift_code, logo, created_at, updated_at
	FROM accounts.banks;
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var bank entity.Bank
		var shortName sql.NullString
		err := rows.Scan(&bank.Id, &bank.Name, &shortName, &bank.BIN, &bank.SwiftCode, &bank.Logo, &bank.CreatedAt, &bank.UpdatedAt)
		if err == nil {
			if shortName.Valid {
				bank.ShortName = shortName.String
			}
			banks = append(banks, bank)
		}
	}

	repo.log.Println(len(banks))
	repo.log.Println(banks)

	return banks, nil
}

func (repo PsqlRepo) FindBankById(id uuid.UUID) (*entity.Bank, error) {
	var bank entity.Bank

	var shortName sql.NullString

	err := repo.db.QueryRow(`
	SELECT id, name, short_name, bin, swift_code, logo, created_at, updated_at
	FROM accounts.banks
	WHERE id = $1::UUID
	`, id).Scan(
		&bank.Id,
		&bank.Name,
		&shortName,
		&bank.BIN,
		&bank.SwiftCode,
		&bank.Logo,
		&bank.CreatedAt,
		&bank.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if shortName.Valid {
		bank.ShortName = shortName.String
	}

	return &bank, err
}
