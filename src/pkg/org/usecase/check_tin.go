package usecase

import (
	"auth/src/pkg/org/core/entity"
)

func (uc Usecase) CheckTIN(tin string) (*entity.Organization, error) {

	return uc.tinchecker.CheckTIN(tin, uc)
}
