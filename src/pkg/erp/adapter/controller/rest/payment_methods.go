package rest

import (
	"auth/src/pkg/account/usecase"
	"auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// CreatePaymentMethod handles the creation of a new payment method.
func (controller Controller) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Create Payment Method Request ||||||||")
	controller.log.Println("Processing Create Payment Method Request")

	// Struct for CreatePaymentMethod request payload
	type CreatePaymentMethodRequest struct {
		Name      string  `json:"name"`
		Type      string  `json:"type"`
		Comission float64 `json:"comission"`
		Details   string  `json:"details"`
		IsActive  bool    `json:"is_active"`
	}

	// Authenticate (AuthN)
	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	// Parse and validate the request body for payment method data
	var req CreatePaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid payment method data.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Usecase [CREATE PAYMENT METHOD]
	err = controller.interactor.CreatePaymentMethod(req.Name, req.Type, req.Comission, req.Details, req.IsActive, session.User.Id)
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

	// Send success response
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Payment method created successfully",
	}, http.StatusCreated)
}

// UpdatePaymentMethod handles updating a specific payment method.
func (controller Controller) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Update Payment Method Request ||||||||")
	controller.log.Println("Processing Update Payment Method Request")

	// Struct for UpdatePaymentMethod request payload
	type UpdatePaymentMethodRequest struct {
		Name      string  `json:"name,omitempty"`
		Type      string  `json:"type,omitempty"`
		Comission float64 `json:"comission"`
		Details   string  `json:"details,omitempty"`
		IsActive  *bool   `json:"is_active,omitempty"`
	}

	// Authenticate (AuthN)
	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	paymentMethodId := r.URL.Query().Get("id")
	paymentMethodID, err := uuid.Parse(paymentMethodId)
	var req UpdatePaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid payment method data.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Usecase [UPDATE PAYMENT METHOD] - Pass fields individually
	err = controller.interactor.UpdatePaymentMethod(paymentMethodID, req.Name, req.Type, req.Comission, req.Details, req.IsActive, session.User.Id)

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

	// Send success response
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Payment method updated successfully",
	}, http.StatusOK)
}

// ListPaymentMethods handles retrieving all payment methods.
func (controller Controller) ListPaymentMethods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Payment Methods Request ||||||||")
	controller.log.Println("Processing List Payment Methods Request")

	// Authenticate (AuthN)
	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	// Usecase [LIST PAYMENT METHODS]
	paymentMethods, err := controller.interactor.ListPaymentMethods(session.User.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusInternalServerError)
		return
	}

	// Send success response with payment methods data
	SendJSONResponse(w, Response{
		Success: true,
		Data:    paymentMethods,
	}, http.StatusOK)
}

// GetPaymentMethod handles retrieving a specific payment method by ID.
func (controller Controller) GetPaymentMethod(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Get Payment Method Request ||||||||")
	controller.log.Println("Processing Get Payment Method Request")

	// Authenticate (AuthN)
	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	paymentMethodId := r.URL.Query().Get("id")
	paymentMethodID, err := uuid.Parse(paymentMethodId)
	// Usecase [GET PAYMENT METHOD]
	paymentMethod, err := controller.interactor.GetPaymentMethod(paymentMethodID, session.User.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusNotFound)
		return
	}

	// Send success response with payment method data
	SendJSONResponse(w, Response{
		Success: true,
		Data:    paymentMethod,
	}, http.StatusOK)
}

// DeactivatePaymentMethod handles deactivating a specific payment method by ID.
func (controller Controller) DeactivatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Deactivate Payment Method Request ||||||||")
	controller.log.Println("Processing Deactivate Payment Method Request")

	// Authenticate (AuthN)
	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	// Retrieve payment method ID from URL parameters
	paymentMethodID := r.URL.Query().Get("id")
	if paymentMethodID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Payment method ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Usecase [DEACTIVATE PAYMENT METHOD]
	err = controller.interactor.DeactivatePaymentMethod(paymentMethodID, session.User.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusInternalServerError)
		return
	}
	// Send success response
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Payment method deactivated successfully",
	}, http.StatusOK)
}
