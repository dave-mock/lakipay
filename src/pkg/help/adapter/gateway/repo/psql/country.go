package psql

import (
	"auth/src/pkg/help/core/entity"
	"fmt"

	"github.com/google/uuid"
)

func (repo PsqlRepo) StorePhonePrefix(phonePrefix entity.PhonePrefix) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.phone_prefixes (prefix, pattern, created_at)
	VALUES ($1, $2, $3);
	`, repo.schema), phonePrefix.Prefix, phonePrefix.Pattern, phonePrefix.CreatedAt)

	return err
}

func (repo PsqlRepo) StoreCountry(country entity.Country) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.countries (id, name, default_name, iso2, flag, phone_prefix ,hidden, created_at)
	VALUES ($1::UUID, $2, $3, $4, $5, $6, $7, $8);
	`, repo.schema), country.Id, country.Name, country.DefaultName, country.Iso2, country.Flag, country.PhonePrefix.Prefix, country.Hidden, country.CreatedAt)

	return err
}

func (repo PsqlRepo) FindPhoneprefixByPrefix(prefix string) (*entity.PhonePrefix, error) {
	var phonePrefix entity.PhonePrefix

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT prefix, pattern, created_at
	FROM %s.phone_prefixes
	WHERE prefix = $1;
	`, repo.schema), prefix).Scan(&phonePrefix.Prefix, &phonePrefix.Pattern, &phonePrefix.CreatedAt)

	return &phonePrefix, err
}

func (repo PsqlRepo) FindCountries() ([]entity.Country, error) {
	var countries []entity.Country = make([]entity.Country, 0)

	rows, err := repo.db.Query(fmt.Sprintf(`
	SELECT id, name, default_name, iso2, flag, hidden, prefix, pattern
	FROM %s.countries
	INNER JOIN %s.phone_prefixes ON %s.phone_prefixes.prefix = phone_prefix;
	`, repo.schema, repo.schema, repo.schema))

	if err != nil {
		return countries, err
	}

	defer rows.Close()

	for rows.Next() {
		var country entity.Country
		err := rows.Scan(&country.Id, &country.Name, &country.DefaultName, &country.Iso2, &country.Flag, &country.Hidden, &country.PhonePrefix.Prefix, &country.PhonePrefix.Pattern)
		if err != nil {
			// Do sth
		} else {
			countries = append(countries, country)
		}
	}

	return countries, err
}

func (repo PsqlRepo) FindSupportedCountries() ([]entity.Country, error) {
	var countries []entity.Country = make([]entity.Country, 0)

	rows, err := repo.db.Query(fmt.Sprintf(`
	SELECT id, name, default_name, iso2, flag, hidden, prefix, pattern
	FROM %s.countries
	INNER JOIN %s.phone_prefixes ON %s.phone_prefixes.prefix = phone_prefix
	WHERE hidden = 'false';
	`, repo.schema, repo.schema, repo.schema))

	if err != nil {
		return countries, err
	}

	defer rows.Close()

	for rows.Next() {
		var country entity.Country
		err := rows.Scan(&country.Id, &country.Name, &country.DefaultName, &country.Iso2, &country.Flag, &country.Hidden, &country.PhonePrefix.Prefix, &country.PhonePrefix.Pattern)
		if err != nil {
			// Do sth
		} else {
			countries = append(countries, country)
		}
	}

	return countries, err
}

func (repo PsqlRepo) FindCountryById(id uuid.UUID) (*entity.Country, error) {
	var country *entity.Country

	err := repo.db.QueryRow(`
	SELECT id, name, default_name, iso2, flag, hidden, prefix, pattern
	FROM public.countries
	INNER JOIN public.phone_prefixes ON public.phone_prefixes.prefix = phone_prefix
	WHERE id = $1::UUID;
	`, id).Scan(&country.Id, &country.Name, &country.DefaultName, &country.Iso2, &country.Flag, &country.Hidden, &country.PhonePrefix.Prefix, &country.PhonePrefix.Pattern)

	if err != nil {
		return country, err
	}

	return country, err
}
