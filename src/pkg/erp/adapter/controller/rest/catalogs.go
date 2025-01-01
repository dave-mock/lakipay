package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"auth/src/pkg/erp/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (controller Controller) CreateCatalog(w http.ResponseWriter, r *http.Request) {
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
	type CreateCatalogRequest struct {
		MerchantID  string `json:"merchant_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	var req CreateCatalogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid request body.",
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Status != "active" && req.Status != "archived" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Status must be 'active' or 'archived'.",
			},
		}, http.StatusBadRequest)
		return
	}

	catalog, err := controller.interactor.CreateCatalog(
		session.User.Id,
		req.Name,
		req.Description,
		req.Status,
	)
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
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Catalog created successfully",
		Data:    catalog,
	}, http.StatusCreated)
}

func (controller Controller) ListMerchantCatalogs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Merchant Catalogs Request ||||||||")
	controller.log.Println("Processing List Merchant Catalogs Request")
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
	merchantID := r.URL.Query().Get("merchant_id")
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
	catalogs, err := controller.interactor.ListMerchantCatalogs(merchantID, session.User.Id)
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
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Catalogs retrieved successfully",
		Data:    catalogs,
	}, http.StatusOK)
}

func (controller Controller) ListCatalogs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Catalogs Request ||||||||")
	controller.log.Println("Processing List Catalogs Request")
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
	catalogs, err := controller.interactor.ListCatalogs(session.User.Id)
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
	SendJSONResponse(w, Response{
		Success: true,
		Message: "All catalogs retrieved successfully",
		Data:    catalogs,
	}, http.StatusOK)
}

func (controller Controller) GetCatalog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Get Catalog Request ||||||||")
	controller.log.Println("Processing Get Catalog Request")

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

	var requestBody struct {
		CatalogID string `json:"catalogId"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.CatalogID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Catalog ID is required in the request body.",
			},
		}, http.StatusBadRequest)
		return
	}

	catalog, err := controller.interactor.GetCatalog(requestBody.CatalogID, session.User.Id)
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

	SendJSONResponse(w, Response{
		Success: true,
		Data:    catalog,
	}, http.StatusOK)
}

func (controller Controller) ArchiveCatalog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Archive Catalog Request ||||||||")
	controller.log.Println("Processing Archive Catalog Request")
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
	catalogID := r.URL.Query().Get("id")
	if catalogID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Catalog ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}
	err = controller.interactor.ArchiveCatalog(catalogID, session.User.Id)
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
		Message: "Catalog archived successfully",
	}, http.StatusOK)
}

func (controller Controller) UpdateCatalog(w http.ResponseWriter, r *http.Request) {
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

	type UpdateCatalogRequest struct {
		CatalogID   string `json:"catalog_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	var req UpdateCatalogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid request body.",
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Status != "active" && req.Status != "archived" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Status must be 'active' or 'archived'.",
			},
		}, http.StatusBadRequest)
		return
	}

	catalog, err := controller.interactor.UpdateCatalog(
		session.User.Id,
		req.CatalogID,
		req.Name,
		req.Description,
		req.Status,
	)
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

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Catalog updated successfully",
		Data:    catalog,
	}, http.StatusOK)
}

func (controller Controller) DeleteCatalog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Delete Catalog Request ||||||||")
	controller.log.Println("Processing Delete Catalog Request")
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

	type DeleteCatalogRequest struct {
		CatalogID string `json:"catalog_id"`
	}
	var req DeleteCatalogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid request body.",
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if req.CatalogID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Catalog ID is required in the request body.",
			},
		}, http.StatusBadRequest)
		return
	}

	if err := controller.interactor.DeleteCatalog(req.CatalogID, session.User.Id); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusInternalServerError)
		return
	}
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Catalog deleted successfully",
	}, http.StatusOK)
}
