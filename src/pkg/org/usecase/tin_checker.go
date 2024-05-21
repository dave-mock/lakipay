package usecase

import "auth/src/pkg/org/core/entity"

type TINChecker interface {
	CheckTIN(tin string, usecase Usecase) (*entity.Organization, error)
}
