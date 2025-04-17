package processors

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ProcessTeleBirr(id string, amount float64, phone string) error {

	url := "https://api.lakipay.co/account/telebirr/ussd-push"
	method := http.MethodPost

	pld := map[string]interface{}{
		"CommandID":                "InitTrans_BuyGoodsForCustomer",
		"OriginatorConversationID": id,
		"ThirdPartyID":             "Laki-Pay",
		"Password":                 "jBq7JfxTs0C5ji0VPKakmRSgBbeh4NO0juJ1LXnPIOw=",
		"ResultURL":                "https://api.lakipay.co/api/v1/checkout/transactions/notify",
		"Timestamp":                "20150101010101",
		"IdentifierType":           12,
		"Identifier":               "51437701",
		"SecurityCredential":       "SvO8Px+vMnrGxTnZFw5tjUopf08sr1GZxs7qE3lmzG0=",
		"ShortCode":                "514377",
		"PrimaryParty":             phone,
		"ReceiverParty":            "514377",
		"Amount":                   amount,
		"Currency":                 "ETB",
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
