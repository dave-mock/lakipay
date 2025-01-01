package usecase

import (
	"auth/src/pkg/access_control/core/entity"
	"fmt"

	"github.com/google/uuid"
)

func (uc Usecase) CreateResource(name, description string, operations []uuid.UUID) (*entity.Resource, error) {
	const ErrFailedToCreateResource = "FAILED_TO_CREATE_RESOURCE"
	if name == "" {
		uc.log.Println("CREATE RESOURCE ERROR: Invalid name")
		return nil, Error{
			Type:    ErrFailedToCreateResource,
			Message: "Name is required",
		}
	}
	if description == "" {
		uc.log.Println("CREATE RESOURCE ERROR: Invalid description")
		return nil, Error{
			Type:    ErrFailedToCreateResource,
			Message: "Description is required",
		}
	}

	resource, err := uc.repo.CreateResource(name, description, operations)
	if err != nil {
		uc.log.Println("CREATE RESOURCE ERROR: Failed to create resource")
		return nil, Error{
			Type:    ErrFailedToCreateResource,
			Message: err.Error(),
		}
	}

	uc.log.Println("CREATE RESOURCE SUCCESS: Resource successfully created")
	return resource, nil
}

func (uc Usecase) UpdateResource(resourceID uuid.UUID, name, description string, operations []uuid.UUID) (*entity.Resource, error) {
	const ErrFailedToUpdateResource = "FAILED_TO_UPDATE_RESOURCE"

	if name == "" {
		uc.log.Println("UPDATE RESOURCE ERROR |||||| Invalid name")
		return nil, Error{
			Type:    ErrFailedToUpdateResource,
			Message: "Name is required",
		}
	}
	if description == "" {
		uc.log.Println("UPDATE RESOURCE ERROR |||||| Invalid description")
		return nil, Error{
			Type:    ErrFailedToUpdateResource,
			Message: "Description is required",
		}
	}

	resource, err := uc.repo.GetResourceByID(resourceID)
	if err != nil {
		uc.log.Printf("UPDATE RESOURCE ERROR: Resource with ID %v not found", resourceID)
		return nil, Error{
			Type:    ErrFailedToUpdateResource,
			Message: fmt.Sprintf("Resource with ID %v does not exist", resourceID),
		}
	}

	resource.Name = name
	resource.Description = description
	updatedResource, err := uc.repo.UpdateResource(resourceID, name, description, operations)
	if err != nil {
		uc.log.Println("UPDATE RESOURCE ERROR |||| Failed to update resource")
		return nil, Error{
			Type:    ErrFailedToUpdateResource,
			Message: err.Error(),
		}
	}

	uc.log.Printf("UPDATE RESOURCE SUCCESS|||| Resource with ID %v successfully updated", resourceID)
	return updatedResource, nil
}

func (uc Usecase) GetResourceByID(resourceID uuid.UUID) (*entity.Resource, error) {
	resource, err := uc.repo.GetResourceByID(resourceID)
	if err != nil {
		uc.log.Println("GET RESOURCE ERROR||| Resource not found")
		return nil, Error{
			Type:    "FAILED_TO_GET_RESOURCE",
			Message: err.Error(),
		}
	}
	return resource, nil
}

func (uc Usecase) ListResources() ([]entity.Resource, error) {
	const ErrFailedToListResources = "FAILED_TO_LIST_RESOURCES"
	resources, err := uc.repo.ListResources()
	if err != nil {
		uc.log.Println("LIST RESOURCES ERROR  ||||| Failed to list resources")
		return nil, Error{
			Type:    ErrFailedToListResources,
			Message: err.Error(),
		}
	}

	uc.log.Println("LIST RESOURCES SUCCESS |||||| Resources successfully listed")
	return resources, nil
}

func (uc Usecase) DeleteResource(resourceID uuid.UUID) error {
	const ErrFailedToDeleteResource = "FAILED_TO_DELETE_RESOURCE"
	err := uc.repo.DeleteResource(resourceID)
	if err != nil {
		uc.log.Println("DELETE RESOURCE ERROR  |||||| Failed to delete resource")
		return Error{
			Type:    ErrFailedToDeleteResource,
			Message: err.Error(),
		}
	}

	uc.log.Println("DELETE RESOURCE SUCCESS ||||||||| Resource successfully deleted")
	return nil
}
