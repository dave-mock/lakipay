package rest

import (
	// auth "auth/src/pkg/auth/adapter/controller/procedure"
	"auth/src/pkg/storage/usecase"
	"log"
	"net/http"
)

type Controller struct {
	log        *log.Logger
	sm         *http.ServeMux
	interactor usecase.Interactor
	// auth       auth.Controller
}

func New(log *log.Logger, sm *http.ServeMux, interactor usecase.Interactor) Controller {
	controller := Controller{log: log, interactor: interactor}

	// Route
	sm.HandleFunc("/storage", func(w http.ResponseWriter, r *http.Request) {
		controller.log.Println("Storage")
		switch r.Method {
		case http.MethodPost:
			{
				controller.Upload(w, r)
			}
		}
	})

	controller.sm = sm

	return controller
}
