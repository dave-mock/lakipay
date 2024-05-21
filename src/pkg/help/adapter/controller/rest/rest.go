package rest

import (
	"auth/src/pkg/help/usecase"
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

type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	sm         *http.ServeMux
}

func New(log *log.Logger, interactor usecase.Interactor, sm *http.ServeMux) Controller {
	var controller = Controller{log: log, interactor: interactor}

	// [TODO] Routing
	sm.HandleFunc("/help/countries", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetCountries(w, r)
			}
		case http.MethodPost:
			{
				controller.GetAddCountry(w, r)
			}

		}
	})

	sm.HandleFunc("/help/countries/results", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetSearchAndFilterCountries(w, r)
			}

		}
	})

	sm.HandleFunc("/help/supported-countries", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetSupportedCountries(w, r)
			}
		}
	})

	// Currencies
	sm.HandleFunc("/help/currencies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetCurrencies(w, r)
			}
		case http.MethodPost:
			{
				controller.GetAddCurrency(w, r)
			}
		}
	})

	controller.sm = sm

	return controller
}

func SendJSONResponse(w http.ResponseWriter, data interface{}, status int) {
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

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}
