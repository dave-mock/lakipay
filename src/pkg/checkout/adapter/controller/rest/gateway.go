package rest

import (
	"auth/src/pkg/checkout/core/entity"
	"net/http"
	"time"
)

type Gateway struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`

	Icon string `json:"icon"`
	Type string `json:"type"`

	CanProcess bool `json:"can_process"`
	CanSettle  bool `json:"can_settle"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func encodeGatewayFromEntity(v entity.Gateway) Gateway {
	return Gateway{
		Id:         v.Id,
		Key:        v.Key,
		Name:       v.Name,
		Acronym:    v.Acronym,
		Icon:       v.Icon,
		Type:       string(v.Type),
		CanProcess: v.CanProcess,
		CanSettle:  v.CanSettle,
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  v.UpdatedAt,
	}
}

func (controller Controller) GetGateways(w http.ResponseWriter, r *http.Request) {

	res, err := controller.interactor.GetGateways()
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.Error(),
				Message: err.Error(),
			},
		}, 400)
		return
	}

	var serRes []Gateway = make([]Gateway, 0)

	for i := 0; i < len(res); i++ {
		serRes = append(serRes, encodeGatewayFromEntity(res[i]))
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    serRes,
	}, 200)
}
