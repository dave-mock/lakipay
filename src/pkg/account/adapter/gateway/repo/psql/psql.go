package psql

import (
	"auth/src/pkg/account/usecase"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type PsqlRepo struct {
	log *log.Logger
	db  *sql.DB
}

func New(log *log.Logger, db *sql.DB) (usecase.Repo, error) {

	var _schema = "accounts"
	// Map of table name with the corresponding sql
	var _tableScripts map[string]string = map[string]string{
		"accounts": fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.accounts
		(
			id uuid NOT NULL PRIMARY KEY,
			title character varying(50),
			"type" character varying(50) NOT NULL,
			"default" boolean NOT NULL DEFAULT 'FALSE',
			user_id uuid NOT NULL,
			verified boolean NOT NULL DEFAULT 'FALSE',
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"stored_accounts": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.stored_accounts
		(
			account_id uuid NOT NULL PRIMARY KEY,
			balance real NOT NULL
		);`, _schema),
		"bank_accounts": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.bank_accounts
		(
			account_id uuid NOT NULL PRIMARY KEY,
			account_number character varying(255) NOT NULL,
			holder_name character varying(255) NOT NULL,
			holder_phone character varying(255) NOT NULL,
			bank_id uuid NOT NULL
		);`, _schema),
		"banks": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.banks
		(
			id uuid NOT NULL PRIMARY KEY,
			name character varying(255) NOT NULL,
			short_name character varying(255) NOT NULL,
			bin character varying(8),
			swift_code character varying(8) NOT NULL,
			logo text NOT NULL,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"transactions": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.transactions
		(
			id uuid NOT NULL PRIMARY KEY,
			"from" uuid,
			"to" uuid,
			"type" character varying(255) NOT NULL,
			medium character varying(255) NOT NULL,
			comment character varying(255) ,
			tag character varying(255) ,
			verified boolean NOT NULL DEFAULT 'FALSE',
			reference character varying(12) NOT NULL,			ttl BIGINT,
			commission DOUBLE PRECISION,
			details JSONB,
			error_message TEXT,
			confirm_timestamp TIMESTAMPTZ,
			bank_reference VARCHAR(255),
			payment_method VARCHAR(255),
			test BOOLEAN,
			description TEXT,	
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

	return PsqlRepo{log: log, db: db}, nil
}
