package procedure

import (
	"auth/src/pkg/auth/core/entity"
	"fmt"

	"github.com/google/uuid"
)

func (controller Controller) GetCheckAuth(token string) (*entity.Session, error) {
	fmt.Println("========================= ", token)
	session, err := controller.interactor.CheckSession(token)
	fmt.Println("==================")
	if err != nil {
		return nil, Error{
			Type:    "UNAUTHORIZED",
			Message: err.Error(),
		}
	}

	return session, nil
}

func (controller Controller) HasPermission(userID uuid.UUID, requiredPermission entity.Permission) (bool, error) {
	fmt.Println("========================= ", userID, requiredPermission)
	hasPermission, err := controller.interactor.CheckPermission(userID, requiredPermission)
	if err != nil {
		return false, Error{
			Type:    "FORBIDDEN",
			Message: err.Error(),
		}
	}
	return hasPermission, nil
}
