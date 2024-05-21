package psql

import (
	"auth/src/pkg/org/core/entity"
	"fmt"

	"github.com/lib/pq"
)

func (repo PsqlRepo) StoreLegalCondition(i entity.LegalCondition) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.legal_conditions (id, name, country_whitelist, country_blacklist, created_at)
	VALUES ($1::UUID,$2,$3,$4,$5)
	`, repo.schema), i.Id, i.Name, pq.StringArray(i.CountryWhitelist), pq.StringArray(i.CountryBlacklist), i.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (repo PsqlRepo) FindLegalConditions() ([]entity.LegalCondition, error) {
	var legalConditions []entity.LegalCondition = make([]entity.LegalCondition, 0)

	rows, err := repo.db.Query(fmt.Sprintf(`
	SELECT id, name, country_whitelist, country_blacklist, created_at, updated_at
	FROM %s.legal_conditions
	`, repo.schema))

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var legalCondition entity.LegalCondition
		err := rows.Scan(&legalCondition.Id, &legalCondition.Name, pq.Array(&legalCondition.CountryWhitelist), pq.Array(&legalCondition.CountryBlacklist), &legalCondition.CreatedAt, &legalCondition.UpdatedAt)
		if err != nil {
			repo.log.Println(err)
		} else {
			legalConditions = append(legalConditions, legalCondition)
		}
	}

	return legalConditions, nil
}

func (repo PsqlRepo) FindLegalConditionByName(name string) (*entity.LegalCondition, error) {
	var legalCondition entity.LegalCondition

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT id, name, country_whitelist, country_blacklist, created_at, updated_at
	FROM %s.legal_conditions
	WHERE name = $1
	`, repo.schema), name).Scan(&legalCondition.Id, &legalCondition.Name, pq.Array(&legalCondition.CountryWhitelist), pq.Array(&legalCondition.CountryBlacklist), &legalCondition.CreatedAt, &legalCondition.UpdatedAt)

	repo.log.Println("legal condition find err")
	repo.log.Println(err)
	repo.log.Println(legalCondition)

	if err != nil {
		return nil, err
	}

	return &legalCondition, nil
}
