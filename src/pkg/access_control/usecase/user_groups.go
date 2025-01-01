package usecase

import (
	"github.com/google/uuid"
)

func (uc Usecase) GrantPermissionToUser(userID uuid.UUID, permissionID uuid.UUID) error {
	const ErrFailedToGrantPermissionToUser = "FAILED_TO_GRANT_PERMISSION_TO_USER"
	err := uc.repo.GrantPermissionToUser(userID, permissionID)
	if err != nil {
		uc.log.Println("GRANT PERMISSION TO USER ERROR |||||| Failed to grant permission to user")
		return Error{
			Type:    ErrFailedToGrantPermissionToUser,
			Message: err.Error(),
		}
	}

	uc.log.Println("GRANT PERMISSION TO USER SUCCESS |||||| Permission successfully granted to user")
	return nil
}

func (uc Usecase) RevokePermissionFromUser(userID uuid.UUID, permissionID uuid.UUID) error {
	const ErrFailedToRevokePermissionFromUser = "FAILED_TO_REVOKE_PERMISSION_FROM_USER"
	err := uc.repo.RevokePermissionFromUser(userID, permissionID)
	if err != nil {
		uc.log.Println("REVOKE PERMISSION FROM USER ERROR||||||| Failed to revoke permission from user")
		return Error{
			Type:    ErrFailedToRevokePermissionFromUser,
			Message: err.Error(),
		}
	}

	uc.log.Println("REVOKE PERMISSION FROM USER SUCCESS |||||||| Permission successfully revoked from user")
	return nil
}

func (uc Usecase) GrantPermissionToGroup(groupID uuid.UUID, permissionID uuid.UUID) error {
	const ErrFailedToGrantPermissionToGroup = "FAILED_TO_GRANT_PERMISSION_TO_GROUP"
	err := uc.repo.GrantPermissionToGroup(groupID, permissionID)
	if err != nil {
		uc.log.Println("GRANT PERMISSION TO GROUP ERROR ||||||  Failed to grant permission to group")
		return Error{
			Type:    ErrFailedToGrantPermissionToGroup,
			Message: err.Error(),
		}
	}

	uc.log.Println("GRANT PERMISSION TO GROUP SUCCESS |||||||| Permission successfully granted to group")
	return nil
}

func (uc Usecase) RevokePermissionFromGroup(groupID uuid.UUID, permissionID uuid.UUID) error {
	const ErrFailedToRevokePermissionFromGroup = "FAILED_TO_REVOKE_PERMISSION_FROM_GROUP"
	err := uc.repo.RevokePermissionFromGroup(groupID, permissionID)
	if err != nil {
		uc.log.Println("REVOKE PERMISSION FROM GROUP ERROR |||||||| Failed to revoke permission from group")
		return Error{
			Type:    ErrFailedToRevokePermissionFromGroup,
			Message: err.Error(),
		}
	}

	uc.log.Println("REVOKE PERMISSION FROM GROUP SUCCESS||||||||| Permission successfully revoked from group")
	return nil
}
