package entity

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	Id              uuid.UUID
	Name            string
	Description     string
	Logo            string
	Capital         float64
	RegDate         time.Time
	Country         string
	Category        *Category
	LegalCondition  *LegalCondition
	Taxes           []OrganizationTax
	Associates      []Associate // List of users
	Departments     []Department
	Details         interface{}
	Status          VerificationStatus
	RetentionStatus RetentionStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Associate struct {
	UserId   uuid.UUID
	Position string
}

type OrganizationTax struct {
	Tax    Tax
	File   string
	Status VerificationStatus
}

// ETH - ORG
type EthBusOrg struct {
	TIN     string
	TINFile string
	RegNo   string
	RegFile string
	Status  VerificationStatus
}

//

type RetentionStatus struct {
	CanRetain bool
	File      string
}
