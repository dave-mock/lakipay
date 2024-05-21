package repo

import (
	"auth/src/pkg/auth/usecase"
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

func NewPsqlRepo(log *log.Logger, db *sql.DB) (usecase.Repo, error) {

	log.SetPrefix("[AUTH] [GATEWAY] [REPO] [PSQL] ")

	var _schema = "auth"
	// Map of table name with the corresponding sql
	var _tableScripts map[string]string = map[string]string{
		"pre_sessions": fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.pre_sessions
		(
			id uuid NOT NULL PRIMARY KEY,
			token character varying(255) NOT NULL,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"devices": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.devices
		(
			id uuid NOT NULL PRIMARY KEY,
    		ip character varying(15) NOT NULL,
    		name character varying(255) NOT NULL,
    		agent character varying(255) NOT NULL,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"device_auths": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.device_auths
		(
			id uuid NOT NULL PRIMARY KEY,
    		token character varying(255) NOT NULL,
    		device uuid NOT NULL,
    		status bool NOT NULL,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"phones": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.phones
		(
			id uuid NOT NULL PRIMARY KEY,
    		prefix character varying(3) NOT NULL,
    		number character varying(15) NOT NULL,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"phone_auths": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.phone_auths
		(
			id uuid NOT NULL PRIMARY KEY,
			token character varying(255) NOT NULL,
    		phone_id uuid NOT NULL,
			code character varying(255) NOT NULL,
			status bool NOT NULL,
    		method character varying(255) NOT NULL,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"phone_identities": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.phone_identities
		(
			id uuid NOT NULL PRIMARY KEY,
			user_id uuid NOT NULL,
    		phone_id uuid NOT NULL,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"users": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.users
		(
			id uuid NOT NULL PRIMARY KEY,
			sir_name character varying(15),
			first_name character varying(15) NOT NULL,
			last_name character varying(15),
			gender character varying(10),
			date_of_birth timestamp with time zone,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"sessions": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.sessions
		(
			id uuid NOT NULL PRIMARY KEY,
			token character varying(255),
			user_id uuid NOT NULL,
			device_id uuid NOT NULL,
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"password_identities": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.password_identities
		(
			id uuid NOT NULL PRIMARY KEY,
			user_id uuid NOT NULL,
    		password character varying(255) NOT NULL,
			hint character varying(255),
			created_at timestamp without time zone NOT NULL,
			updated_at timestamp without time zone NOT NULL DEFAULT 'NOW()'
		);`, _schema),
		"password_auths": fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.password_auths
		(
			id uuid NOT NULL PRIMARY KEY,
			token character varying(255) NOT NULL,
    		password_id uuid NOT NULL,
			status bool NOT NULL,
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

/*

	// Device


	// // Phone


*/
