package basic

import (
	account "auth/src/pkg/account/usecase"
	"auth/src/pkg/checkout/service"
	"log"
)

type BasicCheckoutService struct {
	log     *log.Logger
	repo    service.CheckoutRepo
	account account.Interactor
}

func New(log *log.Logger, repo service.CheckoutRepo, account account.Interactor) service.CheckoutInteractor {
	return BasicCheckoutService{log, repo, account}
}
