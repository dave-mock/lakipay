package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (controller Controller) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Create Customer Request ||||||||")
	controller.log.Println("Processing Create Customer Request")

	// Authenticate (AuthN)
	authHeader := r.Header.Get("Authorization")
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Authentication token missing or malformed in header",
			},
		}, http.StatusUnauthorized)
		return
	}
	token := tokenParts[1]

	// Validate token
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "AUTHENTICATION_ERROR",
				Message: "Invalid token or session",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Decode request body
	var req struct {
		MerchantID    string `json:"merchant_id"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		PhoneNumber   string `json:"phone_number,omitempty"`
		Address       string `json:"address,omitempty"`
		DateOfBirth   string `json:"date_of_birth,omitempty"`
		Status        string `json:"status,omitempty"`
		LoyaltyPoints int    `json:"loyalty_points,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Failed to parse request body",
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Field validation
	if req.MerchantID == "" || req.Name == "" || req.Email == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "MerchantID, Name, and Email are required fields",
			},
		}, http.StatusBadRequest)
		return
	}

	// Convert MerchantID to UUID
	merchantUUID, err := uuid.Parse(req.MerchantID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid MerchantID format",
			},
		}, http.StatusBadRequest)
		return
	}

	// Create customer in use case
	customer, err := controller.interactor.CreateCustomer(
		uuid.New(),        // Generate a new customer ID
		req.Name,          // Name
		req.Email,         // Email
		req.PhoneNumber,   // Phone
		req.Address,       // Address
		req.LoyaltyPoints, // LoyaltyPoints
		req.DateOfBirth,   // DateOfBirth
		req.Status,        // Status
		session.User.Id,   // createdBy (from session)
		merchantUUID,      // MerchantID converted to UUID
	)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "CREATION_ERROR",
				Message: "Failed to create customer",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Success response
	SendJSONResponse(w, Response{
		Success: true,
		Data:    customer,
		Message: "Customer Creted Successfully",
	}, http.StatusCreated)
}

func (controller Controller) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Update Customer Request ||||||||")
	controller.log.Println("Processing Update Customer Request")

	// Authenticate (AuthN)
	authHeader := r.Header.Get("Authorization")
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Authentication token missing or malformed in header",
			},
		}, http.StatusUnauthorized)
		return
	}
	token := tokenParts[1]

	// Validate token
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "AUTHENTICATION_ERROR",
				Message: "Invalid token or session",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var req struct {
		Name          string `json:"name,omitempty"`
		Email         string `json:"email,omitempty"`
		Phone         string `json:"phone,omitempty"`
		Address       string `json:"address,omitempty"`
		LoyaltyPoints int    `json:"loyalty_points,omitempty"`
		DateOfBirth   string `json:"date_of_birth,omitempty"`
		Status        string `json:"status,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid customer data",
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update customer in use case
	customerID := session.User.Id
	customer, err := controller.interactor.UpdateCustomer(
		customerID, req.Name, req.Email, req.Phone, req.Address, req.LoyaltyPoints,
		req.DateOfBirth, req.Status, session.User.Id, session.User.Id,
	)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UPDATE_ERROR",
				Message: "Failed to update customer",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Success response with updated customer details
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Customer updated successfully",
		Data:    customer,
	}, http.StatusOK)
}

/*
func (controller Controller) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Create Merchant Request ||||||||")
	controller.log.Println("Processing Create Merchant Request")

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

	// Define the CreateMerchantRequest struct inside the function
	type CreateMerchantRequest struct {
		Name               string `json:"name"`
		BusinessName       string `json:"business_name"`
		RegistrationNumber string `json:"registration_number"`
		Address            string `json:"address"`
		ContactEmail       string `json:"contact_email"`
		ContactPhone       string `json:"contact_phone"`
	}

	var req CreateMerchantRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
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
	controller.log.Printf("log Create Merchant Request ... %+v", req)
	if req.Name == "" || req.BusinessName == "" || req.RegistrationNumber == "" || req.ContactEmail == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "All fields except address and contact phone are required.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Usecase [CREATE MERCHANT]
	merchantID, err := controller.interactor.CreateMerchant(
		session.User.Id,
		req.Name,
		req.BusinessName,
		req.RegistrationNumber,
		req.Address,
		req.ContactEmail,
		req.ContactPhone,
	)
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

	// Send success response with the created merchant ID
	SendJSONResponse(w, Response{
		Success: true,
		Data:    map[string]interface{}{"MerchantID": merchantID},
	}, http.StatusCreated)
}

// ListMerchants lists all merchants.
func (controller Controller) ListMerchants(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Merchants Request ||||||||")
	controller.log.Println("Processing List Merchants Request")

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

	// Usecase [LIST MERCHANTS]
	merchants, err := controller.interactor.ListMerchants(session.User.Id)
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

	// Send success response with the list of merchants
	SendJSONResponse(w, Response{
		Success: true,
		Data:    merchants,
	}, http.StatusOK)
}
// ListMerchantCustomers lists customers for a specific merchant.
func (controller Controller) ListMerchantCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Merchant Customers Request ||||||||")
	controller.log.Println("Processing List Merchant Customers Request")

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

	merchantID := r.URL.Query().Get("id")
	// Usecase [LIST MERCHANT CUSTOMERS]
	customers, err := controller.interactor.ListMerchantCustomers(merchantID, session.User.Id)
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

	// Send success response with the list of customers
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Customers retrieved successfully",
		Data:    customers,
	}, http.StatusOK)
}
// DeactivateMerchant deactivates a specific merchant.
func (controller Controller) DeactivateCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Deactivate Merchant Request ||||||||")
	controller.log.Println("Processing Deactivate Merchant Request")

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

	// Retrieve the merchant ID from the URL parameters
	merchantID := r.URL.Query().Get("id")
	if merchantID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Printf("log Deactivate Merchant Request for MerchantID: %s", merchantID)
	// Usecase [DEACTIVATE MERCHANT]
	err = controller.interactor.DeactivateCustomers(session.User.Id, merchantID)
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

	// Send success response indicating the merchant was deactivated
	SendJSONResponse(w, Response{
		Success: true,
		Data:    map[string]interface{}{"MerchantID": merchantID},
	}, http.StatusOK)
}

// GetMerchant retrieves a specific merchant by ID.
func (controller Controller) GetMerchant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Get Merchant Request ||||||||")
	controller.log.Println("Processing Get Merchant Request")

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

	// Retrieve the merchant ID from the URL parameters
	merchantID := r.URL.Query().Get("id")
	if merchantID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Usecase [GET MERCHANT]
	merchant, err := controller.interactor.GetMerchant(session.User.Id, merchantID)
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

	// Send success response with the retrieved merchant details
	SendJSONResponse(w, Response{
		Success: true,
		Data:    merchant,
	}, http.StatusOK)
}

func (controller Controller) UpdateMerchant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Create Merchant Request ||||||||")
	controller.log.Println("Processing Create Merchant Request")

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

	// Define the CreateMerchantRequest struct inside the function
	type CreateMerchantRequest struct {
		Name               string `json:"name"`
		BusinessName       string `json:"business_name"`
		RegistrationNumber string `json:"registration_number"`
		Address            string `json:"address"`
		ContactEmail       string `json:"contact_email"`
		ContactPhone       string `json:"contact_phone"`
	}

	// Decode the request from the body
	var req CreateMerchantRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
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

	// Log the request data
	controller.log.Printf("log Create Merchant Request ... %+v", req)

	// Validate the request fields
	if req.Name == "" || req.BusinessName == "" || req.RegistrationNumber == "" || req.ContactEmail == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "All fields except address and contact phone are required.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Usecase [CREATE MERCHANT]
	merchantID, err := controller.interactor.CreateMerchant(
		session.User.Id,
		req.Name,
		req.BusinessName,
		req.RegistrationNumber,
		req.Address,
		req.ContactEmail,
		req.ContactPhone,
	)
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

	// Send success response with the created merchant ID
	SendJSONResponse(w, Response{
		Success: true,
		Data:    map[string]interface{}{"MerchantID": merchantID},
	}, http.StatusCreated)
}

// DeactivateMerchant deactivates a specific merchant.
func (controller Controller) DeactivateMerchant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Deactivate Merchant Request ||||||||")
	controller.log.Println("Processing Deactivate Merchant Request")

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

	// Retrieve the merchant ID from the URL parameters
	merchantID := r.URL.Query().Get("id")
	if merchantID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Log the request data
	controller.log.Printf("log Deactivate Merchant Request for MerchantID: %s", merchantID)

	// Usecase [DEACTIVATE MERCHANT]
	err = controller.interactor.DeactivateMerchant(session.User.Id, merchantID)
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

	// Send success response indicating the merchant was deactivated
	SendJSONResponse(w, Response{
		Success: true,
		Data:    map[string]interface{}{"MerchantID": merchantID},
	}, http.StatusOK)
} */

/*
func (controller Controller) GetApiKeys(w http.ResponseWriter, r *http.Request) {

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
	fmt.Println("||| start auth")
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

	private_key, err := controller.interactor.GetApiKeysUsecase(session.User.Id)
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

	type Response2 struct {
		PrivateKey string `json:"private_key"`
	}

	var res Response2
	res.PrivateKey = private_key
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) GetApplyForToken(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

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

	if req.Username == "" || req.Password == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Username and password are required",
			},
		}, http.StatusBadRequest)
		return
	}

	private_key, err := controller.interactor.ApplyForTokenUsecase(req.Username, req.Password)
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

	type Response2 struct {
		Token  string `json:"token"`
		Detail struct {
			Type string `json:"type"`
			Info string `json:"info"`
		} `json:"detail"`
	}

	var res Response2
	res.Token = private_key
	res.Detail.Type = "Bearer Token"
	res.Detail.Info = "Add the Bearer token to the header to authorize."
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)
}

func (controller Controller) GetSendOtp(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var req Request
	fmt.Println("|||||||||||||||||||||||| one")

	token := strings.Split(r.Header.Get("Authorization"), " ")

	if len(token) == 2 {
		req.Token = token[1]
	}

	fmt.Println("|||||||||||||||||||||||| two", req.Token)

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)
	fmt.Println("|||||||||||||||||||||||| 3")

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.SendOtpUsecase(session.User.Id)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: false,
		Data:    token2,
	}, http.StatusOK)

}

func (controller Controller) GetSetFingerPrint(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	var data interface{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&data)

	// fmt.Println("||||||||||||||||||||||||||||||||||| ", data)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error()},
		}, http.StatusBadRequest)
		return
	}

	token2, err := controller.interactor.SendSetFIngerPrintUsecase(session.User.Id, data)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: false,
		Data:    token2,
	}, http.StatusOK)

}

func (controller Controller) GetGenerateChallenge(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var requestOne struct {
		// Username  string `json:"username"`
		// Challenge string `json:"challenge"`
		DeviceID string `json:"device_id"`
	}
	var res struct {
		Challenge string `json:"challenge"`
	}
	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&requestOne)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.SendGenerateChallenge(session.User.Id, requestOne.DeviceID)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	res.Challenge = token2
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) GetVerifySignatureHandler(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var requestOne struct {
		// Username  string `json:"username"`
		Challenge string `json:"challenge"`
		Signature string `json:"signature"`
	}

	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&requestOne)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.GetverifySignature(session.User.Id, requestOne.Challenge, requestOne.Signature)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: false,
		Data:    token2,
	}, http.StatusOK)

}

func (controller Controller) GetstorePublicKeyHandler(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var requestOne struct {
		DeviceID  string `json:"device_id"`
		PublicKey string `json:"public_key"`
	}

	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&requestOne)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.GetstorePublicKeyHandler(requestOne.PublicKey, session.User.Id, requestOne.DeviceID)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    token2,
	}, http.StatusOK)

}
*/
