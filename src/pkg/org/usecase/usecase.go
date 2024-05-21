package usecase

import "log"

// Errors

type Error struct {
	Type    string
	Message string
}

func (err Error) Error() string {
	return err.Message
}

type Usecase struct {
	log        *log.Logger
	repo       Repo
	tinchecker TINChecker
}

func New(log *log.Logger, repo Repo, tinChecker TINChecker) Interactor {
	return Usecase{log: log, repo: repo, tinchecker: tinChecker}
}
