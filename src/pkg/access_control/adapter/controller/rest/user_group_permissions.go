package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// -------------- Permission Management --------------

func (controller Controller) GrantPermissionToGroup(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [GrantPermissionToUser] ")
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

	type GrantPermissionPayload struct {
		GroupID      string `json:"group_id"`
		PermissionID string `json:"permission_id"`
	}

	var req GrantPermissionPayload
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		controller.log.Println("Error decoding request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Failed to decode request payload",
			},
		}, http.StatusBadRequest)
		return
	}

	groupID, err := uuid.Parse(req.GroupID)
	if err != nil || groupID == uuid.Nil {
		controller.log.Println("Invalid user ID format:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_GROUP_ID",
				Message: "Invalid or empty group ID",
			},
		}, http.StatusBadRequest)
		return
	}

	permissionID, err := uuid.Parse(req.PermissionID)
	if err != nil || permissionID == uuid.Nil {
		controller.log.Println("Invalid permission ID format:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_PERMISSION_ID",
				Message: "Invalid or empty permission ID",
			},
		}, http.StatusBadRequest)
		return
	}

	err = controller.interactor.GrantPermissionToGroup(groupID, permissionID)
	if err != nil {
		controller.log.Println("Failed to grant permission:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "GRANT_PERMISSION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "Permission granted to user successfully",
	}, http.StatusOK)
}

func (controller Controller) RevokePermissionFromGroup(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [RevokePermissionFromUser] ")
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
	type RevokePermissionPayload struct {
		GroupID      string `json:"group_id"`
		PermissionID string `json:"permission_id"`
	}

	var req RevokePermissionPayload
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		controller.log.Println("Error decoding request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Failed to decode request payload",
			},
		}, http.StatusBadRequest)
		return
	}

	groupID, err := uuid.Parse(req.GroupID)
	if err != nil || groupID == uuid.Nil {
		controller.log.Println("Invalid user ID format:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_GROUP_ID",
				Message: "Invalid or empty group ID",
			},
		}, http.StatusBadRequest)
		return
	}

	permissionID, err := uuid.Parse(req.PermissionID)
	if err != nil || permissionID == uuid.Nil {
		controller.log.Println("Invalid permission ID format:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_PERMISSION_ID",
				Message: "Invalid or empty permission ID",
			},
		}, http.StatusBadRequest)
		return
	}

	err = controller.interactor.RevokePermissionFromGroup(groupID, permissionID)
	if err != nil {
		controller.log.Println("Failed to revoke permission:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "REVOKE_PERMISSION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "Permission revoked from user successfully",
	}, http.StatusOK)
}

func (controller Controller) GetGroupPermissions(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [GetGroupPermissions] ")
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
	type GetGroupPermissionsPayload struct {
		GroupID string `json:"group_id"`
	}

	var req GetGroupPermissionsPayload
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		controller.log.Println("Error decoding request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Failed to parse request payload. Ensure the request body is valid JSON.",
			},
		}, http.StatusBadRequest)
		return
	}

	groupID, err := uuid.Parse(req.GroupID)
	if err != nil {
		controller.log.Println("Invalid group ID format:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_GROUP_ID",
				Message: "The provided group ID is invalid. Please provide a valid UUID.",
			},
		}, http.StatusBadRequest)
		return
	}

	permissions, err := controller.interactor.ListGroupPermissions(groupID)
	if err != nil {
		if err.Error() == "no permissions found for the provided group ID" {
			controller.log.Println("No permissions found for group ID", groupID)
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "NO_PERMISSIONS_FOUND",
					Message: "No permissions found for the specified group.",
				},
			}, http.StatusNotFound)
			return
		}

		controller.log.Println("Error fetching permissions:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "USER_PERMISSION_FETCH_ERROR",
				Message: "An error occurred while retrieving permissions.",
			},
		}, http.StatusInternalServerError)
		return
	}

	controller.log.Println("Permissions retrieved successfully for group ID:", groupID)
	SendJSONResponse(w, Response{
		Success: true,
		Data:    permissions,
	}, http.StatusOK)
}
