package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// -------------- User-Group Management --------------

func (controller Controller) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [AddUserToGroup] ")
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
	type AddUserToGroupPayload struct {
		UserID  string `json:"user_id"`
		GroupID string `json:"group_id"`
	}

	var req AddUserToGroupPayload
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
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

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		controller.log.Println("Invalid user ID format:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_USER_ID",
				Message: "User ID must be a valid UUID",
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
				Message: "Group ID must be a valid UUID",
			},
		}, http.StatusBadRequest)
		return
	}

	err = controller.interactor.AddUserToGroup(userID, groupID)
	if err != nil {
		controller.log.Println("Error adding user to group:", err)
		if strings.Contains(err.Error(), "does not exist") {
			if strings.Contains(err.Error(), "user") {
				SendJSONResponse(w, Response{
					Success: false,
					Error: &Error{
						Type:    "USER_NOT_FOUND",
						Message: "User not found",
					},
				}, http.StatusNotFound)
			} else if strings.Contains(err.Error(), "group") {
				SendJSONResponse(w, Response{
					Success: false,
					Error: &Error{
						Type:    "GROUP_NOT_FOUND",
						Message: "Group not found",
					},
				}, http.StatusNotFound)
			}
			return
		} else if strings.Contains(err.Error(), "already in group") {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "USER_ALREADY_IN_GROUP",
					Message: "User is already in this group",
				},
			}, http.StatusConflict)
			return
		}

		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "USER_GROUP_ADD_ERROR",
				Message: "Failed to add user to group",
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "User added to group successfully",
	}, http.StatusOK)
}

func (controller Controller) RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [RemoveUserFromGroup] ")
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
	type RemoveUserFromGroupPayload struct {
		UserID  string `json:"user_id"`
		GroupID string `json:"group_id"`
	}

	var req RemoveUserFromGroupPayload
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		controller.log.Println("Error decoding request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Failed to parse request body: " + err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		controller.log.Println("Invalid user ID format:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_USER_ID",
				Message: "User ID must be a valid UUID. " + err.Error(),
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
				Message: "Group ID must be a valid UUID. " + err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	err = controller.interactor.RemoveUserFromGroup(userID, groupID)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with ID %v does not exist", userID) {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "USER_NOT_FOUND",
					Message: fmt.Sprintf("User with ID %v does not exist", userID),
				},
			}, http.StatusNotFound)
			return
		}

		if err.Error() == fmt.Sprintf("group with ID %v does not exist", groupID) {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "GROUP_NOT_FOUND",
					Message: fmt.Sprintf("Group with ID %v does not exist", groupID),
				},
			}, http.StatusNotFound)
			return
		}

		if err.Error() == fmt.Sprintf("user with ID %v is not in group with ID %v", userID, groupID) {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "USER_NOT_IN_GROUP",
					Message: fmt.Sprintf("User with ID %v is not in group with ID %v", userID, groupID),
				},
			}, http.StatusBadRequest)
			return
		}

		controller.log.Println("Error removing user from group:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "USER_GROUP_REMOVE_ERROR",
				Message: "Failed to remove " + err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "User removed from group successfully",
	}, http.StatusOK)
}

func (controller Controller) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [GetUserGroups] ")
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
	type GetUserGroupsPayload struct {
		UserID string `json:"user_id"`
	}

	var req GetUserGroupsPayload
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		controller.log.Println("Error decoding request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Failed to parse request body. " + err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		controller.log.Println("Invalid user ID format:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_USER_ID",
				Message: "User ID must be a valid UUID. " + err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	groups, err := controller.interactor.ListUserGroups(userID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "USER_GROUP_NOT_FOUND",
					Message: fmt.Sprintf("User with ID %v is not a member of any group", userID),
				},
			}, http.StatusNotFound)
			return
		}
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "USER_GROUP_FETCH_ERROR",
				Message: "Failed to fetch user groups. " + err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	// If everything is okay, return the response with the groups
	SendJSONResponse(w, Response{
		Success: true,
		Data:    groups,
	}, http.StatusOK)
}
