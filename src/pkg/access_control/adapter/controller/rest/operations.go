package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// -------------- Operations Management --------------

func (controller Controller) CreateOperations(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [CreateOperations] ")

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
	fmt.Print(session)

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		controller.log.Println("Error decoding request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	operations, err := controller.interactor.CreateOperations(req.Name, req.Description)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "Operations_CREATION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    operations,
	}, http.StatusCreated)
}

func (controller Controller) UpdateOperations(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [UpdateOperations] ")

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
	fmt.Print(session)

	var req struct {
		OperationID string `json:"operation_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		controller.log.Println("Error reading request body:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Error reading request body",
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Println("Raw Body:", string(body))

	r.Body = io.NopCloser(bytes.NewReader(body))
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		controller.log.Println("Error decoding request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid JSON payload",
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Println("Received OperationID", req.OperationID)
	operationsUUID, err := uuid.Parse(req.OperationID)
	if err != nil {
		controller.log.Println("Error parsing operation_id:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid operation_id format",
			},
		}, http.StatusBadRequest)
		return
	}

	operations, err := controller.interactor.UpdateOperations(operationsUUID, req.Name, req.Description)
	if err != nil {
		controller.log.Printf("Error updating operations: %v", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "Operations_UPDATE_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    operations,
	}, http.StatusOK)
}

func (controller Controller) DeleteOperations(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [DeleteOperations] ")

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
	fmt.Print(session)

	var req struct {
		OperationsID string `json:"operation_id"`
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		controller.log.Println("Error decoding request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Println("Received OperationsID:", req.OperationsID)
	if req.OperationsID == "" {
		controller.log.Println("Error: OperationsID is empty")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "OperationsID is empty",
			},
		}, http.StatusBadRequest)
		return
	}

	operationsUUID, err := uuid.Parse(req.OperationsID)
	if err != nil {
		controller.log.Println("Error parsing operations_id:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid operations_id format",
			},
		}, http.StatusBadRequest)
		return
	}

	err = controller.interactor.DeleteOperations(operationsUUID)
	if err != nil {
		controller.log.Printf("Error deleting operation: %v", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "Operations_DELETE_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "Operation deleted successfully",
	}, http.StatusOK)
}

func (controller Controller) GetSingleOperations(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [GetSingleOperations] ")

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
	fmt.Print(session)

	var req struct {
		OperationsID uuid.UUID `json:"operations_id"`
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		controller.log.Println("Error decoding request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	operations, err := controller.interactor.GetOperationsByID(req.OperationsID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "Operations_NOT_FOUND",
				Message: err.Error(),
			},
		}, http.StatusNotFound)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    operations,
	}, http.StatusOK)
}

func (controller Controller) ListAllOperations(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [ListAllOperations] ")

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
	fmt.Print(session)

	operations, err := controller.interactor.ListOperations()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "Operations_LIST_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    operations,
	}, http.StatusOK)
}
