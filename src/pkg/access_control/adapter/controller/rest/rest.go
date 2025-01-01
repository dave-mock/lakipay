package rest

import (
	"auth/src/pkg/access_control/usecase"
	auth "auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

type CreateGroupPayload struct {
	Title string `json:"title"`
}

type UpdateGroupPayload struct {
	GroupID string `json:"group_id"`
	Title   string `json:"title"`
}

type IDPayload struct {
	ID string `json:"id"`
}

type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	auth       auth.Controller
	sm         *http.ServeMux
}

func New(log *log.Logger, interactor usecase.Interactor, sm *http.ServeMux, auth auth.Controller) Controller {
	controller := Controller{log: log, interactor: interactor, auth: auth, sm: sm}
	// Merchant Management
	/* 	sm.HandleFunc("/merchants", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateMerchant(w, r)
		}
	}) */
	/* 	sm.HandleFunc("/merchants/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListMerchants(w, r)
		case http.MethodPut:
			controller.UpdateMerchant(w, r)
		case http.MethodDelete:
			controller.DeactivateMerchant(w, r)
		}
	}) */
	sm.HandleFunc("/user/groups/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateGroup(w, r)
		}
	})
	sm.HandleFunc("/user/groups/update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.UpdateGroup(w, r)
		}
	})
	sm.HandleFunc("/user/groups/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.DeleteGroup(w, r)
		}
	})
	sm.HandleFunc("/user/groups/list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetAllGroups(w, r)
		}
	})
	sm.HandleFunc("/user/groups/view", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetSingleGroup(w, r)
		}
	})

	sm.HandleFunc("/user/permissions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreatePermission(w, r)
		}
	})
	sm.HandleFunc("/user/permissions/update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.UpdatePermission(w, r)
		}
	})
	sm.HandleFunc("/user/permissions/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.DeletePermission(w, r)
		}
	})
	sm.HandleFunc("/user/permissions/list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListPermissions(w, r)
		}
	})
	sm.HandleFunc("/user/permissions/view", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetSinglePermission(w, r)
		}
	})

	// User-Group Assignments
	sm.HandleFunc("/user/user_groups/assign", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.AddUserToGroup(w, r)
		}
	})
	sm.HandleFunc("/user/user-groups/remove", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.RemoveUserFromGroup(w, r)
		}
	})
	sm.HandleFunc("/user/user-groups/view", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetUserGroups(w, r)
		}
	})
	sm.HandleFunc("/groups/users/list-users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetGroupUsers(w, r)
		}
	})

	// User-Permission Assignments
	sm.HandleFunc("/user/user-permissions/assign", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GrantPermissionToUser(w, r)
		}
	})
	sm.HandleFunc("/user/user-permissions/revoke", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.RevokePermissionFromUser(w, r)
		}
	})
	sm.HandleFunc("/user/use-permissions/view", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetUserPermissions(w, r)
		}
	})

	// User-Permission Assignments
	sm.HandleFunc("/group/group-permissions/assign", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GrantPermissionToGroup(w, r)
		}
	})
	sm.HandleFunc("/group/group-permissions/revoke", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.RevokePermissionFromGroup(w, r)
		}
	})
	sm.HandleFunc("/group/group-permissions/view", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetGroupPermissions(w, r)
		}
	})

	sm.HandleFunc("/user/resources/add", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateResource(w, r)
		}
	})
	sm.HandleFunc("/user/resources/update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.UpdateResource(w, r)
		}
	})
	sm.HandleFunc("/user/resources/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.DeleteResource(w, r)
		}
	})
	sm.HandleFunc("/user/resources/list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListAllResources(w, r)
		}
	})

	sm.HandleFunc("/user/operations/add", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateOperations(w, r)
		}
	})
	sm.HandleFunc("/user/operations/update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.UpdateOperations(w, r)
		}
	})
	sm.HandleFunc("/user/operations/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.DeleteOperations(w, r)
		}
	})
	sm.HandleFunc("/user/operations/list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListAllOperations(w, r)
		}
	})

	sm.HandleFunc("/user/resources/view", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetSingleResource(w, r)
		}
	})
	sm.HandleFunc("/account/users/list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetUsersList(w, r)
		}
	})
	return controller
}

func SendJSONResponse(w http.ResponseWriter, data Response, status int) {
	serData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(serData)
}
