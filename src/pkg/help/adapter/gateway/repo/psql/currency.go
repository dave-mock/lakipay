package psql

import (
	"auth/src/pkg/help/core/entity"

	"github.com/google/uuid"
)

func (repo PsqlRepo) StoreCurrency(currency entity.Currency) error {

	_, err := repo.db.Exec(`
	INSERT INTO help.currencies (id, name, currency, symbol, rate, base_id, hidden ,created_at)
	VALUES ($1::UUID,$2,$3,$4,$5,$6::UUID,$7, $8);
	`, currency.Id, currency.Name, currency.Currency, currency.Symbol, currency.Rate, currency.BaseId, currency.Hidden, currency.CreatedAt)

	for i := 0; i < len(currency.Countries); i++ {
		_, err = repo.db.Exec(`
		INSERT INTO help.country_currencies (id, currency, country, "default")
		VALUES (gen_random_uuid(), $1::UUID, $2::UUID, $3);
		`, currency.Id, currency.Countries[i].Id, currency.Countries[i].Default)
		if err != nil {
			repo.log.Println(err)
		}
	}

	return err
}

func (repo PsqlRepo) FindCurrencies() ([]entity.Currency, error) {
	var currencies []entity.Currency = make([]entity.Currency, 0)

	rows, err := repo.db.Query(`
	SELECT id, name, currency, symbol, rate, base_id, hidden ,created_at, updated_at
	FROM help.currencies
	`)

	if err != nil {
		return currencies, err
	}

	defer rows.Close()

	for rows.Next() {
		var currency entity.Currency
		err := rows.Scan(&currency.Id, &currency.Name, &currency.Currency, &currency.Symbol, &currency.Rate, &currency.BaseId, &currency.Hidden, &currency.CreatedAt, &currency.UpdatedAt)
		if err != nil {
			// Do sth
		} else {
			// Find countries
			var countries []entity.CurrencyCountry = make([]entity.CurrencyCountry, 0)
			countryRows, err := repo.db.Query(`
			SELECT country, "default", countries.name, countries.default_name, countries.iso2, countries.flag
			FROM help.country_currencies
			INNER JOIN help.countries ON help.countries.id = country
			WHERE currency = $1;
			`, currency.Id)
			if err != nil {
				repo.log.Println(err)
			}
			for countryRows.Next() {
				var country entity.CurrencyCountry
				err = countryRows.Scan(&country.Id, &country.Default, &country.Name, &country.DefaultName, &country.Iso2, &country.Flag)
				if err != nil {
					repo.log.Println(err)
				}
				countries = append(countries, country)
			}
			currency.Countries = countries
			currencies = append(currencies, currency)
		}
	}

	return currencies, err
}

func (repo PsqlRepo) FindSupportedCurrencies() ([]entity.Currency, error) {
	var currencies []entity.Currency = make([]entity.Currency, 0)

	rows, err := repo.db.Query(`
	SELECT id, name, currency, symbol, rate, base_id, hidden ,created_at, updated_at
	FROM help.currencies
	WHERE hidden = 'false';
	`)

	if err != nil {
		return currencies, err
	}

	defer rows.Close()

	for rows.Next() {
		var currency entity.Currency
		err := rows.Scan(&currency.Id, &currency.Name, &currency.Currency, &currency.Symbol, &currency.Rate, &currency.BaseId, &currency.Hidden, &currency.CreatedAt, &currency.UpdatedAt)
		if err != nil {
			// Do sth
		} else {
			currencies = append(currencies, currency)
		}
	}

	return currencies, err
}

func (repo PsqlRepo) FindCurrencyById(id uuid.UUID) (*entity.Currency, error) {
	var currency *entity.Currency

	err := repo.db.QueryRow(`
	SELECT id, name, currency, symbol, rate, base_id, hidden ,created_at, updated_at
	FROM help.currencies
	WHERE id = $1::UUID;
	`, id).Scan(&currency.Id, &currency.Name, &currency.Currency, &currency.Symbol, &currency.Rate, &currency.BaseId, &currency.Hidden, &currency.CreatedAt, &currency.UpdatedAt)

	if err != nil {
		return currency, err
	}

	return currency, err
}
