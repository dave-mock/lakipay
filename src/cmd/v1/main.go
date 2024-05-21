package main

import (
	"log"
	"os"

	// HTTP Server

	"auth/src/pkg/auth/infra/network/http"
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

	// [STORAGE]
	storage "auth/src/pkg/storage/adapter/controller/rest"
	storageRepo "auth/src/pkg/storage/adapter/gateway/repo/psql"
	storageUsecase "auth/src/pkg/storage/usecase"
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
	s := http.New(log)
	defer s.Serve()

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
	_accRepo, err := accRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		acc.New(log, accUsecase.New(log, _accRepo), s.ServeMux, procedure.New(log, authUsecase.New(log, _authRepo, sms.New(log))))
	}

}
