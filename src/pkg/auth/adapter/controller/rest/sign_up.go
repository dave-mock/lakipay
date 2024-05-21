package rest

import (
	"auth/src/pkg/auth/usecase"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (controller Controller) GetSignUp(w http.ResponseWriter, r *http.Request) {
	// Request
	type Request struct {
		Token string `json:"token"`
		Phone struct {
			Prefix string `json:"prefix"`
			Number string `json:"number"`
		} `json:"phone"`
		SirName   string `json:"sir_name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	// Response
	type Response struct {
		Success  bool        `json:"success"`
		NextStep string      `json:"next_step,omitempty"`
		Data     interface{} `json:"data,omitempty"`
		Error    error       `json:"error,omitempty"`
	}

	// Parse request
	var req Request

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Usecase operations
	// create user
	_, err = controller.interactor.CreateUser(req.Token, req.SirName, req.FirstName, req.LastName)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Create session
	log.Println("requesting for session creation")
	session, at, err := controller.interactor.CreateSession(req.Token)
	if err != nil {
		if err, ok := err.(usecase.Error); ok {
			switch err.Type {
			case "SET_PASSWORD":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "SET_PASSWORD",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			case "CHECK_PASSWORD":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "CHECK_PASSWORD",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			case "SIGN_UP":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "SIGN_UP",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			}
		}
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "UNSPECIFIED",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Return Response
	SendJSONResponse(w, Response{
		Success: true,
		Data: AuthResponse{
			Token: &struct {
				Active  string "json:\"active\""
				Refresh string "json:\"refresh\""
			}{
				Active:  at,
				Refresh: session.Token,
			},
			User: &struct {
				Id        uuid.UUID "json:\"id\""
				SirName   string    "json:\"sir_name,omitempty\""
				FirstName string    "json:\"first_name\""
				LastName  string    "json:\"last_name,omitempty\""
			}{
				Id:        session.User.Id,
				SirName:   session.User.SirName,
				FirstName: session.User.FirstName,
				LastName:  session.User.LastName,
			},
		},
	}, http.StatusOK)
}

/*



 */
