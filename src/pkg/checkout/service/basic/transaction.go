package basic

import (
	accountEntity "auth/src/pkg/account/core/entity"
	"auth/src/pkg/checkout/core/entity"
	"auth/src/pkg/checkout/service/basic/processors"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
)

func (service BasicCheckoutService) InitTransaction(
	to string,
	medium string,
	amount float64,
	redirects struct {
		Success string
		Cancel  string
		Decline string
	},
	details map[string]interface{},
) (*entity.Transaction, error) {
	// Find recipient
	recepientId, err := uuid.Parse(to)
	if err != nil {
		return nil, err
	}
	accs, err := service.account.GetUserAccounts(recepientId)
	if err != nil {
		return nil, err
	}

	var acc *accountEntity.Account

	for _, v := range accs {
		if v.Type == accountEntity.STORED && v.Default {
			acc = &v
		}
	}

	if acc == nil {
		return nil, errors.New("could not find associated account for the specified recepient")
	}

	// Check gateway support
	gateway, err := service.FindGatewayByKey(medium)
	if err != nil {
		return nil, err
	}

	// if !gateway.CanProcess {
	// 	return nil, errors.New("the selected medium can not process transaction")
	// }

	// Init Transaction
	txn := entity.Transaction{
		Id:  xid.New().String(),
		For: to,
		To:  acc.Id.String(),
		Pricing: struct {
			Amount float64
			Fees   []map[string]float64
		}{
			Amount: amount,
			Fees: []map[string]float64{
				{
					"transaction": amount * (2.75 / 100),
				},
			},
		},
		Ttl:     60 * 60,
		Details: details,
		GateWay: gateway.Key,
		Status: struct {
			Value   entity.TransactionStatus
			Message string
		}{
			Value:   entity.TxnPending,
			Message: "Waiting for user's confirmation",
		},
		CreatedAt: time.Now(),
	}

	// Store Transaction
	err = service.repo.StoreTransaction(txn)
	if err != nil {
		return &txn, err
	}

	// Return

	return &txn, nil
}
func (service BasicCheckoutService) ConfirmTransaction(id string) (any, error) {

	// Find transaction
	txn, err := service.repo.FindTransaction(id)
	if err != nil {
		return nil, err
	}

	var data any
	var res map[string]interface{}

	switch txn.GateWay {
	case "CYBERSOURCE":
		{
			// Process Cybersource
			data, err = processors.ProcessCybersource(txn.Id, txn.Pricing.Amount, "https://api.lakipay.co:32000")
			if err != nil {
				return nil, err
			}
			txn.Status = struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value:   entity.TxnProcessing,
				Message: "Waiting for payment completion",
			}
			err = service.repo.UpdateTransaction(*txn)
			if err != nil {
				return nil, err
			}

			res = map[string]interface{}{
				"type": "WEB",
				"data": data,
			}

		}
	case "CBE_BIRR":
		{
			err = processors.ProcessCBRBirr(txn.Id, txn.Pricing.Amount+txn.Pricing.Fees[0]["transaction"], txn.Details["phone"].(string))
			if err != nil {
				return nil, err
			}

			txn.Status = struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value:   entity.TxnProcessing,
				Message: "Waiting for payment completion",
			}

			err = service.repo.UpdateTransaction(*txn)
			if err != nil {
				return nil, err
			}

			res = map[string]interface{}{
				"type": "USSD",
			}
		}
	case "TELEBIRR":
		{
			err = processors.ProcessTeleBirr(txn.Id, txn.Pricing.Amount+txn.Pricing.Fees[0]["transaction"], txn.Details["phone"].(string))
			if err != nil {
				return nil, err
			}

			txn.Status = struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value:   entity.TxnProcessing,
				Message: "Waiting for payment completion",
			}

			err = service.repo.UpdateTransaction(*txn)
			if err != nil {
				return nil, err
			}

			res = map[string]interface{}{
				"type": "USSD",
			}
		}
	case "M_PESA":
		{
			err = processors.ProcessMPesa(txn.Id, txn.Pricing.Amount+txn.Pricing.Fees[0]["transaction"], txn.Details["phone"].(string))
			if err != nil {
				return nil, err
			}

			txn.Status = struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value:   entity.TxnProcessing,
				Message: "Waiting for payment completion",
			}

			err = service.repo.UpdateTransaction(*txn)
			if err != nil {
				return nil, err
			}

			res = map[string]interface{}{
				"type": "USSD",
			}
		}
	}

	return res, nil
}

func (service BasicCheckoutService) GetTransaction(id string) (entity.Transaction, error) {
	txn, err := service.repo.FindTransaction(id)

	return *txn, err
}

func (service BasicCheckoutService) UpdatePaymentStatus(id string, status struct {
	Value   entity.TransactionStatus
	Message string
}) error {
	// Find transaction
	txn, err := service.repo.FindTransaction(id)
	if err != nil {
		return err
	}

	txn.Status = status
	err = service.repo.UpdateTransaction(*txn)
	if err != nil {
		return err
	}

	if txn.NotifyUrl != "" {
		// Call webhook
		client := http.Client{
			Timeout: 3 * time.Second,
		}

		serBody, err := json.Marshal(txn)
		if err != nil {
			return err
		}

		req, _ := http.NewRequest(http.MethodPost, txn.NotifyUrl, bytes.NewBuffer(serBody))
		req.Header.Add("Content-Type", "application/json")

		go func() {
			client.Do(req)
		}()
	}

	return nil
}

func (service BasicCheckoutService) InitDirectTransaction(to string, medium string, amount float64, details map[string]interface{}, callback string) (any, error) {

	// Find recipient
	recepientId, err := uuid.Parse(to)
	if err != nil {
		return nil, err
	}
	accs, err := service.account.GetUserAccounts(recepientId)
	if err != nil {
		return nil, err
	}

	var acc *accountEntity.Account

	for _, v := range accs {
		if v.Type == accountEntity.STORED && v.Default {
			acc = &v
		}
	}

	if acc == nil {
		return nil, errors.New("could not find associated account for the specified recepient")
	}

	// Check gateway support
	gateway, err := service.FindGatewayByKey(medium)
	if err != nil {
		return nil, err
	}

	// if !gateway.CanProcess {
	// 	return nil, errors.New("the selected medium can not process transaction")
	// }

	// Init Transaction
	txn := entity.Transaction{
		Id:  xid.New().String(),
		For: to,
		To:  acc.Id.String(),
		Pricing: struct {
			Amount float64
			Fees   []map[string]float64
		}{
			Amount: amount,
			Fees: []map[string]float64{
				{
					"transaction": amount * (2.75 / 100),
				},
			},
		},
		NotifyUrl: callback,
		Ttl:       60 * 60,
		Details:   details,
		GateWay:   gateway.Key,
		Status: struct {
			Value   entity.TransactionStatus
			Message string
		}{
			Value:   entity.TxnPending,
			Message: "Waiting for user's confirmation",
		},
		CreatedAt: time.Now(),
	}

	var data any
	var res map[string]interface{}

	switch txn.GateWay {
	case "CYBERSOURCE":
		{
			// Process Cybersource
			data, err = processors.ProcessCybersource(txn.Id, txn.Pricing.Amount, "https://api.lakipay.co:32000")
			if err != nil {
				return nil, err
			}
			txn.Status = struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value:   entity.TxnProcessing,
				Message: "Waiting for payment completion",
			}

			res = map[string]interface{}{
				"type": "WEB",
				"data": data,
			}

		}
	case "CBE_BIRR":
		{
			err = processors.ProcessCBRBirr(txn.Id, txn.Pricing.Amount+txn.Pricing.Fees[0]["transaction"], txn.Details["phone"].(string))
			if err != nil {
				return nil, err
			}

			txn.Status = struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value:   entity.TxnProcessing,
				Message: "Waiting for payment completion",
			}

			res = map[string]interface{}{
				"type": "USSD",
			}
		}
	}

	err = service.repo.StoreTransaction(txn)
	if err != nil {
		return nil, err
	}

	return res, nil
}
