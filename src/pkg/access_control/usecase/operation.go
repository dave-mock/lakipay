package usecase

import (
	"auth/src/pkg/access_control/core/entity"
	"fmt"

	"github.com/google/uuid"
)

func (uc Usecase) CreateOperations(name, description string) (*entity.Operation, error) {
	const ErrFailedToCreateOperations = "FAILED_TO_CREATE_Operations"
	if name == "" {
		uc.log.Println("CREATE Operations ERROR Invalid name")
		return nil, Error{
			Type:    ErrFailedToCreateOperations,
			Message: "Name is required",
		}
	}
	if description == "" {
		uc.log.Println("CREATE Operations ERROR Invalid description")
		return nil, Error{
			Type:    ErrFailedToCreateOperations,
			Message: "Description is required",
		}
	}

	operation, err := uc.repo.CreateOperations(name, description)
	if err != nil {
		uc.log.Println("CREATE Operations ERROR Failed to create operation")
		return nil, Error{
			Type:    ErrFailedToCreateOperations,
			Message: err.Error(),
		}
	}

	uc.log.Println("CREATE Operations SUCCESS Operations successfully created")
	return operation, nil
}

func (uc Usecase) UpdateOperations(operationID uuid.UUID, name, description string) (*entity.Operation, error) {
	const ErrFailedToUpdateOperations = "FAILED_TO_UPDATE_Operations"

	if name == "" {
		uc.log.Println("UPDATE Operations ERROR Invalid name")
		return nil, Error{
			Type:    ErrFailedToUpdateOperations,
			Message: "Name is required",
		}
	}
	if description == "" {
		uc.log.Println("UPDATE Operations ERROR Invalid description")
		return nil, Error{
			Type:    ErrFailedToUpdateOperations,
			Message: "Description is required",
		}
	}

	operation, err := uc.repo.GetOperationsByID(operationID)
	if err != nil {
		uc.log.Printf("UPDATE Operations ERROR Operation with ID %v not found", operationID)
		return nil, Error{
			Type:    ErrFailedToUpdateOperations,
			Message: fmt.Sprintf("Operation with ID %v does not exist", operationID),
		}
	}

	operation.Name = name
	operation.Description = description
	updatedOperation, err := uc.repo.UpdateOperations(operationID, name, description)
	if err != nil {
		uc.log.Println("UPDATE Operations ERROR Failed to update operation")
		return nil, Error{
			Type:    ErrFailedToUpdateOperations,
			Message: err.Error(),
		}
	}

	uc.log.Printf("UPDATE Operations SUCCESS Operation with ID %v successfully updated", operationID)
	return updatedOperation, nil
}

func (uc Usecase) GetOperationsByID(operationID uuid.UUID) (*entity.Operation, error) {
	operation, err := uc.repo.GetOperationsByID(operationID)
	if err != nil {
		uc.log.Println("GET Operations ERROR Operation not found")
		return nil, Error{
			Type:    "FAILED_TO_GET_Operations",
			Message: err.Error(),
		}
	}
	return operation, nil
}

func (uc Usecase) ListOperations() ([]*entity.Operation, error) {
	const ErrFailedToListOperations = "FAILED_TO_LIST_Operations"
	operations, err := uc.repo.ListOperations()
	if err != nil {
		uc.log.Println("LIST Operations ERROR ailed to list operations")
		return nil, Error{
			Type:    ErrFailedToListOperations,
			Message: err.Error(),
		}
	}

	var operationsPtr []*entity.Operation
	for _, operation := range operations {
		operationsPtr = append(operationsPtr, operation)
	}

	uc.log.Println("LIST Operations SUCCESS Operations successfully listed")
	return operationsPtr, nil
}

func (uc Usecase) DeleteOperations(operationID uuid.UUID) error {
	const ErrFailedToDeleteOperations = "FAILED_TO_DELETE_Operations"
	err := uc.repo.DeleteOperations(operationID)
	if err != nil {
		uc.log.Println("DELETE Operations ERROR Failed to delete operation")
		return Error{
			Type:    ErrFailedToDeleteOperations,
			Message: err.Error(),
		}
	}

	uc.log.Println("DELETE Operations SUCCESS Operation successfully deleted")
	return nil
}
