package usecase

import (
	"auth/src/pkg/access_control/core/entity"
	"fmt"

	"github.com/google/uuid"
)

func (uc Usecase) CreatePermission(resourceID string, resource uuid.UUID, operation []uuid.UUID, effect string) (*entity.Permission, error) {
	const ErrFailedToCreatePermission = "FAILED_TO_CREATE_PERMISSION"
	permission, err := uc.repo.CreatePermission(resourceID, resource, operation, effect)
	if err != nil {
		uc.log.Println("CREATE PERMISSION ERROR: Failed to create permission")
		return nil, Error{
			Type:    ErrFailedToCreatePermission,
			Message: err.Error(),
		}
	}

	uc.log.Println("CREATE PERMISSION SUCCESS: Permission successfully created")
	return permission, nil
}

func (uc Usecase) UpdatePermission(permissionID uuid.UUID, resourceID string, resource uuid.UUID, operation []uuid.UUID, effect string) (*entity.Permission, error) {
	const ErrFailedToUpdatePermission = "FAILED_TO_UPDATE_PERMISSION"
	permission, err := uc.repo.UpdatePermission(permissionID, resourceID, resource, operation, effect)
	if err != nil {
		uc.log.Println("UPDATE PERMISSION ERROR: Failed to update permission")
		return nil, Error{
			Type:    ErrFailedToUpdatePermission,
			Message: err.Error(),
		}
	}
	uc.log.Println("UPDATE PERMISSION SUCCESS: Permission successfully updated")
	return permission, nil
}

func (uc Usecase) DeletePermission(resourceID uuid.UUID, permissionID uuid.UUID) error {
	const ErrFailedToDeletePermission = "FAILED_TO_DELETE_PERMISSION"
	err := uc.repo.DeletePermission(resourceID, permissionID)
	if err != nil {
		uc.log.Println("DELETE PERMISSION ERROR: Failed to delete permission")
		return fmt.Errorf("%s: %v", ErrFailedToDeletePermission, err)
	}

	uc.log.Println("DELETE PERMISSION SUCCESS: Permission deleted successfully")
	return nil
}

func (uc Usecase) ListPermissions() ([]entity.Permission, error) {
	const ErrFailedToListPermissions = "FAILED_TO_LIST_PERMISSIONS"
	permissions, err := uc.repo.ListPermissions()
	if err != nil {
		uc.log.Println("LIST PERMISSIONS ERROR: Failed to retrieve permissions")
		return nil, Error{
			Type:    ErrFailedToListPermissions,
			Message: err.Error(),
		}
	}

	uc.log.Println("LIST PERMISSIONS SUCCESS: Permissions successfully retrieved")
	return permissions, nil
}

func (uc Usecase) ListUsers() ([]entity.User, error) {
	const ErrFailedToListUsers = "FAILED_TO_LIST_USERS"
	users, err := uc.repo.ListUsers()
	if err != nil {
		uc.log.Println("LIST USERS ERROR: Failed to retrieve users")
		return nil, Error{
			Type:    ErrFailedToListUsers,
			Message: err.Error(),
		}
	}

	uc.log.Println("LIST USERS SUCCESS: Users successfully retrieved")
	return users, nil
}

func (uc Usecase) ListUserPermissions(userID uuid.UUID) ([]entity.Permission, error) {
	const ErrFailedToListUserPermissions = "FAILED_TO_LIST_USER_PERMISSIONS"
	permissions, err := uc.repo.ListUserPermissions(userID)
	if err != nil {
		uc.log.Println("LIST USER PERMISSIONS ERROR: Failed to retrieve user permissions")
		return nil, Error{
			Type:    ErrFailedToListUserPermissions,
			Message: err.Error(),
		}
	}

	uc.log.Println("LIST USER PERMISSIONS SUCCESS: User permissions successfully retrieved")
	return permissions, nil
}

func (uc Usecase) ListGroupPermissions(groupID uuid.UUID) ([]entity.Permission, error) {
	const ErrFailedToListGroupPermissions = "FAILED_TO_LIST_GROUP_PERMISSIONS"
	permissions, err := uc.repo.ListGroupPermissions(groupID)
	if err != nil {
		uc.log.Println("LIST GROUP PERMISSIONS ERROR: Failed to retrieve group permissions")
		return nil, Error{
			Type:    ErrFailedToListGroupPermissions,
			Message: err.Error(),
		}
	}

	uc.log.Println("LIST GROUP PERMISSIONS SUCCESS: Group permissions successfully retrieved")
	return permissions, nil
}
