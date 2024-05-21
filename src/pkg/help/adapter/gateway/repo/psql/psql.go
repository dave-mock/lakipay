package psql

import (
	"auth/src/pkg/help/usecase"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type PsqlRepo struct {
	log    *log.Logger
	db     *sql.DB
	schema string
}

func New(log *log.Logger, db *sql.DB) (usecase.Repo, error) {

	var _schema = "help"
	var _tableScripts map[string]string = map[string]string{
		"phone_prefixes": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.phone_prefixes
		(
			prefix character varying(4),
			pattern character varying(255),
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"countries": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.countries
		(
			id uuid NOT NULL PRIMARY KEY,
			name character varying(56),
			default_name character varying(255),
			iso2 character varying(2) NOT NULL,
			flag text NOT NULL,
			phone_prefix character varying(4) NOT NULL,
			hidden bool NOT NULL DEFAULT 'TRUE',
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"currencies": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.currencies
		(
			id uuid NOT NULL PRIMARY KEY,
			name character varying(56),
			currency character varying(255),
			symbol character varying(255),
			rate real,
			base_id uuid,
			hidden bool NOT NULL DEFAULT 'TRUE',
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"country_currencies": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.country_currencies
		(
			id uuid NOT NULL PRIMARY KEY,
			currency uuid,
			country uuid,
			"default" bool NOT NULL DEFAULT 'FALSE'
		);`, _schema),
	}

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s;`, _schema))
	if err != nil {
		return nil, err
	}

	// Initialize tables
	for _, v := range _tableScripts {
		_, err = tx.Exec(v)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return PsqlRepo{log: log, db: db, schema: _schema}, nil
}
