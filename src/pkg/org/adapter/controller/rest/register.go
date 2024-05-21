package rest

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (controller Controller) GetRegisterOrganization(w http.ResponseWriter, r *http.Request) {
	// Register Company
	type Request struct {
		Name           string    `json:"name"`
		Capital        float64   `json:"capital"`
		RegDate        time.Time `json:"reg_date"`
		Type           uuid.UUID `json:"type"`
		LegalCondition uuid.UUID `json:"legal_condition"`
		Country        string    `json:"country"`
		Taxes          []struct {
			Id    uuid.UUID
			Files []string
		} `json:"taxes"`
		Departments []struct {
			Name       string `json:"name"`
			Categories []struct {
				Id uuid.UUID
			} `json:"categories"`
			Details interface{} `json:"details"`
		} `json:"departments"`
		Details interface{} `json:"details"`
	}
}

func (controller Controller) GetInitOrgRegistration(w http.ResponseWriter, r *http.Request) {
	// time.Sleep(10 * time.Second)
	http.Redirect(w, r, "http://127.0.0.1:3000", http.StatusSeeOther)
}
