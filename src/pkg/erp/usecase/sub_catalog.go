package usecase

import (
	"auth/src/pkg/erp/core/entity"

	"github.com/google/uuid"
)

func (uc Usecase) CreateSubCatalog(subCatalog *entity.SubCatalog) error {
	const ErrFailedToCreateSubCatalog = "FAILED_TO_CREATE_SUB_CATALOG"

	if subCatalog.Name == "" || subCatalog.Description == "" || subCatalog.Status == "" {
		uc.log.Println("CREATE SUB CATALOG ERROR: Invalid input parameters")
		return Error{
			Type:    ErrFailedToCreateSubCatalog,
			Message: "All fields (name, description, status) are required",
		}
	}

	err := uc.repo.CreateSubCatalog(subCatalog)
	if err != nil {
		uc.log.Println("CREATE SUB CATALOG ERROR: Failed to store sub-catalog")
		return Error{
			Type:    ErrFailedToCreateSubCatalog,
			Message: err.Error(),
		}
	}

	uc.log.Println("CREATE SUB CATALOG SUCCESS: Sub-catalog successfully created")
	return nil
}

func (uc Usecase) ListSubCatalogs(userId uuid.UUID) ([]entity.SubCatalog, error) {
	const ErrCouldNotFindSubCatalogs = "COULD_NOT_FIND_SUB_CATALOGS"

	subCatalogs, err := uc.repo.ListSubCatalogs(userId)
	if err != nil {
		uc.log.Println("LIST SUB CATALOGS ERROR: Failed to retrieve sub-catalogs")
		return nil, Error{
			Type:    ErrCouldNotFindSubCatalogs,
			Message: err.Error(),
		}
	}

	uc.log.Println("LIST SUB CATALOGS SUCCESS")
	return subCatalogs, nil
}

// GetSubCatalog retrieves details of a sub-catalog by its ID.
func (uc Usecase) GetSubCatalog(subCatalogID uuid.UUID, userId uuid.UUID) (*entity.SubCatalog, error) {
	const ErrSubCatalogNotFound = "SUB_CATALOG_NOT_FOUND"

	// Retrieve the sub-catalog from the repository
	subCatalog, err := uc.repo.GetSubCatalog(subCatalogID, userId)
	if err != nil {
		uc.log.Println("GET SUB CATALOG ERROR: Sub-catalog not found")
		return nil, Error{
			Type:    ErrSubCatalogNotFound,
			Message: err.Error(),
		}
	}

	uc.log.Println("GET SUB CATALOG SUCCESS")
	return subCatalog, nil
}

func (uc Usecase) ArchiveSubCatalog(subCatalogID uuid.UUID, userId uuid.UUID) error {
	const ErrFailedToArchiveSubCatalog = "FAILED_TO_ARCHIVE_SUB_CATALOG"
	err := uc.repo.ArchiveSubCatalog(subCatalogID, userId)
	if err != nil {
		uc.log.Println("ARCHIVE SUB CATALOG ERROR: Failed to archive sub-catalog")
		return Error{
			Type:    ErrFailedToArchiveSubCatalog,
			Message: err.Error(),
		}
	}

	uc.log.Println("ARCHIVE SUB CATALOG SUCCESS")
	return nil
}

// UpdateSubCatalog updates an existing sub-catalog.
func (uc Usecase) UpdateSubCatalog(userId uuid.UUID, subCatalog *entity.SubCatalog) (*entity.SubCatalog, error) {
	const ErrFailedToUpdateSubCatalog = "FAILED_TO_UPDATE_SUB_CATALOG"

	// Validate inputs before proceeding
	if subCatalog.Name == "" || subCatalog.Description == "" || subCatalog.Status == "" {
		uc.log.Println("UPDATE SUB CATALOG ERROR: Invalid input parameters")
		return nil, Error{
			Type:    ErrFailedToUpdateSubCatalog,
			Message: "Name, description, and status cannot be empty",
		}
	}

	// Update the sub-catalog using the repository
	err := uc.repo.UpdateSubCatalog(subCatalog)
	if err != nil {
		uc.log.Println("UPDATE SUB CATALOG ERROR: Failed to update sub-catalog")
		return nil, Error{
			Type:    ErrFailedToUpdateSubCatalog,
			Message: err.Error(),
		}
	}

	uc.log.Println("UPDATE SUB CATALOG SUCCESS")
	return subCatalog, nil
}

func (uc Usecase) DeleteSubCatalog(subCatalogID uuid.UUID, userId uuid.UUID) error {
	const ErrFailedToDeleteSubCatalog = "FAILED_TO_DELETE_SUB_CATALOG"

	err := uc.repo.DeleteSubCatalog(subCatalogID, userId)
	if err != nil {
		uc.log.Println("DELETE SUB CATALOG ERROR: Failed to delete sub-catalog")
		return Error{
			Type:    ErrFailedToDeleteSubCatalog,
			Message: err.Error(),
		}
	}

	uc.log.Println("DELETE SUB CATALOG SUCCESS")
	return nil
}
