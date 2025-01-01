package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// -------------- Permission Management --------------

func (controller Controller) CreatePermission(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [CreatePermission] ")
	authHeader := r.Header.Get("Authorization")
	if len(strings.Split(authHeader, " ")) != 2 {
		controller.log.Println("Authentication token missing in header.")
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
		controller.log.Printf("Authentication failed: %v\n", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	controller.log.Printf("||||||||||| Auth session %+v\n", session)

	var req struct {
		ResourceID string      `json:"resource_id"`
		Resource   uuid.UUID   `json:"resource"`
		Operation  []uuid.UUID `json:"operations"`
		Effect     string      `json:"effect"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		controller.log.Printf("Error reading body %v\n", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Error reading body.",
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Printf("Raw body: %s\n", string(body))
	err = json.Unmarshal(body, &req)
	if err != nil {
		controller.log.Printf("Error decoding request body %v\n", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Printf("Decoded request %+v\n", req)
	if len(req.Operation) == 0 {
		controller.log.Println("Operations array is empty")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Operations cannot be empty.",
			},
		}, http.StatusBadRequest)
		return
	}

	permission, err := controller.interactor.CreatePermission(req.ResourceID, req.Resource, req.Operation, req.Effect)
	if err != nil {
		controller.log.Printf("Error creating permission: %v\n", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "PERMISSION_CREATION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	controller.log.Printf("Created permission %+v\n", permission)

	SendJSONResponse(w, Response{
		Success: true,
		Data:    permission,
	}, http.StatusCreated)
}

func (controller Controller) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [UpdatePermission] ")

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

	fmt.Println(session)

	var req struct {
		PermissionID uuid.UUID   `json:"permission_id"`
		ResourceID   string      `json:"resource_id"`
		Resource     uuid.UUID   `json:"resource"`
		Operations   []uuid.UUID `json:"operations"`
		Effect       string      `json:"effect"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		controller.log.Printf("Error reading body %v\n", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Error reading body.",
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Printf("Raw body: %s\n", string(body))
	err = json.Unmarshal(body, &req)
	if err != nil {
		controller.log.Printf("Error decoding request body %v\n", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Printf("Decoded request %+v\n", req)
	if len(req.Operations) == 0 {
		controller.log.Println("Operations array is empty")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Operations cannot be empty.",
			},
		}, http.StatusBadRequest)
		return
	}

	permission, err := controller.interactor.UpdatePermission(req.PermissionID, req.ResourceID, req.Resource, req.Operations, req.Effect)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "PERMISSION_UPDATE_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    permission,
	}, http.StatusOK)
}

func (controller Controller) GetSinglePermission(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [GetSinglePermission] ")

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

	fmt.Println(session)

	permissions, err := controller.interactor.ListPermissions()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "PERMISSION_LIST_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    permissions,
	}, http.StatusOK)
}

func (controller Controller) ListPermissions(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [ListPermissions] ")

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

	fmt.Println(session)

	permissions, err := controller.interactor.ListPermissions()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "PERMISSION_LIST_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    permissions,
	}, http.StatusOK)
}

func (controller Controller) DeletePermission(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [DeletePermission] ")

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

	fmt.Println(session)

	var req struct {
		ResourceID   string `json:"resource_id"`
		PermissionID string `json:"permission_id"`
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Println("Error decoding request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Failed to parse request body",
			},
		}, http.StatusBadRequest)
		return
	}

	resourceID, err := uuid.Parse(req.ResourceID)
	if err != nil {
		controller.log.Println("Invalid resource ID format:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_RESOURCE_ID",
				Message: "Resource ID must be a valid UUID",
			},
		}, http.StatusBadRequest)
		return
	}

	permissionID, err := uuid.Parse(req.PermissionID)
	if err != nil {
		controller.log.Println("Invalid permission ID format:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_PERMISSION_ID",
				Message: "Permission ID must be a valid UUID",
			},
		}, http.StatusBadRequest)
		return
	}

	err = controller.interactor.DeletePermission(resourceID, permissionID)
	if err != nil {
		controller.log.Println("Error deleting permission:", err)
		var statusCode int
		errorType := "PERMISSION_DELETE_ERROR"
		switch {
		case strings.Contains(err.Error(), "resource"):
			statusCode = http.StatusNotFound
			errorType = "RESOURCE_NOT_FOUND"
		case strings.Contains(err.Error(), "permission"):
			statusCode = http.StatusNotFound
			errorType = "PERMISSION_NOT_FOUND"
		default:
			statusCode = http.StatusInternalServerError
		}
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    errorType,
				Message: err.Error(),
			},
		}, statusCode)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "Permission deleted successfully",
	}, http.StatusOK)
}
