package rest

import (
	"auth/src/pkg/auth/usecase"
	"auth/src/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (controller Controller) GetSet2FA(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
		Hint     string `json:"hint"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Usecase
	// Creare 2FA
	_, err = controller.interactor.InitPasswordAuth(req.Token, req.Password, req.Hint)
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

	//
	// Usecase

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
				UserType  string    "json:\"user_type,omitempty\""
			}{
				Id:        session.User.Id,
				SirName:   session.User.SirName,
				FirstName: session.User.FirstName,
				LastName:  session.User.LastName,
				UserType:  session.User.UserType,
			},
		},
	}, http.StatusOK)
}

func (controller Controller) GetCheck2FA(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Usecase
	// Check Password
	err = controller.interactor.AuthPassword(req.Token, req.Password)
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
			log.Println(err.Type)
			log.Println(err.Message)
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
				UserType  string    "json:\"user_type,omitempty\""
			}{
				Id:        session.User.Id,
				SirName:   session.User.SirName,
				FirstName: session.User.FirstName,
				LastName:  session.User.LastName,
				UserType:  session.User.UserType,
			},
		},
	}, http.StatusOK)
}

func (controller Controller) GetDecryptData(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Data string `json:"data"`
	}
	fmt.Println("###################################################")
	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	fmt.Println("################################################### , two")

	decryptedData, err := utils.AesDecription(req.Data)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return

	}

	parts := strings.Split(decryptedData, ",")

	merchantIdString := parts[0]

	parsedUUID, err := uuid.Parse(merchantIdString)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}
	fmt.Println("################################################### , three", parsedUUID)

	user, err := controller.interactor.GetUserById(parsedUUID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}
	fmt.Println("################################################### , four")

	type Response2 struct {
		Data     string `json:"data"`
		FullName string `json:"full_name"`
	}

	var res Response2

	res.FullName = user.SirName + " " + user.FirstName + " " + user.LastName
	res.Data = decryptedData

	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}
