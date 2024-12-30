package rest

import (
	"auth/src/pkg/auth/adapter/controller/procedure"
	"auth/src/pkg/auth/core/entity"
	"auth/src/pkg/erp/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (controller Controller) CreateProduct(w http.ResponseWriter, r *http.Request) {
	controller.log.Println("Processing Create Product Request")
	if r.Method != http.MethodPost {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "METHOD_NOT_ALLOWED",
				Message: "Invalid HTTP method, expected POST",
			},
		}, http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	token := tokenParts[1]
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
		Resource:           "products",
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
				Message: "You do not have permission to create a product.",
			},
		}, http.StatusForbidden)
		return
	}

	type CreateProductRequest struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description,omitempty"`
		Price       float64 `json:"price" binding:"required"`
		Currency    string  `json:"currency" binding:"required"`
		SKU         string  `json:"sku" binding:"required"`
		Weight      float64 `json:"weight,omitempty"`
		Dimensions  string  `json:"dimensions,omitempty"`
		ImageURL    string  `json:"image_url,omitempty"`
	}

	var req CreateProductRequest
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

	if req.Price <= 0 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Price must be greater than zero.",
			},
		}, http.StatusBadRequest)
		return
	}

	err = controller.interactor.CreateProduct(
		req.Name,
		req.Description,
		req.Price,
		req.Currency,
		req.SKU,
		req.Weight,
		req.Dimensions,
		req.ImageURL,
		session.User.Id,
	)

	if err != nil {
		if strings.Contains(err.Error(), "DUPLICATE_SKU") {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "DUPLICATE_SKU",
					Message: fmt.Sprintf("A product with the SKU '%s' already exists.", req.SKU),
				},
			}, http.StatusConflict)
			return
		}
		controller.log.Println("Error creating product:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "An error occurred while creating the product.",
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Product created successfully.",
	}, http.StatusCreated)
}
func (controller Controller) ListProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle List Products Request ||||||||")
	controller.log.Println("Processing List Products Request")
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
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}
	requiredPermission := entity.Permission{
		Resource:           "Users",                 // Resource to access
		Operation:          "create",                // Operation to perform (create)
		ResourceIdentifier: "/users/profile/update", // Specific resource identifier
		Effect:             "allow",                 // Effect ("allow" or "deny")
	}

	hasPermission, err := controller.auth.HasPermission(session.User.Id, requiredPermission)
	if err != nil || !hasPermission {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "FORBIDDEN",
				Message: "You do not have permission to read users' profile information.",
			},
		}, http.StatusForbidden)
		return
	}

	userType := session.User.UserType
	products, err := controller.interactor.ListProducts(session.User.Id, userType)
	if err != nil {
		switch err := err.(type) {
		case usecase.Error:
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    err.Type,
					Message: err.Message,
				},
			}, http.StatusInternalServerError)
		default:
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "INTERNAL_SERVER_ERROR",
					Message: err.Error(),
				},
			}, http.StatusInternalServerError)
		}
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    products,
	}, http.StatusOK)
}

func (controller Controller) GetProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Get Product Request ||||||||")
	controller.log.Println("Processing Get Product Request")

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
		Resource:           "Products",
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
				Message: "You do not have permission to read products.",
			},
		}, http.StatusForbidden)
		return
	}

	productIDStr := r.URL.Query().Get("id")
	if productIDStr == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Product ID is required.",
			},
		}, http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Product ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	product, err := controller.interactor.GetProduct(productID, session.User.Id)
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
		Message: "Product retrieved successfully",
		Data:    product,
	}, http.StatusOK)
}

func (controller Controller) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Update Product Request ||||||||")
	controller.log.Println("Processing Update Product Request")
	authHeader := r.Header.Get("Authorization")
	if len(strings.Split(authHeader, " ")) != 2 {
		controller.log.Println("Authorization token missing or malformed")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	token := strings.Split(authHeader, " ")[1]
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		controller.log.Println("Failed to authorize token:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	/* requiredPermission := entity.Permission{
		Resource:           "products",
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
				Message: "You do not have permission to update products.",
			},
		}, http.StatusForbidden)
		return
	} */
	type UpdateProductRequest struct {
		ProductID   string  `json:"product_id"`
		Name        string  `json:"name,omitempty"`
		Description string  `json:"description,omitempty"`
		Price       float64 `json:"price,omitempty"`
		Currency    string  `json:"currency,omitempty"`
		SKU         string  `json:"sku,omitempty"`
		Weight      float64 `json:"weight,omitempty"`
		Dimensions  string  `json:"dimensions,omitempty"`
		ImageURL    string  `json:"image_url,omitempty"`
		Status      string  `json:"status,omitempty"`
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Println("Invalid request body:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid request body.",
			},
		}, http.StatusBadRequest)
		return
	}

	if req.ProductID == "" {
		controller.log.Println("Product ID missing in the request body")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Product ID is required in the request body.",
			},
		}, http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	err = controller.interactor.UpdateProduct(productID, req.Name, req.Description, req.Price, req.Currency, req.SKU, req.Weight, req.Dimensions, req.ImageURL, req.Status, session.User.Id)
	if err != nil {
		controller.log.Println("Failed to update product:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	controller.log.Println("Product updated successfully")
	SendJSONResponse(w, Response{
		Success: true,
		Message: "Product updated successfully",
	}, http.StatusOK)
}

func (controller Controller) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || Handle Deactivate Product Request ||||||||")
	controller.log.Println("Processing Deactivate Product Request")

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
	/*
		requiredPermission := entity.Permission{
			Resource:           "products",
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
					Message: "You do not have permission to delete products.",
				},
			}, http.StatusForbidden)
			return
		} */
	var requestBody struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Unable to read request body or missing ID",
			},
		}, http.StatusBadRequest)
		return
	}

	UpdateProductID, err := uuid.Parse(requestBody.ID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Product ID format.",
			},
		}, http.StatusBadRequest)
		return
	}
	err = controller.interactor.DeleteProduct(UpdateProductID, session.User.Id)
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
		Message: "Product deactivated successfully",
	}, http.StatusOK)
}

func (controller Controller) CountMerchantProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Count Merchant Orders Request ||||||||")
	controller.log.Println("Processing Count Merchant Orders Request")
	authorizationHeader := r.Header.Get("Authorization")
	if len(strings.Split(authorizationHeader, " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in the header",
			},
		}, http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorizationHeader, " ")[1]
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

	/* 	requiredPermission := entity.Permission{
	   		Resource:           "products",
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
	   				Message: "You do not have permission to count products.",
	   			},
	   		}, http.StatusForbidden)
	   		return
	   	} */
	totalProductsCount, err := controller.interactor.CountMerchantProducts(session.User.Id)
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
		Message: "Total Product retrieved successfully",
		Data: map[string]int{
			"total_products": totalProductsCount,
		},
	}, http.StatusOK)
}
func (controller Controller) AddProductToCatalog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Add Product to Catalog Request ||||||||")
	controller.log.Println("Processing Add Product to Catalog Request")

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

	/* requiredPermission := entity.Permission{
		Resource:           "products",
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
				Message: "You do not have permission to add product to catalogs.",
			},
		}, http.StatusForbidden)
		return
	} */
	type AddProductToCatalogRequest struct {
		ProductID    string `json:"product_id"`
		CatalogID    string `json:"catalog_id"`
		DisplayOrder int    `json:"display_order,omitempty"`
	}

	var requestPayload AddProductToCatalogRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Unable to read request body.",
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &requestPayload); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid JSON format.",
			},
		}, http.StatusBadRequest)
		return
	}

	if requestPayload.ProductID == "" || requestPayload.CatalogID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Product ID and Catalog ID are required.",
			},
		}, http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(requestPayload.ProductID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Product ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	catalogID, err := uuid.Parse(requestPayload.CatalogID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Catalog ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	_, err = controller.interactor.AddProductToCatalog(
		productID,
		catalogID,
		requestPayload.DisplayOrder,
		session.User.Id,
	)

	if err != nil {
		if errors.Is(err, usecase.ErrProductDoesNotExist) {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "PRODUCT_NOT_FOUND",
					Message: "Product with the given ID does not exist.",
				},
			}, http.StatusNotFound)
		} else {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "SERVER_ERROR",
					Message: err.Error(),
				},
			}, http.StatusInternalServerError)
		}
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Message: "Product added to catalog successfully",
	}, http.StatusOK)
}

func (controller Controller) RemoveProductFromCatalog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || handle Remove Product from Catalog Request ||||||||")
	controller.log.Println("Processing Remove Product from Catalog Request")
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
	/*
		requiredPermission := entity.Permission{
			Resource:           "products",
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
					Message: "You do not have permission to remove products from catalogs.",
				},
			}, http.StatusForbidden)
			return
		} */

	type RemoveProductFromCatalogRequest struct {
		ProductID string `json:"product_id"`
		CatalogID string `json:"catalog_id"`
	}

	var requestPayload RemoveProductFromCatalogRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Unable to read request body.",
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, &requestPayload); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid JSON format.",
			},
		}, http.StatusBadRequest)
		return
	}

	if requestPayload.ProductID == "" || requestPayload.CatalogID == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Product ID and Catalog ID are required.",
			},
		}, http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(requestPayload.ProductID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Product ID format.",
			},
		}, http.StatusBadRequest)
		return
	}

	catalogID, err := uuid.Parse(requestPayload.CatalogID)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Invalid Catalog ID format.",
			},
		}, http.StatusBadRequest)
		return
	}
	_, err = controller.interactor.RemoveProductFromCatalog(
		productID,
		catalogID,
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
		Message: "Product removed from catalog successfully",
	}, http.StatusOK)
}
