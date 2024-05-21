package usecase

import (
	"auth/src/pkg/org/core/entity"
	"time"

	"github.com/google/uuid"
)

type OrgInteractor interface {
	GetOrganization(token string, id uuid.UUID) (*entity.Organization, error)
}

func (uc Usecase) GetOrganization(token string, id uuid.UUID) (*entity.Organization, error) {

	return &entity.Organization{
		Id:             uuid.New(),
		Name:           "LakiPay Finanacial Technologies SC",
		Description:    "A payment gateway service",
		Logo:           "http://",
		Capital:        20000000.00,
		RegDate:        time.Now(),
		Country:        "ET",
		Category:       &entity.Category{},
		LegalCondition: &entity.LegalCondition{},
		Taxes:          []entity.OrganizationTax{},
		Associates:     []entity.Associate{},
		Departments:    []entity.Department{},
		Details: entity.EthBusOrg{
			TIN:     "0000000000",
			TINFile: "http://",
			RegNo:   "00000000",
			RegFile: "http://",
			Status:  entity.VerificationStatus{},
		},
		Status:          entity.VerificationStatus{},
		RetentionStatus: entity.RetentionStatus{},
		CreatedAt:       time.Now(),
	}, nil
}
