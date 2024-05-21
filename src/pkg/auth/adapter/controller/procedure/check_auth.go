package procedure

import "auth/src/pkg/auth/core/entity"

func (controller Controller) GetCheckAuth(token string) (*entity.Session, error) {
	session, err := controller.interactor.CheckSession(token)
	if err != nil {
		return nil, Error{
			Type:    "UNAUTHORIZED",
			Message: err.Error(),
		}
	}

	return session, nil
}
