package psql

import (
	"auth/src/pkg/checkout/service"
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

type CheckoutPSQLRepo struct {
	log *log.Logger
	db  *sql.DB
}

func New(log *log.Logger, db *sql.DB) (service.CheckoutRepo, error) {
	log.SetPrefix("[AUTH] [GATEWAY] [REPO] [PSQL] ")

	_initFile, err := filepath.Abs("src/pkg/checkout/adapter/gateway/repo/psql/init.sql")
	if err != nil {
		return nil, err
	}
	script, err := os.ReadFile(_initFile)
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(string(script))
	if err != nil {
		log.Println(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return CheckoutPSQLRepo{log, db}, nil
}
