package rest

import (
	"auth/src/pkg/checkout/service"
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}

// Response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

func SendJSONResponse(w http.ResponseWriter, data Response, status int) {
	serData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(serData)
}

type Controller struct {
	log        *log.Logger
	sm         *http.ServeMux
	interactor service.CheckoutInteractor
}

func New(log *log.Logger, sm *http.ServeMux, interactor service.CheckoutInteractor) Controller {

	var controller Controller = Controller{log: log, interactor: interactor}

	// [ROUTING]
	//
	// [GATEWAYS]
	sm.HandleFunc("/api/v1/checkout/gateways", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetGateways(w, r)
			}
		}
	})

	// [TRANSACTION]
	sm.HandleFunc("/api/v1/checkout/init", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetInitTransaction(w, r)
			}
		}
	})
	sm.HandleFunc("/api/v1/checkout/process", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetConfirmTransaction(w, r)
			}
		}
	})
	sm.HandleFunc("/api/v1/checkout/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetTransaction(w, r)
			}
		}
	})

	sm.HandleFunc("/api/v1/checkout/transactions/notify", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetHandleTransactionNotification(w, r)
				return
			}
		}
	})

	controller.sm = sm

	return controller
}
