package rest

import (
	"auth/src/pkg/org/core/entity"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Organization
type Organization struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Logo        string    `json:"logo"`
	Capital     float64   `json:"capital"`
	RegDate     time.Time `json:"reg_date"`
	Country     string    `json:"country"`
	Category    *struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	LegalCondtion *struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"legal_condition"`
	Taxes           []struct{}         `json:"taxes"`
	Associates      []struct{}         `json:"associates"`
	Departments     []Department       `json:"departments"`
	Details         interface{}        `json:"details"`
	Status          VerificationStatus `json:"status"`
	RetentionStatus RetentionStatus    `json:"retention_status"`
	CreatedAt       time.Time          `json:"created_at"`
}

func newOrgFromEntity(entity entity.Organization) Organization {

	var org Organization

	org.Id = entity.Id.String()
	org.Name = entity.Name
	org.Description = entity.Description
	org.Logo = entity.Logo
	org.Capital = entity.Capital
	org.RegDate = entity.RegDate
	org.Country = entity.Country
	org.Category = &struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
	}{}
	org.LegalCondtion = &struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
	}{}
	org.Taxes = []struct{}{}
	org.Associates = []struct{}{}
	org.Departments = []Department{}
	org.Details = entity.Details
	org.Status = VerificationStatus{}
	org.RetentionStatus = RetentionStatus{}
	org.CreatedAt = entity.CreatedAt

	return org
}

type VerificationStatus struct {
	Verified bool   `json:"verified"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type RetentionStatus struct {
	CanRetain bool   `json:"can_retain"`
	File      string `json:"file"`
}

// Eth Bus Org
type EthBusOrg struct {
	TIN     string             `json:"tin"`
	TinFile string             `json:"tin_file"`
	RegNo   string             `json:"reg_no"`
	RegFile string             `json:"reg_file"`
	Status  VerificationStatus `json:"status"`
}

func (controller Controller) GetOrganization(w http.ResponseWriter, r *http.Request) {

	// Request
	type Request struct {
		OrgId uuid.UUID
		Token string
	}

	var req Request

	// OrgId

	// Token
	if r.Header.Get("Authorization") != "" && len(strings.Split(r.Header.Get("Authorization"), " ")) == 2 {
		req.Token = strings.Split(r.Header.Get("Authorization"), " ")[1]
	}

	// uc
	org, err := controller.interactor.GetOrganization(req.Token, req.OrgId)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   Error{},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    newOrgFromEntity(*org),
	}, http.StatusOK)
}
