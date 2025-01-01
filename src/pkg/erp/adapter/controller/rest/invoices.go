package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (controller Controller) CreateMerchantInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| ||  Create Invoice Request ||||||||")
	controller.log.Println("Processing create Invoice Request")
	authorizationHeader := r.Header.Get("Authorization")
	if len(strings.Split(authorizationHeader, " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorizationHeader, " ")[1]
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		controller.log.Printf("Error authenticating token: %v", err)

		var authErr procedure.Error
		if ok := errors.As(err, &authErr); ok {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    authErr.Type,
					Message: authErr.Message,
				},
			}, http.StatusUnauthorized)
			return
		}

		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "AUTHENTICATION_ERROR",
				Message: "Failed to authenticate the token.",
			},
		}, http.StatusUnauthorized)
		return
	}
	var requestBody struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		controller.log.Printf("Error parsing request body: %v", err)

		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid JSON body format.",
			},
		}, http.StatusBadRequest)
		return
	}

	if requestBody.ID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Order ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	OrderID, err := uuid.Parse(requestBody.ID)
	if err != nil {
		controller.log.Printf("Error parsing Order ID: %v", err)

		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Order ID format.",
			},
		}, http.StatusBadRequest)
		return
	}
	orders, err := controller.interactor.CreateMerchantInvoice(session.User.Id, OrderID)
	if err != nil {
		controller.log.Printf("Error creating invoice: %v", err)
		if strings.Contains(err.Error(), "no orders found") {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "ORDER_NOT_FOUND",
					Message: "No orders found for the provided order ID",
				},
			}, http.StatusNotFound)
			return
		}

		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "An unexpected error occurred while creating the invoice.",
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Invoice Created successfully",
		Data:    orders,
	}, http.StatusOK)
}
