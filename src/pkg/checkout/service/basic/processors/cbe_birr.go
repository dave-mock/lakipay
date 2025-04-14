package processors

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ProcessCBRBirr(id string, amount float64, phone string) error {
	url := "http://196.190.251.169:33180/api/v1/payments/create"
	method := http.MethodPost

	pld := map[string]interface{}{
		"amount":        amount,
		"description":   "Payment for service",
		"referenceId":   id,
		"callbackUrl":   "https://api.lakipay.co:32000/api/v1/checkout/transactions/notify",
		"phoneNumber":   phone,
		"merchantId":    "TsegawTest",
		"merchantKey":   "piHTZDyB1jh8OK1Jb04EcA==",
		"terminalId":    "202526",
		"credentialKey": "A49euY3AvwQCPt8FOKsOSQ==",
		"userId":        "LakiOperator",
	}

	serPld, err := json.Marshal(pld)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(serPld)))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
