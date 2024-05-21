package rest

import (
	"net/http"
	"strings"
)

func (controller Controller) GetCheckSession(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Token string
	}

	var req Request

	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) != 2 {
		SendJSONResponse(w, Error{}, http.StatusUnauthorized)
		return
	}

	req.Token = strings.Split(token, " ")[1]

	_, err := controller.interactor.CheckSession(req.Token)
	if err != nil {
		SendJSONResponse(w, Error{
			Type:    err.Error(),
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "",
	}, http.StatusOK)
}
