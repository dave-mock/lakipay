package procedure

import (
	"auth/src/pkg/auth/usecase"
	"log"
)

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}

type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
}

func New(log *log.Logger, interactor usecase.Interactor) Controller {
	return Controller{log: log, interactor: interactor}
}
