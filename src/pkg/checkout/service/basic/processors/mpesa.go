package processors

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ProcessMPesa(id string, amount float64, phone string) error {
	url := "https://"
	method := http.MethodPost

	pld := map[string]interface{}{}

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
