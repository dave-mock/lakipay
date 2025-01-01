package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// -------------- Resource Management --------------
func (controller Controller) CreateResource(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [CreateResource] ")

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

	type CreateResourcePayload struct {
		Name        string      `json:"name"`
		Operations  []uuid.UUID `json:"operations"`
		Description string      `json:"description"`
	}

	var req CreateResourcePayload
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

	resource, err := controller.interactor.CreateResource(req.Name, req.Description, req.Operations)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "RESOURCE_CREATION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    resource,
	}, http.StatusCreated)
}

func (controller Controller) UpdateResource(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [UpdateResource] ")

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

	type UpdateResourcePayload struct {
		ResourceID  uuid.UUID   `json:"resource_id"`
		Name        string      `json:"name"`
		Operations  []uuid.UUID `json:"operations"`
		Description string      `json:"description"`
	}

	var req UpdateResourcePayload
	defer r.Body.Close()
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

	if req.ResourceID == uuid.Nil {
		controller.log.Println("Error: Invalid or missing resource_id")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "resource_id is required and must be valid",
			},
		}, http.StatusBadRequest)
		return
	}
	controller.log.Printf("Received UpdateResource request with ResourceID %s", req.ResourceID)
	resource, err := controller.interactor.UpdateResource(req.ResourceID, req.Name, req.Description, req.Operations)
	if err != nil {
		controller.log.Printf("Error updating resource: %v", err)
		if strings.Contains(err.Error(), "does not exist") {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "RESOURCE_NOT_FOUND",
					Message: fmt.Sprintf("Resource with ID %s does not exist", req.ResourceID),
				},
			}, http.StatusNotFound)
			return
		}

		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "RESOURCE_UPDATE_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    resource,
	}, http.StatusOK)
}

func (controller Controller) DeleteResource(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [DeleteResource] ")

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

	type DeleteResourcePayload struct {
		ResourceID uuid.UUID `json:"resource_id"`
	}

	var req DeleteResourcePayload
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

	err = controller.interactor.DeleteResource(req.ResourceID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "RESOURCE_DELETE_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "Resource deleted successfully",
	}, http.StatusOK)
}

func (controller Controller) GetSingleResource(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [GetSingleResource] ")

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

	type GetResourcePayload struct {
		ResourceID uuid.UUID `json:"resource_id"`
	}

	var req GetResourcePayload
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

	resource, err := controller.interactor.GetResourceByID(req.ResourceID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "RESOURCE_NOT_FOUND",
				Message: err.Error(),
			},
		}, http.StatusNotFound)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    resource,
	}, http.StatusOK)
}

func (controller Controller) ListAllResources(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [ListAllResources] ")

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

	resources, err := controller.interactor.ListResources()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "RESOURCE_LIST_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    resources,
	}, http.StatusOK)
}
