package processors

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

func sign(fields map[string]string, secretKey string) string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var encodedFields []string
	for _, k := range keys {
		encodedFields = append(encodedFields, k+"="+fields[k])
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(strings.Join(encodedFields, ",")))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

func ProcessCybersource(id string, amount float64, host string) (string, error) {
	var url string
	var err error

	UNSIGNED_FIELD_NAMES := []string{}
	SIGNED_FIELD_NAMES := []string{
		"access_key",
		"amount",
		"currency",
		"locale",
		"payment_method",
		"profile_id",
		"reference_number",
		"signed_date_time",
		"signed_field_names",
		"transaction_type",
		"transaction_uuid",
		"unsigned_field_names",
	}
	DEVICE_FINGERPRINT := uuid.New()

	reqParams := map[string]string{
		"access_key":           "66ad734a971a3f79b84f183c4e52b790",
		"amount":               fmt.Sprintf("%v", amount),
		"currency":             "ETB",
		"locale":               "en-US",
		"payment_method":       "card",
		"profile_id":           "2674CE2D-EA15-4D9F-85D0-713FC5F6329F",
		"reference_number":     DEVICE_FINGERPRINT.String(),
		"signed_date_time":     time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"signed_field_names":   strings.Join(SIGNED_FIELD_NAMES, ","),
		"transaction_type":     "sale",
		"transaction_uuid":     id,
		"unsigned_field_names": strings.Join(UNSIGNED_FIELD_NAMES, ","),
	}

	walletId := id
	f, err := os.Create(fmt.Sprintf("./public/%s.html", walletId))
	if err := os.WriteFile(
		fmt.Sprintf("./public/%s.html", walletId),
		[]byte(fmt.Sprintf(`
								<!DOCTYPE html>
								<html lang="en">
									<head>
										<meta charset="UTF-8">
										<!-- <meta name="viewport" content="width=vw, initial-scale=1.0"> -->
										<title>Cybersource Payment</title>
										<script src="https://cdnjs.cloudflare.com/ajax/libs/crypto-js/4.0.0/crypto-js.min.js"></script>
									</head>
									<body>
										<form id='payment_form' method='post' action="https://secureacceptance.cybersource.com/pay">
											<input type='hidden' id='access_key' name='access_key' />
											<input type='hidden' id='locale' name='locale' />
											<input type='hidden' id='payment_method' name='payment_method' />
											<input type='hidden' id='profile_id' name='profile_id' />
											<input type='hidden' id='signature' name='signature' />
											<input type='hidden' id='signed_date_time' name='signed_date_time' />
											<input type='hidden' id='signed_field_names' name='signed_field_names' />
											<input type='hidden' id='transaction_type' name='transaction_type' />
											<input type='hidden' id='transaction_uuid' name='transaction_uuid' />
											<input type='hidden' id='unsigned_field_names' name='unsigned_field_names' />
											<input type='hidden' id='reference_number' name='reference_number' />
											<input type='hidden' id='amount' name='amount' />
											<input type='hidden' id='currency' name='currency' />
										</form>
									</body>
									<script type="text/javascript">
										function load() {
											document.getElementById("access_key").value = "%[1]s";
											document.getElementById("locale").value = "%[2]s";
											document.getElementById("payment_method").value = "%[3]s";
											document.getElementById("profile_id").value = "%[4]s";
											document.getElementById("signature").value = "%[5]s";
											document.getElementById("signed_date_time").value = "%[6]s";
											document.getElementById("signed_field_names").value = "%[7]s";
											document.getElementById("transaction_type").value = "%[8]s";
											document.getElementById("transaction_uuid").value = "%[9]s";
											document.getElementById("unsigned_field_names").value = "%[10]s";
											document.getElementById("reference_number").value = "%[11]s";
											document.getElementById("amount").value = "%[12]s";
											document.getElementById("currency").value = "%[13]s";
											document.getElementById("payment_form").submit();
										}
										window.onload = load;
									</script>
								</html>
								`,
			reqParams["access_key"],
			reqParams["locale"],
			reqParams["payment_method"],
			reqParams["profile_id"],
			sign(reqParams, "328be89eb0ca4c53845594974c09a17cd4aa0c561bc147f0b12f6fd612cb85a21ddd85b92e0b4ff3b39153c2ad904c9966375e00572f450ba68b19a72c2055a11ea890496b9a4eaab8748fc93d7bef65a37c01d94d6d43c2ab8e7e90314ce098c6978a1d2ceb4e17b0a995d46f90099676bb45f923e64258b7d6d856e00487cc"),
			reqParams["signed_date_time"],
			reqParams["signed_field_names"],
			reqParams["transaction_type"],
			reqParams["transaction_uuid"],
			reqParams["unsigned_field_names"],
			reqParams["reference_number"],
			reqParams["amount"],
			reqParams["currency"],
		)),
		0666,
	); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	url = fmt.Sprintf("%s/static/%s.html", host, walletId)

	return url, err
}
