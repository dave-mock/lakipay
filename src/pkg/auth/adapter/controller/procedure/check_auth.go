package procedure

import (
	"auth/src/pkg/auth/core/entity"
	"fmt"
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
