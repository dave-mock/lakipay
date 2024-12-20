package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"auth/src/pkg/erp/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"auth/src/pkg/erp/core/entity"

	"github.com/google/uuid"
)

func (controller Controller) CreateSubCatalog(w http.ResponseWriter, r *http.Request) {
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

	// Validate the token
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

	type CreateSubCatalogRequest struct {
		MerchantID  string    `json:"merchant_id"`
		CatalogID   uuid.UUID `json:"catalog_id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
	}

	var req CreateSubCatalogRequest
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

	subCatalog := &entity.SubCatalog{
		ID:           uuid.New(),
		CatalogID:    req.CatalogID,
		SubCatalogID: uuid.New(),
		MerchantID:   session.User.Id,
		Name:         req.Name,
		Description:  req.Description,
		Status:       req.Status,
		CreatedBy:    session.User.Id,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Usecase [CREATE SUBCATALOG]
	err = controller.interactor.CreateSubCatalog(subCatalog)
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
		Message: "Subcatalog created successfully",
		Data:    subCatalog,
	}, http.StatusCreated)
}

func (controller Controller) ListSubCatalogs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List SubCatalogs Request ||||||||")
	controller.log.Println("Processing List SubCatalogs Request")
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

	// Usecase [LIST ALL SUBCATALOGS]
	subCatalogs, err := controller.interactor.ListSubCatalogs(session.User.Id)
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
		Message: "All subcatalogs retrieved successfully",
		Data:    subCatalogs,
	}, http.StatusOK)
}

func (controller Controller) GetSubCatalog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Get SubCatalog Request ||||||||")
	controller.log.Println("Processing Get SubCatalog Request")
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
		SubCatalogID string `json:"sub_catalog_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.SubCatalogID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "SubCatalog ID is required in the request body.",
			},
		}, http.StatusBadRequest)
		return
	}

	subCatalogID, err := uuid.Parse(requestBody.SubCatalogID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid SubCatalog ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Usecase [GET SUBCATALOG]
	subCatalog, err := controller.interactor.GetSubCatalog(subCatalogID, session.User.Id)
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
		Data:    subCatalog,
	}, http.StatusOK)
}

func (controller Controller) ArchiveSubCatalog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || Handle Archive SubCatalog Request ||||||||")
	controller.log.Println("Processing Archive SubCatalog Request")

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

	subCatalogID := r.URL.Query().Get("id")
	if subCatalogID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "SubCatalog ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	subCatalogUUID, err := uuid.Parse(subCatalogID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid SubCatalog ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Usecase [ARCHIVE SUBCATALOG]
	err = controller.interactor.ArchiveSubCatalog(subCatalogUUID, session.User.Id)
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
		Message: "Subcatalog archived successfully.",
	}, http.StatusOK)
}

func (controller Controller) UpdateSubCatalog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || Handle Update SubCatalog Request ||||||||")
	controller.log.Println("Processing Update SubCatalog Request")

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in the header",
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

	type UpdateSubCatalogRequest struct {
		MerchantID   string    `json:"merchant_id"`
		CatalogID    uuid.UUID `json:"catalog_id"`
		SubCatalogID uuid.UUID `json:"sub_catalog_id"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		Status       string    `json:"status"`
	}

	var req UpdateSubCatalogRequest
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

	if req.SubCatalogID == uuid.Nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "SubCatalog ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	if req.CatalogID == uuid.Nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Catalog ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

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

	subCatalog := &entity.SubCatalog{
		ID:          req.SubCatalogID,
		CatalogID:   req.CatalogID,
		MerchantID:  session.User.Id,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		UpdatedBy:   session.User.Id,
		UpdatedAt:   time.Now(),
	}

	updatedSubCatalog, err := controller.interactor.UpdateSubCatalog(session.User.Id, subCatalog)
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
		Message: "Subcatalog updated successfully",
		Data:    updatedSubCatalog,
	}, http.StatusOK)
}

func (controller Controller) DeleteSubCatalog(w http.ResponseWriter, r *http.Request) {
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
		SubCatalogID string `json:"sub_catalog_id"`
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

	if req.SubCatalogID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Catalog ID is required in the request body.",
			},
		}, http.StatusBadRequest)
		return
	}

	// Usecase [DELETE CATALOG]
	SubCatalogID, err := uuid.Parse(req.SubCatalogID)
	if err := controller.interactor.DeleteSubCatalog(SubCatalogID, session.User.Id); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusInternalServerError)
		return
	}

	// Send success response
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Catalog deleted successfully",
	}, http.StatusOK)
}
