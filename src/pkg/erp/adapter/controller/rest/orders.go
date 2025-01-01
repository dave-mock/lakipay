package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"auth/src/pkg/erp/core/entity"
	"auth/src/pkg/erp/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (controller Controller) CreateOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || Handle Create Order Request ||||||||")
	controller.log.Println("Processing Create Order Request")

	authHeader := r.Header.Get("Authorization")
	if len(strings.Split(authHeader, " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Authentication token missing in header.",
			},
		}, http.StatusUnauthorized)
		return
	}

	token := strings.Split(authHeader, " ")[1]
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

	var createOrderRequest struct {
		OrderTypeID     string             `json:"order_type_id"`
		TotalAmount     float64            `json:"total_amount"`
		Currency        string             `json:"currency"`
		Medium          string             `json:"medium"`
		ShippingAddress string             `json:"shipping_address"`
		BillingAddress  string             `json:"billing_address"`
		OrderItems      []entity.OrderItem `json:"order_items"`
		Discounts       []entity.Discount  `json:"discounts"`
		Taxes           []entity.Tax       `json:"taxes"`
	}
	/* requiredPermission := entity.Permission{
		Resource:           "warehouses",
		Operation:          "delete",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to create an order.",
			},
		}, http.StatusForbidden)
		return
	} */
	if err := json.NewDecoder(r.Body).Decode(&createOrderRequest); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid order data.",
			},
		}, http.StatusBadRequest)
		return
	}

	createdOrder, err := controller.interactor.CreateOrder(
		session.User.Id,
		session.User.Id.String(),
		createOrderRequest.OrderTypeID,
		createOrderRequest.TotalAmount,
		createOrderRequest.Currency,
		createOrderRequest.Medium,
		createOrderRequest.ShippingAddress,
		createOrderRequest.BillingAddress,
		createOrderRequest.OrderItems,
		createOrderRequest.Discounts,
		createOrderRequest.Taxes,
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

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Order created successfully",
		Data:    createdOrder,
	}, http.StatusOK)
}

func (controller Controller) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || Handle Update Order Request ||||||||")
	controller.log.Println("Processing Update Order Request")

	authHeader := r.Header.Get("Authorization")
	if len(strings.Split(authHeader, " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Authentication token missing in header.",
			},
		}, http.StatusUnauthorized)
		return
	}

	token := strings.Split(authHeader, " ")[1]
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
	/* requiredPermission := entity.Permission{
		Resource:           "warehouses",
		Operation:          "delete",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to update an order.",
			},
		}, http.StatusForbidden)
		return
	}
	*/
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Order ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	var updateOrderRequest struct {
		MerchantID      string             `json:"merchant_id"`
		OrderTypeID     string             `json:"order_type_id"`
		TotalAmount     float64            `json:"total_amount"`
		Currency        string             `json:"currency"`
		Medium          string             `json:"medium"`
		ShippingAddress string             `json:"shipping_address"`
		BillingAddress  string             `json:"billing_address"`
		OrderItems      []entity.OrderItem `json:"order_items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateOrderRequest); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid order update data.",
			},
		}, http.StatusBadRequest)
		return
	}

	merchantID, err := uuid.Parse(updateOrderRequest.MerchantID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Merchant ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	updatedOrder, err := controller.interactor.UpdateOrder(
		orderID,
		merchantID,
		session.User.Id,
		updateOrderRequest.OrderTypeID,
		updateOrderRequest.TotalAmount,
		updateOrderRequest.Currency,
		updateOrderRequest.Medium,
		updateOrderRequest.ShippingAddress,
		updateOrderRequest.BillingAddress,
		updateOrderRequest.OrderItems,
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

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Order updated successfully",
		Data:    updatedOrder,
	}, http.StatusOK)
}

func (controller Controller) CancelOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Cancel Order Request ||||||||")
	controller.log.Println("Processing Cancel Order Request")

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
	/* requiredPermission := entity.Permission{
		Resource:           "warehouses",
		Operation:          "delete",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to cancel an order.",
			},
		}, http.StatusForbidden)
		return
	} */
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Order ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	err = controller.interactor.CancelOrder(session.User.Id, orderID)
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

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Order canceled successfully",
	}, http.StatusOK)
}

func (controller Controller) CreateCartOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || Handle Add Item to Cart Request ||||||||")
	controller.log.Println("Processing Create Order Request")

	authHeader := r.Header.Get("Authorization")
	if len(strings.Split(authHeader, " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Authentication token missing in header.",
			},
		}, http.StatusUnauthorized)
		return
	}

	token := strings.Split(authHeader, " ")[1]
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
	/* requiredPermission := entity.Permission{
		Resource:           "warehouses",
		Operation:          "delete",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to create an cart order.",
			},
		}, http.StatusForbidden)
		return
	} */
	var createOrderRequest struct {
		OrderTypeID     string             `json:"order_type_id"`
		TotalAmount     float64            `json:"total_amount"`
		Currency        string             `json:"currency"`
		Medium          string             `json:"medium"`
		ShippingAddress string             `json:"shipping_address"`
		BillingAddress  string             `json:"billing_address"`
		OrderItems      []entity.OrderItem `json:"order_items"`
		Discounts       []entity.Discount  `json:"discounts"`
		Taxes           []entity.Tax       `json:"taxes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createOrderRequest); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid order data.",
			},
		}, http.StatusBadRequest)
		return
	}

	createdOrder, err := controller.interactor.CreateOrder(
		session.User.Id,
		session.User.Id.String(),
		createOrderRequest.OrderTypeID,
		createOrderRequest.TotalAmount,
		createOrderRequest.Currency,
		createOrderRequest.Medium,
		createOrderRequest.ShippingAddress,
		createOrderRequest.BillingAddress,
		createOrderRequest.OrderItems,
		createOrderRequest.Discounts,
		createOrderRequest.Taxes,
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

	SendJSONResponse(w, Response{
		Success: true,
		Message: "|||||||  Order added to cart |||||||",
		Data:    createdOrder,
	}, http.StatusOK)
}

func (controller Controller) UpdateCartOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || Handle Update cart request ||||||||")
	controller.log.Println("Processing Update Order Request")

	authHeader := r.Header.Get("Authorization")
	if len(strings.Split(authHeader, " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Authentication token missing in header.",
			},
		}, http.StatusUnauthorized)
		return
	}

	token := strings.Split(authHeader, " ")[1]
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
	/* requiredPermission := entity.Permission{
		Resource:           "warehouses",
		Operation:          "delete",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to update the cart.",
			},
		}, http.StatusForbidden)
		return
	}
	*/
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Order ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	var updateOrderRequest struct {
		MerchantID      string             `json:"merchant_id"`
		OrderTypeID     string             `json:"order_type_id"`
		TotalAmount     float64            `json:"total_amount"`
		Currency        string             `json:"currency"`
		Medium          string             `json:"medium"`
		ShippingAddress string             `json:"shipping_address"`
		BillingAddress  string             `json:"billing_address"`
		OrderItems      []entity.OrderItem `json:"order_items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateOrderRequest); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid order update data.",
			},
		}, http.StatusBadRequest)
		return
	}

	merchantID, err := uuid.Parse(updateOrderRequest.MerchantID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Merchant ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	updatedOrder, err := controller.interactor.UpdateOrder(
		orderID,
		merchantID,
		session.User.Id,
		updateOrderRequest.OrderTypeID,
		updateOrderRequest.TotalAmount,
		updateOrderRequest.Currency,
		updateOrderRequest.Medium,
		updateOrderRequest.ShippingAddress,
		updateOrderRequest.BillingAddress,
		updateOrderRequest.OrderItems,
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

	SendJSONResponse(w, Response{
		Success: true,
		Message: "cart updated successfully",
		Data:    updatedOrder,
	}, http.StatusOK)
}

func (controller Controller) CancelCartOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Cancel Order Request ||||||||")
	controller.log.Println("Processing Cancel Order Request")

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
	/* 	requiredPermission := entity.Permission{
	   		Resource:           "warehouses",
	   		Operation:          "delete",
	   		ResourceIdentifier: "*",
	   		Effect:             "allow",
	   	}
	   	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	   	if err != nil || !hasPermission {
	   		SendJSONResponse(w, Response{
	   			Success: false,
	   			Error: &Error{
	   				Type:    "FORBIDDEN",
	   				Message: "You do not have permission to cancel an cart order.",
	   			},
	   		}, http.StatusForbidden)
	   		return
	   	} */
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Order ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Pass the order ID to the use case layer for cancellation
	err = controller.interactor.CancelOrder(session.User.Id, orderID)
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
		Message: "Order canceled successfully",
	}, http.StatusOK)
}
func (controller Controller) ListOrderItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Order Items Request ||||||||")
	controller.log.Println("Processing List Order Items Request")
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
	/* requiredPermission := entity.Permission{
		Resource:           "orders",
		Operation:          "read",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to list order items.",
			},
		}, http.StatusForbidden)
		return
	} */
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Order ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	items, err := controller.interactor.ListOrderItems(session.User.Id, orderID)
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

	// Send success response with items data
	SendJSONResponse(w, Response{
		Success: true,
		Data:    items,
		Message: "Order items retrieved successfully",
	}, http.StatusOK)
}

func (controller Controller) ListMerchantOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Merchant Orders Request ||||||||")
	controller.log.Println("Processing List Merchant Orders Request")

	// Authenticate (AuthN)
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

	// Validate token
	token := strings.Split(authorizationHeader, " ")[1]
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
	/* requiredPermission := entity.Permission{
		Resource:           "orders",
		Operation:          "read",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to merchant orders.",
			},
		}, http.StatusForbidden)
		return
	} */
	orders, err := controller.interactor.ListMerchantOrders(session.User.Id)
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

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Orders retrieved successfully",
		Data:    orders,
	}, http.StatusOK)
}
func (controller Controller) CountMerchantOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Count Merchant Orders Request ||||||||")
	controller.log.Println("Processing Count Merchant Orders Request")
	authorizationHeader := r.Header.Get("Authorization")
	if len(strings.Split(authorizationHeader, " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in the header",
			},
		}, http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorizationHeader, " ")[1]
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

	/* 	requiredPermission := entity.Permission{
	   		Resource:           "orders",
	   		Operation:          "read",
	   		ResourceIdentifier: "*",
	   		Effect:             "allow",
	   	}
	   	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	   	if err != nil || !hasPermission {
	   		SendJSONResponse(w, Response{
	   			Success: false,
	   			Error: &Error{
	   				Type:    "FORBIDDEN",
	   				Message: "You do not have permission to count orders.",
	   			},
	   		}, http.StatusForbidden)
	   		return
	   	} */
	totalOrdersCount, err := controller.interactor.CountMerchantOrders(session.User.Id)
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

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Total Orders retrieved successfully",
		Data: map[string]int{
			"total_orders": totalOrdersCount,
		},
	}, http.StatusOK)
}

func (controller Controller) CountMerchantCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Count Merchant Orders Request ||||||||")
	controller.log.Println("Processing Count Merchant Orders Request")

	// Authenticate (AuthN)
	authorizationHeader := r.Header.Get("Authorization")
	if len(strings.Split(authorizationHeader, " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in the header",
			},
		}, http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorizationHeader, " ")[1]
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

	/* requiredPermission := entity.Permission{
		Resource:           "orders",
		Operation:          "read",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to count merchant customers.",
			},
		}, http.StatusForbidden)
		return
	}
	*/
	// Fetch total customer count
	totalCustomerCount, err := controller.interactor.CountMerchantCustomers(session.User.Id)
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

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Total Customers retrieved successfully",
		Data: map[string]int{
			"total_customers": totalCustomerCount,
		},
	}, http.StatusOK)
}
func (controller Controller) GetMerchantCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Get Merchant Customers Request ||||||||")
	controller.log.Println("Processing Get Merchant Customers Request")

	// Authenticate (AuthN)
	authorizationHeader := r.Header.Get("Authorization")
	if len(strings.Split(authorizationHeader, " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in the header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(authorizationHeader, " ")[1]
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
	/*
		requiredPermission := entity.Permission{
			Resource:           "orders",
			Operation:          "read",
			ResourceIdentifier: "*",
			Effect:             "allow",
		}
		hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
		if err != nil || !hasPermission {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "FORBIDDEN",
					Message: "You do not have permission to list merchant orders.",
				},
			}, http.StatusForbidden)
			return
		} */
	customerDetailsList, err := controller.interactor.ListMerchantCustomers(session.User.Id)
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

	// Return the list of customers
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Merchant customers retrieved successfully",
		Data: map[string]interface{}{
			"customer_details": customerDetailsList,
		},
	}, http.StatusOK)
}

func (controller Controller) ListOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Orders Request ||||||||")
	controller.log.Println("Processing List Orders Request")

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
	/*
		requiredPermission := entity.Permission{
			Resource:           "orders",
			Operation:          "read",
			ResourceIdentifier: "*",
			Effect:             "allow",
		}
		hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
		if err != nil || !hasPermission {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "FORBIDDEN",
					Message: "You do not have permission to list all orders.",
				},
			}, http.StatusForbidden)
			return
		} */
	fmt.Println(session.User.Id)
	orders, err := controller.interactor.ListOrders()
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

	SendJSONResponse(w, Response{
		Success: true,
		Data:    orders,
		Message: "Orders retrieved successfully",
	}, http.StatusOK)
}

func (controller Controller) GetOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Get Order Request ||||||||")
	controller.log.Println("Processing Get Order Request")
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
	/* requiredPermission := entity.Permission{
		Resource:           "orders",
		Operation:          "read",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to get single orders.",
			},
		}, http.StatusForbidden)
		return
	} */
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Order ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	order, err := controller.interactor.GetOrder(orderID, session.User.Id)
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

	// Send success response with order data
	SendJSONResponse(w, Response{
		Success: true,
		Data:    order,
		Message: "Order retrieved successfully",
	}, http.StatusOK)
}
func (controller Controller) ListCartOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Orders Request ||||||||")
	controller.log.Println("Processing List Orders Request")
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

	/* 	requiredPermission := entity.Permission{
	   		Resource:           "orders",
	   		Operation:          "read",
	   		ResourceIdentifier: "*",
	   		Effect:             "allow",
	   	}
	   	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	   	if err != nil || !hasPermission {
	   		SendJSONResponse(w, Response{
	   			Success: false,
	   			Error: &Error{
	   				Type:    "FORBIDDEN",
	   				Message: "You do not have permission to list cart orders.",
	   			},
	   		}, http.StatusForbidden)
	   		return
	   	} */
	fmt.Println(session.User.Id)
	orders, err := controller.interactor.ListOrders()
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

	// Send success response with order data
	SendJSONResponse(w, Response{
		Success: true,
		Data:    orders,
		Message: "Orders retrieved successfully",
	}, http.StatusOK)
}

func (controller Controller) GetCartOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Get Order Request ||||||||")
	controller.log.Println("Processing Get Order Request")

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
	/* requiredPermission := entity.Permission{
		Resource:           "orders",
		Operation:          "read",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to get cart orders.",
			},
		}, http.StatusForbidden)
		return
	} */

	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Order ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}
	order, err := controller.interactor.GetOrder(orderID, session.User.Id)
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
	SendJSONResponse(w, Response{
		Success: true,
		Data:    order,
		Message: "Order retrieved successfully",
	}, http.StatusOK)
}

func (controller Controller) UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Update Order Item Request ||||||||")
	controller.log.Println("Processing Update Order Item Request")
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
	/* requiredPermission := entity.Permission{
		Resource:           "orders",
		Operation:          "update",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to update order items.",
			},
		}, http.StatusForbidden)
		return
	} */

	ordersID := r.URL.Query().Get("id")
	itemsID := mux.Vars(r)["item_id"]
	orderID, err := uuid.Parse(ordersID)
	var itemUpdate struct {
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
		Discount float64 `json:"discount,omitempty"`
		Tax      float64 `json:"tax,omitempty"`
	}
	itemID, err := uuid.Parse(itemsID)
	if err := json.NewDecoder(r.Body).Decode(&itemUpdate); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid order item update data.",
			},
		}, http.StatusBadRequest)
		return
	}

	err = controller.interactor.UpdateOrderItem(
		orderID,
		itemID,
		session.User.Id,
		itemUpdate.Quantity,
		itemUpdate.Price,
		itemUpdate.Discount,
		itemUpdate.Tax,
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

	// Send success response
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Order item updated successfully",
	}, http.StatusOK)
}

func (controller Controller) RemoveOrderItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Remove Order Item Request ||||||||")
	controller.log.Println("Processing Remove Order Item Request")
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
	/* requiredPermission := entity.Permission{
		Resource:           "orders",
		Operation:          "read",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to update order items.",
			},
		}, http.StatusForbidden)
		return
	} */
	ordersID := r.URL.Query().Get("id")
	itemsID := mux.Vars(r)["item_id"]
	orderID, err := uuid.Parse(ordersID)
	itemID, err := uuid.Parse(itemsID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Item ID format.",
			},
		}, http.StatusBadRequest)
		return
	}
	err = controller.interactor.RemoveOrderItem(orderID, itemID, session.User.Id)
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
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Order item removed successfully",
	}, http.StatusOK)
}
