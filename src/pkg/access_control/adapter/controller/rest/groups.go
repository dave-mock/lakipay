package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// -------------- Group Management --------------

func (controller Controller) CreateGroup(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [CreateGroup] ")
	fmt.Println("||||||| || Handle Create Group Request ||||||||")
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
	fmt.Print(session)
	type CreateGroupPayload struct {
		Title string `json:"title"`
	}

	var req CreateGroupPayload
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
	group, err := controller.interactor.CreateGroup(req.Title)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "GROUP_CREATION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    group,
	}, http.StatusCreated)
}

func (controller Controller) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [UpdateGroup] ")
	fmt.Println("||||||| || Handle Update Group Request ||||||||")
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
	fmt.Print(session)
	type UpdateGroupPayload struct {
		GroupID uuid.UUID `json:"group_id"`
		Title   string    `json:"title"`
	}

	var req UpdateGroupPayload
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

	group, err := controller.interactor.UpdateGroup(req.GroupID, req.Title)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "GROUP_UPDATE_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    group,
	}, http.StatusOK)
}

func (controller Controller) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [DeleteGroup] ")
	fmt.Println("||||||| || Handle Delete Group Request ||||||||")
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
	fmt.Print(session)
	type IDPayload struct {
		ID uuid.UUID `json:"id"`
	}

	var req IDPayload
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

	err = controller.interactor.DeleteGroup(req.ID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "GROUP_DELETION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "Group deleted successfully",
	}, http.StatusOK)
}

func (controller Controller) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [ListGroups] ")
	fmt.Println("||||||| || Handle Get all Groups Request ||||||||")

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
	groups, err := controller.interactor.ListGroups()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "GROUP_LIST_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    groups,
	}, http.StatusOK)
}

func (controller Controller) GetSingleGroup(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[CONTROLLER] [GetSingleGroup] ")
	fmt.Println("||||||| || Handle Get Single Group Request ||||||||")

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
	groupIDParam := r.URL.Query().Get("group_id")
	groupID, err := uuid.Parse(groupIDParam)
	if err != nil {
		controller.log.Println("Invalid group ID:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_GROUP_ID",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	group, err := controller.interactor.ListGroupPermissions(groupID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "GROUP_NOT_FOUND",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    group,
	}, http.StatusOK)
}
