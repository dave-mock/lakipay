package usecase

import (
	"auth/src/pkg/access_control/core/entity"

	"github.com/google/uuid"
)

type Repo interface {
	CreateGroup(title string) (*entity.Group, error)
	UpdateGroup(groupID uuid.UUID, title string) (*entity.Group, error)
	ListGroups() ([]entity.Group, error)
	DeleteGroup(groupID uuid.UUID) error
	CreatePermission(resourceID uuid.UUID, operation, effect string) (*entity.Permission, error)
	UpdatePermission(permissionID uuid.UUID, operation, effect string) (*entity.Permission, error)
	ListPermissions() ([]entity.Permission, error)
	DeletePermission(resourceID uuid.UUID, permissionID uuid.UUID) error
	AddUserToGroup(userID, groupID uuid.UUID) error
	RemoveUserFromGroup(userID, groupID uuid.UUID) error
	ListUserGroups(userID uuid.UUID) ([]entity.Group, error)
	ListGroupUsers(groupID uuid.UUID) ([]entity.User, error)
	GrantPermissionToUser(userID uuid.UUID, permissionID uuid.UUID) error
	RevokePermissionFromUser(userID uuid.UUID, permissionID uuid.UUID) error
	GrantPermissionToGroup(groupID uuid.UUID, permissionID uuid.UUID) error
	RevokePermissionFromGroup(groupID uuid.UUID, permissionID uuid.UUID) error
	ListUserPermissions(userID uuid.UUID) ([]entity.Permission, error)
	ListGroupPermissions(groupID uuid.UUID) ([]entity.Permission, error)
	// CRUD methods for resources
	CreateResource(name, description string) (*entity.Resource, error)
	ListResources() ([]entity.Resource, error)
	UpdateResource(resourceID uuid.UUID, name, description string) (*entity.Resource, error)
	DeleteResource(resourceID uuid.UUID) error
	GetResourceByID(resourceID uuid.UUID) (*entity.Resource, error)
}
