package psql

import (
	"auth/src/pkg/org/core/entity"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (repo PsqlRepo) FindCategories() ([]entity.Category, error) {
	var categories []entity.Category = make([]entity.Category, 0)

	rows, err := repo.db.Query(fmt.Sprintf(`
	SELECT 
		id, 
		name, 
		description,
		icon, 
		parents, 
		country_whitelist, 
		country_blacklist, 
		options, 
		hidden, 
		created_at, 
		updated_at
	FROM %s.categories;
	`, repo.schema))

	if err != nil {
		repo.log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var category entity.Category
		var options []uuid.UUID = make([]uuid.UUID, 0)
		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Description,
			&category.Icon,
			pq.Array(&category.Parents),
			pq.Array(&category.CountryWhitelist),
			pq.Array(&category.CountryBlacklist),
			pq.Array(&options),
			&category.Hidden,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			repo.log.Println(err)
		} else {
			if category.Parents == nil {
				category.Parents = make([]uuid.UUID, 0)
			}
			category.Options = make([]entity.Option, 0)
			for i := 0; i < len(options); i++ {
				var option entity.Option
				var validator []uint8
				err := repo.db.QueryRow(
					fmt.Sprintf(`
						SELECT 
							id, 
							name, 
							description, 
							data_type, 
							represented_in, 
							values, 
							allow_custom_value, 
							validator, 
							created_at, 
							updated_at
						FROM %s.options`,
						repo.schema,
					),
				).Scan(
					&option.Id,
					&option.Name,
					&option.Description,
					&option.DataType,
					&option.RepresentedIn,
					pq.Array(&option.Values),
					&option.AllowCustomValue,
					&validator,
					&option.CreatedAt,
					&option.UpdatedAt,
				)

				if err != nil {
					repo.log.Println(err)
				} else {
					json.Unmarshal(validator, &option.Validator)
					category.Options = append(category.Options, option)
				}
			}
			categories = append(categories, category)
		}

	}

	return categories, nil
}

func (repo PsqlRepo) StoreCategory(category entity.Category) error {

	// Store category

	// Begin transaction
	tx, err := repo.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	// Store options
	for i := 0; i < len(category.Options); i++ {
		validator, _ := json.Marshal(category.Options[i].Validator)
		_, err = tx.Exec(fmt.Sprintf(`
		INSERT INTO %s.options (
			id, 
			name, 
			description, 
			data_type, 
			represented_in, 
			values, 
			allow_custom_value, 
			validator, 
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		`, repo.schema),
			category.Options[i].Id,
			category.Options[i].Name,
			category.Options[i].Description,
			category.Options[i].DataType,
			category.Options[i].RepresentedIn,
			pq.Array(category.Options[i].Values),
			category.Options[i].AllowCustomValue,
			validator,
			category.Options[i].CreatedAt,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Store category

	var options []uuid.UUID = make([]uuid.UUID, 0)
	for i := 0; i < len(category.Options); i++ {
		options = append(options, category.Options[i].Id)
	}

	_, err = tx.Exec(fmt.Sprintf(`
	INSERT INTO %s.categories (
		id, 
		name,
		description,
		icon,
		parents, 
		country_whitelist, 
		country_blacklist, 
		options,
		hidden,
		created_at
	)
	VALUES ($1::UUID, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, repo.schema),
		category.Id,
		category.Name,
		category.Description,
		category.Icon,
		pq.Array(category.Parents),
		pq.StringArray(category.CountryWhitelist),
		pq.StringArray(category.CountryBlacklist),
		pq.Array(options),
		category.Hidden,
		category.CreatedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (repo PsqlRepo) FindCategoryByName(name string) (*entity.Category, error) {
	var category entity.Category
	var options []uuid.UUID = make([]uuid.UUID, 0)

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT 
		id, 
		name, 
		description,
		icon, 
		parents, 
		country_whitelist, 
		country_blacklist, 
		options, 
		hidden, 
		created_at, 
		updated_at
	FROM %s.categories
	WHERE name = $1;
	`, repo.schema), name).Scan(
		&category.Id,
		&category.Name,
		&category.Description,
		&category.Icon,
		pq.Array(&category.Parents),
		pq.Array(&category.CountryWhitelist),
		pq.Array(&category.CountryBlacklist),
		pq.Array(&options),
		&category.Hidden,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	repo.log.Println("category")
	repo.log.Println(category)
	repo.log.Println(err)

	if err != nil {
		return nil, err
	}

	return &category, nil
}
