package rest

import (
	"auth/src/pkg/auth/usecase"
	"encoding/json"
	"net"
	"net/http"
)

func (controller Controller) GetInitAuth(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[AUTH] [ADAPTER] [CONTROLLER] [REST] [GetInitAuth] ")
	// Accept req with
	// Device info
	// 		- IP
	// 		- Name
	// 		- Agent
	// Phone
	//		- Code
	// 		- Number

	// Request
	type Request struct {
		Device struct {
			IP    net.IPAddr `json:"ip"`
			Name  string     `json:"name"`
			Agent string     `json:"agent"`
		} `json:"device"`
		Phone struct {
			Prefix string `json:"prefix"`
			Number string `json:"number"`
		} `json:"phone"`
	}

	// Response
	type Response struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data,omitempty"`
		Error   error       `json:"error,omitempty"`
	}

	// Parse Request
	var req Request

	req.Device = struct {
		IP    net.IPAddr "json:\"ip\""
		Name  string     "json:\"name\""
		Agent string     "json:\"agent\""
	}{
		IP: net.IPAddr{
			IP: net.ParseIP(r.RemoteAddr),
		},
		Name:  "name",
		Agent: "agent",
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req.Phone)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Do usecase operations
	// Initiate Authentication
	controller.log.Println("initiating authentication")
	preSession, err := controller.interactor.InitPreSession()
	if err != nil {
		controller.log.Println("failed to initiate authentication")
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Authenticate Device
	controller.log.Println("authenticating device")
	err = controller.interactor.AuthDevice(preSession.Token, req.Device.IP, req.Device.Name, req.Device.Agent)
	if err != nil {
		// Send error response
		controller.log.Printf("failed authenticating device : %s\n", err.Error())
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Init Phone Authentication
	controller.log.Println("authenticating phone")
	phoneAuth, err := controller.interactor.InitPhoneAuth(preSession.Token, req.Phone.Prefix, req.Phone.Number)
	if err != nil {
		controller.log.Printf("failed authenticating phone : %s\n", err.Error())
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Return Response
	SendJSONResponse(w, Response{
		Success: true,
		Data: struct {
			Token   string `json:"token"`
			Method  string `json:"method"`
			Length  int64  `json:"length"`
			Timeout int64  `json:"timeout"`
		}{
			Token:   preSession.Token,
			Method:  phoneAuth.Method,
			Length:  phoneAuth.Length,
			Timeout: phoneAuth.Timeout,
		},
	}, http.StatusOK)
}
