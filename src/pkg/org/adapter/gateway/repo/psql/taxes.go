package psql

import (
	"auth/src/pkg/org/core/entity"
	"fmt"

	"github.com/lib/pq"
)

// Store tax information
func (repo PsqlRepo) StoreTax(v entity.Tax) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.taxes (id, name, description, rate, "from", country_whitelist, country_blacklist, hidden, created_at)
	VALUES ($1::UUID, $2, $3, $4, $5, $6, $7, $8, $9);
	`, repo.schema), v.Id, v.Name, v.Description, v.Rate, v.From, pq.StringArray(v.CountryWhitelist), pq.StringArray(v.CountryBlacklist), v.Hidden, v.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Find taxes

func (repo PsqlRepo) FindTaxes() ([]entity.Tax, error) {
	var taxes []entity.Tax

	rows, err := repo.db.Query(fmt.Sprintf(`
	SELECT id, name, description, rate, "from", country_whitelist, country_blacklist, hidden, created_at, updated_at
	FROM %s.taxes;
	`, repo.schema))

	if err != nil {
		return taxes, nil
	}

	for rows.Next() {
		var tax entity.Tax
		err := rows.Scan(
			&tax.Id,
			&tax.Name,
			&tax.Description,
			&tax.Rate,
			&tax.From,
			pq.Array(&tax.CountryWhitelist),
			pq.Array(&tax.CountryBlacklist),
			&tax.Hidden,
			&tax.CreatedAt,
			&tax.UpdatedAt,
		)
		if err != nil {
			repo.log.Println(err)
		} else {
			taxes = append(taxes, tax)
		}
	}

	return taxes, nil
}
