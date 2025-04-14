package main

import (
	"log"
	"net/http"
	"os"

	// HTTP Server

	lhttp "auth/src/pkg/auth/infra/network/http"
	// Postgres Storage
	"auth/src/pkg/auth/infra/storage/psql"
	// SMS
	"auth/src/pkg/auth/adapter/gateway/sms"
	// [AUTH]
	"auth/src/pkg/auth/adapter/controller/procedure"
	auth "auth/src/pkg/auth/adapter/controller/rest"
	authRepo "auth/src/pkg/auth/adapter/gateway/repo/psql"
	authUsecase "auth/src/pkg/auth/usecase"

	// [HELP]
	help "auth/src/pkg/help/adapter/controller/rest"
	helpRepo "auth/src/pkg/help/adapter/gateway/repo/psql"
	helpUsecase "auth/src/pkg/help/usecase"

	// // [ORG]
	org "auth/src/pkg/org/adapter/controller/rest"
	orgRepo "auth/src/pkg/org/adapter/gateway/repo/psql"
	"auth/src/pkg/org/adapter/gateway/tin_checker/etrade"
	orgUsecase "auth/src/pkg/org/usecase"

	// [ACCOUNT]
	acc "auth/src/pkg/account/adapter/controller/rest"
	accRepo "auth/src/pkg/account/adapter/gateway/repo/psql"
	accUsecase "auth/src/pkg/account/usecase"

	// [ERP]

	erp "auth/src/pkg/erp/adapter/controller/rest"
	erpRepo "auth/src/pkg/erp/adapter/gateway/repo/psql"
	erpUsecase "auth/src/pkg/erp/usecase"

	// [ACCESS CONTROL]

	access_control "auth/src/pkg/access_control/adapter/controller/rest"
	accessRepo "auth/src/pkg/access_control/adapter/gateway/repo/psql"
	accessUsecase "auth/src/pkg/access_control/usecase"

	// [STORAGE]
	storage "auth/src/pkg/storage/adapter/controller/rest"
	storageRepo "auth/src/pkg/storage/adapter/gateway/repo/psql"
	storageUsecase "auth/src/pkg/storage/usecase"

	// [CHECKOUT]
	checkout "auth/src/pkg/checkout/adapter/controller/rest"
	checkoutRepo "auth/src/pkg/checkout/adapter/gateway/repo/psql"
	checkoutService "auth/src/pkg/checkout/service/basic"
)

func main() {

	log := log.New(os.Stdout, "[LAKIPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	// [Ouput Adapters]
	// [DB] Postgres
	db, err := psql.New(log)
	if err != nil {
		log.Fatal(err.Error())
	}

	// [Input Adapters]
	// [SERVER] HTTP
	s := lhttp.New(log)
	defer s.Serve()

	s.ServeMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	// [Module Adapters]

	// [AUTH]
	_authRepo, err := authRepo.NewPsqlRepo(log, db)
	if err != nil {
		log.Println(err)
	} else {
		auth.New(log, s.ServeMux, authUsecase.New(log, _authRepo, sms.New(log)))
	}

	// [HELP]
	helpRepo, err := helpRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		help.New(log, helpUsecase.New(log, helpRepo), s.ServeMux)
	}

	// [ORG]
	_orgRepo, err := orgRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		org.New(log, orgUsecase.New(log, _orgRepo, etrade.New(log)), s.ServeMux)
	}

	// [STORAGE]
	_storageRepo, err := storageRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		storage.New(log, s.ServeMux, storageUsecase.New(log, _storageRepo))
	}

	// [ACCOUNT]
	var _accRepo accUsecase.Repo
	var accService accUsecase.Interactor
	_accRepo, err = accRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		accService = accUsecase.New(log, _accRepo)
		acc.New(log, accService, s.ServeMux, procedure.New(log, authUsecase.New(log, _authRepo, sms.New(log))))
	}
	// ERP
	_erpRepo, err := erpRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		erp.New(log, erpUsecase.New(log, _erpRepo), s.ServeMux, procedure.New(log, authUsecase.New(log, _authRepo, sms.New(log))))
	}

	// ACCESS CONTROL
	_accessRepo, err := accessRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		access_control.New(log, accessUsecase.New(log, _accessRepo), s.ServeMux, procedure.New(log, authUsecase.New(log, _authRepo, sms.New(log))))
	}

	// Checkout
	_checkoutRepo, err := checkoutRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		checkout.New(log, s.ServeMux, checkoutService.New(log, _checkoutRepo, accService))
	}
}
