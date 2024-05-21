package rest

import (
	"auth/src/pkg/org/core/entity"
	"net/http"
)

func (controller Controller) CheckTIN(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		TIN string
	}

	type Response struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data,omitempty"`
		Error   error       `json:"error,omitempty"`
	}

	var req Request

	req.TIN = r.URL.Query().Get("tin")

	if req.TIN == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: "Please provide a Tax Identification Number",
			},
		}, http.StatusBadRequest)
		return
	}

	org, err := controller.interactor.CheckTIN(req.TIN)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error:   err,
		}, http.StatusBadRequest)
		return
	}

	if org == nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "NO_ORG",
				Message: "No organization information found for the provided TIN",
			},
		}, http.StatusBadRequest)
		return
	}

	var deps []Department = make([]Department, 0)

	for i := 0; i < len(org.Departments); i++ {
		var cats []struct {
			Id   int64  "json:\"id\""
			Name string "json:\"name\""
		} = make([]struct {
			Id   int64  "json:\"id\""
			Name string "json:\"name\""
		}, 0)

		for _, c := range org.Departments[i].Categories {
			cats = append(cats, struct {
				Id   int64  "json:\"id\""
				Name string "json:\"name\""
			}{
				// Id:   c.Id,
				Name: c.Name,
			})
		}

		deps = append(deps, Department{
			Name:       org.Departments[i].Name,
			Categories: cats,
		})
	}

	var orgCategory *struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
	}

	if org.Category != nil {
		orgCategory = &struct {
			Id   string "json:\"id\""
			Name string "json:\"name\""
		}{
			Id:   org.Category.Id.String(),
			Name: org.Category.Name,
		}
	}

	var legalCondition *struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
	}

	if org.LegalCondition != nil {
		legalCondition = &struct {
			Id   string "json:\"id\""
			Name string "json:\"name\""
		}{
			Id:   org.LegalCondition.Id.String(),
			Name: org.LegalCondition.Name,
		}
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: Organization{
			Name:          org.Name,
			Capital:       org.Capital,
			RegDate:       org.RegDate,
			Country:       org.Country,
			Category:      orgCategory,
			LegalCondtion: legalCondition,
			Departments:   deps,
			Details: EthBusOrg{
				TIN:     org.Details.(entity.EthBusOrg).TIN,
				TinFile: org.Details.(entity.EthBusOrg).TINFile,
				RegNo:   org.Details.(entity.EthBusOrg).RegNo,
				RegFile: org.Details.(entity.EthBusOrg).RegFile,
			},
		},
	}, http.StatusOK)
}
