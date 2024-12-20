package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"auth/src/pkg/auth/core/entity"
	"auth/src/pkg/auth/usecase"

	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (controller Controller) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Create Warehouse Request ||||||||")
	controller.log.Println("Processing Create Warehouse Request")
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

	type CreateWarehouseRequest struct {
		Name        string `json:"name"`
		Location    string `json:"location"`
		Capacity    int    `json:"capacity"`
		MerchantID  string `json:"merchant_id"`
		Description string `json:"description"`
	}

	var req CreateWarehouseRequest
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
	controller.log.Printf("|||||||| LOOOOG Create Warehouse Request ... %+v", req)
	if req.Name == "" || req.Location == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Name, Location, and Merchant ID are required.",
			},
		}, http.StatusBadRequest)
		return
	}

	requiredPermission := entity.Permission{
		Resource:           "warehouses",
		Operation:          "create",
		ResourceIdentifier: "*",
		Effect:             "allow",
	}
	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to create a warehouse.",
			},
		}, http.StatusForbidden)
		return
	}
	warehouseID, err := controller.interactor.CreateWarehouse(
		req.Name,
		req.Location,
		req.Capacity,
		req.Description,
		session.User.Id,
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
		Data:    map[string]interface{}{"WarehouseID": warehouseID},
	}, http.StatusCreated)
}

func (controller Controller) ListWarehouses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Warehouses Request ||||||||")
	controller.log.Println("Processing List Warehouses Request")
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

	requiredPermission := entity.Permission{
		Resource:           "warehouses",
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
				Message: "You do not have permission to list a warehouse.",
			},
		}, http.StatusForbidden)
		return
	}
	warehouses, err := controller.interactor.ListWarehouses(session.User.Id)
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

	controller.log.Printf("Retrieved warehouses: %+v", warehouses)
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Warehouses retrieved successfully",
		Data:    warehouses,
	}, http.StatusOK)
}

func (controller Controller) ListMerchantWarehouses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Merchant Warehouses Request ||||||||")
	controller.log.Println("Processing List Merchant Warehouses Request")
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

	type ListMerchantWarehousesRequest struct {
		MerchantID string `json:"merchant_id"`
	}

	var payload ListMerchantWarehousesRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid request payload",
			},
		}, http.StatusBadRequest)
		return
	}

	if payload.MerchantID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	requiredPermission := entity.Permission{
		Resource:           "warehouses",
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
				Message: "You do not have permission to see a warehouse.",
			},
		}, http.StatusForbidden)
		return
	}
	warehouses, err := controller.interactor.ListMerchantWarehouses(session.User.Id)
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

	controller.log.Printf("Retrieved warehouses for merchant ID %s: %+v", payload.MerchantID, warehouses)
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Warehouses retrieved successfully",
		Data:    warehouses,
	}, http.StatusOK)
}

func (controller Controller) GetWarehouse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Get Warehouse Request ||||||||")
	controller.log.Println("Processing Get Warehouse Request")
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

	type GetWarehouseRequest struct {
		WarehouseID string `json:"warehouse_id"`
	}

	var payload GetWarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid request payload",
			},
		}, http.StatusBadRequest)
		return
	}

	if payload.WarehouseID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Warehouse ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	GetWarehouseID, err := uuid.Parse(payload.WarehouseID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Warehouse ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	requiredPermission := entity.Permission{
		Resource:           "warehouses",
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
				Message: "You do not have permission to list and see warehouse",
			},
		}, http.StatusForbidden)
		return
	}

	warehouse, err := controller.interactor.GetWarehouse(GetWarehouseID, session.User.Id)
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

	controller.log.Printf("||||||| Retrieved warehouse||||||| %+v", warehouse)
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Warehouse retrieved successfully",
		Data:    warehouse,
	}, http.StatusOK)
}

func (controller Controller) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Update Warehouse Request ||||||||")
	controller.log.Println("Processing Update Warehouse Request")
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

	requiredPermission := entity.Permission{
		Resource:           "warehouses",
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
				Message: "You do not have permission to update a warehouse.",
			},
		}, http.StatusForbidden)
		return
	}

	type UpdateWarehouseRequest struct {
		WarehouseID string `json:"warehouse_id"`
		Name        string `json:"name,omitempty"`
		Location    string `json:"location,omitempty"`
		Capacity    int    `json:"capacity,omitempty"`
		Description string `json:"description,omitempty"`
	}

	var req UpdateWarehouseRequest
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

	if req.WarehouseID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Warehouse ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	warehouseID, err := uuid.Parse(req.WarehouseID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Warehouse ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Printf("Updating Warehouse with ID: %+v and data: %+v", warehouseID, req)
	updatedWarehouse, err := controller.interactor.UpdateWarehouse(warehouseID, req.Name, req.Location, req.Capacity, req.Description, session.User.Id)
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
		Message: "Warehouse updated successfully",
		Data:    updatedWarehouse,
	}, http.StatusOK)
}

func (controller Controller) DeactivateWarehouse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Deactivate Warehouse Request ||||||||")
	controller.log.Println("Processing Deactivate Warehouse Request")
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

	type DeactivateWarehouseRequest struct {
		WarehouseID string `json:"warehouse_id"`
	}
	var req DeactivateWarehouseRequest
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
	if req.WarehouseID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Warehouse ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}
	DeactivateWarehouseID, err := uuid.Parse(req.WarehouseID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Warehouse ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	requiredPermission := entity.Permission{
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
				Message: "You do not have permission to deactivate warehouses",
			},
		}, http.StatusForbidden)
		return
	}
	_, err = controller.interactor.DeactivateWarehouse(DeactivateWarehouseID, session.User.Id)
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
		Message: "Warehouse deactivated successfully",
	}, http.StatusOK)
}

func (controller Controller) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Delete Warehouse Request ||||||||")
	controller.log.Println("Processing Delete Warehouse Request")
	authHeader := r.Header.Get("Authorization")
	if len(strings.Split(authHeader, " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in the header",
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

	requiredPermission := entity.Permission{
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
				Message: "You do not have permission to delete this warehouse.",
			},
		}, http.StatusForbidden)
		return
	}

	type DeleteWarehouseRequest struct {
		WarehouseID string `json:"warehouse_id"`
	}

	var req DeleteWarehouseRequest
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
	if req.WarehouseID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Warehouse ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}
	WarehouseID, err := uuid.Parse(req.WarehouseID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Warehouse ID format.",
			},
		}, http.StatusBadRequest)
		return
	}
	err = controller.interactor.DeleteWarehouse(WarehouseID, session.User.Id)
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
		Message: "Warehouse deleted successfully",
	}, http.StatusOK)
}
