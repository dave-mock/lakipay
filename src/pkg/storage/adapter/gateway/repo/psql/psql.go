package psql

import (
	"auth/src/pkg/storage/usecase"
	"database/sql"
	"log"
)

type PsqlRepo struct {
	log *log.Logger
	db  *sql.DB
}

func New(log *log.Logger, db *sql.DB) (usecase.Repo, error) {
	return PsqlRepo{log: log, db: db}, nil
}
