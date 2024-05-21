package rest

// func (controller Controller) GetVerifyBankAccount(w http.ResponseWriter, r *http.Request) {

// 	// Request
// 	type Request struct {
// 		BankId        string `json:"bank_id"`
// 		AccountNumber string `json:"account_number"`
// 		Holder        struct {
// 			Name  string `json:"name"`
// 			Phone string `json:"phone"`
// 		} `json:"holder"`
// 	}

// 	var req Request

// 	// Decoder
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&req)
// 	if err != nil {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    "INVALID_REQUEST",
// 				Message: "The request contains one or more invalid parameters please refer to the spec",
// 			},
// 		}, http.StatusBadRequest)
// 		return
// 	}

// 	defer r.Body.Close()

// 	bankId, _ := uuid.Parse(req.BankId)
// 	// Usecase
// 	err = controller.interactor.VerifyBankAccount(bankId, req.AccountNumber, req.Holder.Name, req.Holder.Phone)
// 	if err != nil {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    err.(usecase.Error).Type,
// 				Message: err.(usecase.Error).Message,
// 			},
// 		}, http.StatusBadRequest)
// 		return
// 	}

// 	w.Write([]byte("Verified"))
// }
