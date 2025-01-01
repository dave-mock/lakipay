package usecase

import (
	"auth/src/pkg/erp/core/entity"
	"github.com/google/uuid"
	"time"
)



// Improved CreateCatalog method with better error handling and logging
func (uc Usecase) CreateCatalog(userId uuid.UUID, name string, description string, status string) (*entity.Catalog, error) {
	const ErrFailedToCreateCatalog = "FAILED_TO_CREATE_CATALOG"

	// Basic validation for inputs
	if name == "" || description == "" || status == "" {
		uc.log.Println("CREATE CATALOG ERROR: Invalid input parameters")
		return nil, Error{
			Type:    ErrFailedToCreateCatalog,
			Message: "All fields (name, description, status) are required",
		}
	}

	// Create catalog entity
	catalog := &entity.Catalog{
		Id:          uuid.New(),
		MerchantId:  userId,
		Name:        name,
		Description: description,
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   userId,
		UpdatedBy:   userId,
	}

	// Store catalog in repository
	catalog, err := uc.repo.CreateCatalog(userId, name, description, status)
	if err != nil {
		uc.log.Println("CREATE CATALOG ERROR: Failed to store catalog")
		return nil, Error{
			Type:    ErrFailedToCreateCatalog,
			Message: err.Error(),
		}
	}

	uc.log.Println("CREATE CATALOG SUCCESS: Catalog successfully created")
	return catalog, nil
}

// ListMerchantCatalogs method 
func (uc Usecase) ListMerchantCatalogs(merchantID string, userId uuid.UUID) ([]entity.Catalog, error) {
	const ErrCouldNotFindMerchantCatalogs = "COULD_NOT_FIND_MERCHANT_CATALOGS"

	// Get merchant catalogs from repository
	catalogs, err := uc.repo.ListMerchantCatalogs(merchantID, userId)
	if err != nil {
		uc.log.Println("LIST MERCHANT CATALOGS ERROR: Failed to retrieve catalogs")
		return nil, Error{
			Type:    ErrCouldNotFindMerchantCatalogs,
			Message: err.Error(),
		}
	}

	uc.log.Println("LIST MERCHANT CATALOGS SUCCESS")
	return catalogs, nil
}

// List all catalogs for a user with better error handling
func (uc Usecase) ListCatalogs(userId uuid.UUID) ([]entity.Catalog, error) {
	const ErrCouldNotFindCatalogs = "COULD_NOT_FIND_CATALOGS"

	// List all catalogs for the user
	catalogs, err := uc.repo.ListCatalogs(userId)
	if err != nil {
		uc.log.Println("LIST CATALOGS ERROR: Failed to retrieve catalogs")
		return nil, Error{
			Type:    ErrCouldNotFindCatalogs,
			Message: err.Error(),
		}
	}

	uc.log.Println("LIST CATALOGS SUCCESS")
	return catalogs, nil
}

// Get catalog details by ID with better error handling
func (uc Usecase) GetCatalog(catalogID string, userId uuid.UUID) (*entity.Catalog, error) {
	const ErrCatalogNotFound = "CATALOG_NOT_FOUND"

	// Retrieve the catalog from the repository
	catalog, err := uc.repo.GetCatalog(catalogID, userId)
	if err != nil {
		uc.log.Println("GET CATALOG ERROR: Catalog not found")
		return nil, Error{
			Type:    ErrCatalogNotFound,
			Message: err.Error(),
		}
	}

	uc.log.Println("GET CATALOG SUCCESS")
	return catalog, nil
}

// Archive a catalog with better error handling
func (uc Usecase) ArchiveCatalog(catalogID string, userId uuid.UUID) error {
	const ErrFailedToArchiveCatalog = "FAILED_TO_ARCHIVE_CATALOG"

	// Archive the catalog using the repository method
	err := uc.repo.ArchiveCatalog(catalogID, userId)
	if err != nil {
		uc.log.Println("ARCHIVE CATALOG ERROR: Failed to archive catalog")
		return Error{
			Type:    ErrFailedToArchiveCatalog,
			Message: err.Error(),
		}
	}

	uc.log.Println("ARCHIVED CATALOG SUCCESS")
	return nil
}

// Update catalog with better validation and error handling
func (uc Usecase) UpdateCatalog(userId uuid.UUID, catalogID, name, description, status string) (*entity.Catalog, error) {
	const ErrFailedToUpdateCatalog = "FAILED_TO_UPDATE_CATALOG"

	// Validate inputs before proceeding
	if name == "" || description == "" || status == "" {
		uc.log.Println("UPDATE CATALOG ERROR: Invalid input parameters")
		return nil, Error{
			Type:    ErrFailedToUpdateCatalog,
			Message: "Name, description, and status cannot be empty",
		}
	}

	// Update the catalog using the repository
	catalog, err := uc.repo.UpdateCatalog(userId, catalogID, name, description, status)
	if err != nil {
		uc.log.Println("UPDATE CATALOG ERROR: Failed to update catalog")
		return nil, Error{
			Type:    ErrFailedToUpdateCatalog,
			Message: err.Error(),
		}
	}

	uc.log.Println("UPDATE CATALOG SUCCESS")
	return catalog, nil
}

// Delete catalog with better error handling
func (uc Usecase) DeleteCatalog(catalogID string, userId uuid.UUID) error {
	const ErrFailedToDeleteCatalog = "FAILED_TO_DELETE_CATALOG"

	// Delete the catalog from the repository
	err := uc.repo.DeleteCatalog(catalogID, userId)
	if err != nil {
		uc.log.Println("DELETE CATALOG ERROR: Failed to delete catalog")
		return Error{
			Type:    ErrFailedToDeleteCatalog,
			Message: err.Error(),
		}
	}

	uc.log.Println("DELETE CATALOG SUCCESS")
	return nil
}
