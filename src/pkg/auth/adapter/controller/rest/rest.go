package rest

import (
	"auth/src/pkg/auth/usecase"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
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

// Response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

type User struct {
	Id        uuid.UUID `json:"id"`
	SirName   string    `json:"sir_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

func New(log *log.Logger, sm *http.ServeMux, interactor usecase.Interactor) Controller {

	controller := Controller{log: log, interactor: interactor}

	sm.HandleFunc("/auth/init", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetInitAuth(w, r)
			}
		}
	})

	sm.HandleFunc("/auth/sign-in", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetSignIn(w, r)
			}
		}
	})

	sm.HandleFunc("/auth/sign-up", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetSignUp(w, r)
			}
		}
	})

	sm.HandleFunc("/auth/check", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetCheckSession(w, r)
			}
		}
	})

	sm.HandleFunc("/auth/password", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetSet2FA(w, r)
			}
		}
	})

	sm.HandleFunc("/auth/password/check", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetCheck2FA(w, r)
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
