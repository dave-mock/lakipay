package psql

import (
	"auth/src/pkg/org/usecase"
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

	var _schema = "org"
	// Map of table name with the corresponding sql
	var _tableScripts map[string]string = map[string]string{
		"categories": fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.categories
		(
			id uuid NOT NULL PRIMARY KEY,
			name character varying(255) NOT NULL,
			description text,
			icon text,
			parents uuid[],
			country_whitelist character varying(2)[] NOT NULL,
			country_blacklist character varying(2)[] NOT NULL,
			options uuid[],
			hidden boolean NOT NULL DEFAULT 'TRUE',
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"options": fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.options
		(
			id uuid NOT NULL PRIMARY KEY,
			name character varying(255) NOT NULL,
			description text,
			data_type character varying(255) NOT NULL,
			represented_in character varying(10),
			values text[],
			allow_custom_value boolean NOT NULL DEFAULT 'FALSE',
			validator jsonb,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"legal_conditions": fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.legal_conditions
		(
			id uuid NOT NULL PRIMARY KEY,
			name character varying(255) NOT NULL,
			description text,
			country_whitelist character varying(2)[] NOT NULL,
			country_blacklist character varying(2)[] NOT NULL,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"taxes": fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.taxes
		(
			id uuid NOT NULL PRIMARY KEY,
			name character varying(255) NOT NULL,
			description text,
			rate real NOT NULL,
			"from" character varying(255),
			country_whitelist character varying(2)[] NOT NULL,
			country_blacklist character varying(2)[] NOT NULL,
			hidden boolean NOT NULL DEFAULT 'TRUE',
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
	}

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	// Initialize schema
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

	return PsqlRepo{log, db, _schema}, nil
}
