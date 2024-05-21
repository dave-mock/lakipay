package usecase

import (
	"log"
)

type Usecase struct {
	log  *log.Logger
	repo Repo
}

func New(log *log.Logger, repo Repo) Interactor {
	return Usecase{log, repo}
}
