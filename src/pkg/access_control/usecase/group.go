package usecase

import (
	"auth/src/pkg/access_control/core/entity"

	"github.com/google/uuid"
)

func (uc Usecase) CreateGroup(title string) (*entity.Group, error) {
	const ErrFailedToCreateGroup = "FAILED_TO_CREATE_GROUP"
	if title == "" {
		uc.log.Println("CREATE GROUP ERROR: Invalid title")
		return nil, Error{
			Type:    ErrFailedToCreateGroup,
			Message: "Title is required",
		}
	}

	group, err := uc.repo.CreateGroup(title)
	if err != nil {
		uc.log.Println("CREATE GROUP ERROR: Failed to create group")
		return nil, Error{
			Type:    ErrFailedToCreateGroup,
			Message: err.Error(),
		}
	}

	uc.log.Println("CREATE GROUP SUCCESS: Group successfully created")
	return group, nil
}

func (uc Usecase) UpdateGroup(groupID uuid.UUID, title string) (*entity.Group, error) {
	const ErrFailedToUpdateGroup = "FAILED_TO_UPDATE_GROUP"
	if title == "" {
		uc.log.Println("UPDATE GROUP ERROR: Invalid title")
		return nil, Error{
			Type:    ErrFailedToUpdateGroup,
			Message: "Title is required",
		}
	}
	group, err := uc.repo.UpdateGroup(groupID, title)
	if err != nil {
		uc.log.Println("UPDATE GROUP ERROR: Failed to update group")
		return nil, Error{
			Type:    ErrFailedToUpdateGroup,
			Message: err.Error(),
		}
	}

	uc.log.Println("UPDATE GROUP SUCCESS: Group successfully updated")
	return group, nil
}

func (uc Usecase) DeleteGroup(groupID uuid.UUID) error {
	const ErrFailedToDeleteGroup = "FAILED_TO_DELETE_GROUP"
	err := uc.repo.DeleteGroup(groupID)
	if err != nil {
		uc.log.Println("DELETE GROUP ERROR: Failed to delete group")
		return Error{
			Type:    ErrFailedToDeleteGroup,
			Message: err.Error(),
		}
	}

	uc.log.Println("DELETE GROUP SUCCESS: Group successfully deleted")
	return nil
}

func (uc Usecase) ListGroups() ([]entity.Group, error) {
	const ErrFailedToListGroups = "FAILED_TO_LIST_GROUPS"
	groups, err := uc.repo.ListGroups()
	if err != nil {
		uc.log.Println("|||| Failed to retrieve groups")
		return nil, Error{
			Type:    ErrFailedToListGroups,
			Message: err.Error(),
		}
	}

	uc.log.Println("||||| Groups successfully retrieved")
	return groups, nil
}

func (uc Usecase) AddUserToGroup(userID, groupID uuid.UUID) error {
	const ErrFailedToAddUserToGroup = "FAILED_TO_ADD_USER_TO_GROUP"
	err := uc.repo.AddUserToGroup(userID, groupID)
	if err != nil {
		uc.log.Println("ADD USER TO GROUP ERROR |||||| Failed to add user to group")
		return Error{
			Type:    ErrFailedToAddUserToGroup,
			Message: err.Error(),
		}
	}

	uc.log.Println("ADD USER TO GROUP SUCCESS |||||| User successfully added to group")
	return nil
}

func (uc Usecase) RemoveUserFromGroup(userID, groupID uuid.UUID) error {
	const ErrFailedToRemoveUserFromGroup = "FAILED_TO_REMOVE_USER_FROM_GROUP"
	err := uc.repo.RemoveUserFromGroup(userID, groupID)
	if err != nil {
		uc.log.Println("REMOVE USER FROM GROUP ERROR|||| Failed to remove user from group")
		return Error{
			Type:    ErrFailedToRemoveUserFromGroup,
			Message: err.Error(),
		}
	}

	uc.log.Println("REMOVE USER FROM GROUP SUCCESS ||||||| User successfully removed from group")
	return nil
}

func (uc Usecase) ListUserGroups(userID uuid.UUID) ([]entity.Group, error) {
	const ErrFailedToListUserGroups = "FAILED_TO_LIST_USER_GROUPS"

	groups, err := uc.repo.ListUserGroups(userID)
	if err != nil {
		uc.log.Println("LIST USER GROUPS ERROR|||| Failed to retrieve user groups")
		return nil, Error{
			Type:    ErrFailedToListUserGroups,
			Message: err.Error(),
		}
	}

	uc.log.Println("LIST USER GROUPS SUCCESS |||| User groups successfully retrieved")
	return groups, nil
}

func (uc Usecase) ListGroupUsers(groupID uuid.UUID) ([]entity.User, error) {
	const ErrFailedToListGroupUsers = "FAILED_TO_LIST_GROUP_USERS"

	users, err := uc.repo.ListGroupUsers(groupID)
	if err != nil {
		uc.log.Println("LIST GROUP USERS ERROR |||| Failed to retrieve group users")
		return nil, Error{
			Type:    ErrFailedToListGroupUsers,
			Message: err.Error(),
		}
	}

	uc.log.Println("LIST GROUP USERS SUCCESS |||Group users successfully retrieved")
	return users, nil
}
