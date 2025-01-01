package usecase

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"auth/src/pkg/erp/core/entity"

	"github.com/google/uuid"
)

type JError struct {
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
}

func (e *JError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorType, e.Message)
}

const (
	ErrFailedToCreateProduct            = "FAILED_TO_CREATE_PRODUCT"
	ErrFailedToUpdateProduct            = "FAILED_TO_UPDATE_PRODUCT"
	ErrFailedToDeleteProduct            = "FAILED_TO_DELETE_PRODUCT"
	ErrFailedToListProducts             = "FAILED_TO_LIST_PRODUCTS"
	ErrFailedToGetProduct               = "FAILED_TO_GET_PRODUCT"
	ErrFailedToArchiveProduct           = "FAILED_TO_ARCHIVE_PRODUCT"
	ErrFailedToAddProductToCatalog      = "FAILED_TO_ADD_PRODUCT_TO_CATALOG"
	ErrFailedToRemoveProductFromCatalog = "FAILED_TO_REMOVE_PRODUCT_FROM_CATALOG"
)

var ErrProductDoesNotExist = errors.New("product does not exist")

func (uc Usecase) CreateProduct(name string, description string, price float64, currency string, SKU string, weight float64, dimensions string, imageURL string, userId uuid.UUID) error {
	if err := validateProductInput(
		name,
		description,
		price,
		currency,
		SKU,
		weight,
		dimensions,
		imageURL); err != nil {
		uc.log.Println("ERROR VALIDATING PRODUCT INPUT:", err)
		return fmt.Errorf("%s: %w", ErrFailedToCreateProduct, err)
	}

	uc.log.Println("CREATING PRODUCT")
	err := uc.repo.CreateProduct(
		name,
		description,
		price,
		currency,
		SKU,
		weight,
		dimensions,
		imageURL,
		userId)
	if err != nil {
		uc.log.Println("ERROR CREATING PRODUCT:", err)
		return fmt.Errorf("%s: %w", ErrFailedToCreateProduct, err)
	}

	uc.log.Println("PRODUCT CREATED SUCCESSFULLY")
	return nil
}

func (uc Usecase) GetProduct(productId, userId uuid.UUID) (*entity.Product, error) {
	uc.log.Println("FETCHING PRODUCT")
	product, err := uc.repo.GetProduct(productId, userId)
	if err != nil {
		uc.log.Println("ERROR FETCHING PRODUCT:", err)
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	uc.log.Println("Product fetched successfully")
	return product, nil
}

func (uc Usecase) UpdateProduct(
	productIDParsed uuid.UUID,
	name, description string,
	price float64,
	currency string,
	SKU string,
	weight float64,
	dimensions string,
	imageURL, status string,
	userId uuid.UUID,
) error {
	product, err := uc.repo.GetProduct(productIDParsed, userId)
	if err != nil {
		uc.log.Println("ERROR FETCHING PRODUCT FOR UPDATE:", err)
		return fmt.Errorf("%s: %w", ErrFailedToUpdateProduct, err)
	}

	if product == nil {
		uc.log.Println("PRODUCT NOT FOUND FOR UPDATE")
		return fmt.Errorf("product not found")
	}

	uc.log.Println("UPDATING PRODUCT")
	err = uc.repo.UpdateProduct(
		productIDParsed,
		name,
		description,
		price,
		currency,
		SKU,
		weight,
		dimensions,
		imageURL,
		status,
		userId)

	// Handle  error for duplicate SKU
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			uc.log.Println("ERROR UPDATING PRODUCT: duplicate SKU")
			return fmt.Errorf("duplicate SKU: product with this SKU already exists")
		}
		uc.log.Println("ERROR UPDATING PRODUCT:", err)
		return fmt.Errorf("%s: %w", ErrFailedToUpdateProduct, err)
	}

	uc.log.Println("PRODUCT UPDATED SUCCESSFULLY")
	return nil
}

func (uc Usecase) DeleteProduct(productID uuid.UUID, userId uuid.UUID) error {
	uc.log.Println("DELETING PRODUCT")
	err := uc.repo.DeleteProduct(productID, userId)
	if err != nil {
		uc.log.Println("ERROR DELETING PRODUCT:", err)
		return fmt.Errorf("%s: %w", ErrFailedToDeleteProduct, err)
	}

	uc.log.Println("PRODUCT DELETED SUCCESSFULLY")
	return nil
}

func (uc Usecase) ListProducts(userId uuid.UUID, userType string) ([]entity.Product, error) {
	uc.log.Println("FETCHING PRODUCT LIST")
	products, err := uc.repo.ListProducts(userId, userType)
	if err != nil {
		uc.log.Println("ERROR FETCHING PRODUCT LIST:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToListProducts, err)
	}

	uc.log.Println("PRODUCT LIST FETCHED SUCCESSFULLY")
	return products, nil
}

func (uc Usecase) ListAllProducts() ([]entity.Product, error) {
	uc.log.Println("FETCHING ALL PRODUCTS")
	products, err := uc.repo.ListAllProducts()
	if err != nil {
		uc.log.Println("ERROR FETCHING ALL PRODUCTS:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToListProducts, err)
	}

	uc.log.Println("ALL PRODUCTS FETCHED SUCCESSFULLY")
	return products, nil
}

func (uc Usecase) ArchiveProduct(productId uuid.UUID) error {
	uc.log.Println("ARCHIVING PRODUCT")
	err := uc.repo.ArchiveProduct(productId)
	if err != nil {
		uc.log.Println("ERROR ARCHIVING PRODUCT:", err)
		return fmt.Errorf("%s: %w", ErrFailedToArchiveProduct, err)
	}

	uc.log.Println("PRODUCT ARCHIVED SUCCESSFULLY")
	return nil
}
func (uc Usecase) AddProductToCatalog(productId, catalogId uuid.UUID, displayOrder int, userId uuid.UUID) ([]entity.ProductCatalog, error) {
	uc.log.Println("ADDING PRODUCT TO CATALOG")
	catalogs, err := uc.repo.AddProductToCatalog(
		productId,
		catalogId,
		displayOrder,
		userId,
	)
	if err != nil {
		uc.log.Println("ERROR ADDING PRODUCT TO CATALOG:", err)
		if errors.Is(err, ErrProductDoesNotExist) {
			return nil, fmt.Errorf("%w: product not found", err)
		}
		return nil, fmt.Errorf("%s: %w", ErrFailedToAddProductToCatalog, err)
	}

	uc.log.Println("PRODUCT ADDED TO CATALOG SUCCESSFULLY")
	return catalogs, nil
}

func (uc Usecase) RemoveProductFromCatalog(productId, catalogId uuid.UUID, userId uuid.UUID) ([]entity.ProductCatalog, error) {
	uc.log.Println("REMOVING PRODUCT FROM CATALOG")
	productCatalogs, err := uc.repo.RemoveProductFromCatalog(
		productId,
		catalogId,
		userId)
	if err != nil {
		uc.log.Println("ERROR REMOVING PRODUCT FROM CATALOG:", err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToRemoveProductFromCatalog, err)
	}

	uc.log.Println("PRODUCT REMOVED FROM CATALOG SUCCESSFULLY")
	return productCatalogs, nil
}

func validateProductInput(name, description string, price float64, currency, SKU string, weight float64, dimensions, imageURL string) error {
	if name == "" {
		return &JError{
			ErrorType: "INVALID_INPUT",
			Message:   "product name cannot be empty",
		}
	}
	if description == "" {
		return &JError{
			ErrorType: "INVALID_INPUT",
			Message:   "product description cannot be empty",
		}
	}
	if price <= 0 {
		return &JError{
			ErrorType: "INVALID_INPUT",
			Message:   "product price must be greater than zero",
		}
	}
	if currency == "" {
		return &JError{
			ErrorType: "INVALID_INPUT",
			Message:   "currency cannot be empty",
		}
	}
	skuRegex := `^[A-Za-z0-9_-]+$`
	if matched, _ := regexp.MatchString(skuRegex, SKU); !matched {
		return &JError{
			ErrorType: "INVALID_INPUT",
			Message:   "invalid SKU format",
		}
	}
	if weight <= 0 {
		return &JError{
			ErrorType: "INVALID_INPUT",
			Message:   "product weight must be greater than zero",
		}
	}
	if dimensions == "" {
		return &JError{
			ErrorType: "INVALID_INPUT",
			Message:   "product dimensions cannot be empty",
		}
	}
	if imageURL == "" {
		return &JError{
			ErrorType: "INVALID_INPUT",
			Message:   "image URL cannot be empty",
		}
	}
	return nil
}
