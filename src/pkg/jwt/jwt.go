package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

/*

type Payload struct {
	Username string `json:"username"`
	Policies []interface{}
	ReqLimit int   `json:"reqLimit"`
	Iat      int64 `json:"iat"`
	Exp      int64 `json:"exp"`
}

*/

type Payload struct {
	Policies []interface{}
	ReqLimit int
	Iat      int64
	Exp      int64
	Public   interface{}
}

// Base64Encode takes in a string and returns a base 64 encoded string
func Base64Encode(src string) string {
	return strings.
		TrimRight(base64.URLEncoding.
			EncodeToString([]byte(src)), "=")
}

// Base64Encode takes in a base 64 encoded string and returns the //actual string or an error of it fails to decode the string
func Base64Decode(src string) (string, error) {
	if l := len(src) % 4; l > 0 {
		src += strings.Repeat("=", 4-l)
	}
	decoded, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		errMsg := fmt.Errorf("decoding Error %s", err)
		return "", errMsg
	}
	return string(decoded), nil
}

// Hash generates a Hmac256 hash of a string using a secret
func Hash(src string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(src))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// isValidHash validates a hash againt a value
func isValidHash(value string, hash string, secret string) bool {
	return hash == Hash(value, secret)
}

//    NB: Hash cannot be reversed all you can do is hash the same character and compare it with a hashed value. If it evaluates to true, then the character is a what is in the hash.
//   The isValidHash function only hashes the value with the secret and comared it with the hash
//   Above we created two methods, One for generating an HS256 hash and the other for validating a string against a hash.

// Encode generates a jwt.
func Encode(payload Payload, secret string) string {
	type Header struct {
		Alg string
		Typ string
	}
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}
	str, _ := json.Marshal(header)
	_header := Base64Encode(string(str))
	encodedPayload, _ := json.Marshal(payload)
	signatureValue := _header + "." +
		Base64Encode(string(encodedPayload))
	return signatureValue + "." + Hash(signatureValue, secret)
}

func Decode(jwt string, secret string) (Payload, error) {
	var payload Payload
	token := strings.Split(jwt, ".")
	// check if the jwt token contains
	// header, payload and token
	if len(token) != 3 {
		splitErr := errors.New("invalid token: token should contain header, payload and secret")
		return payload, splitErr
	}
	// decode payload
	decodedPayload, PayloadErr := Base64Decode(token[1])
	if PayloadErr != nil {
		return payload, fmt.Errorf("invalid payload: %s", PayloadErr.Error())
	}
	// parses payload from string to a struct
	ParseErr := json.Unmarshal([]byte(decodedPayload), &payload)
	if ParseErr != nil {
		return payload, fmt.Errorf("invalid payload: %s", ParseErr.Error())
	}
	// checks if the token has expired.
	if payload.Exp != 0 && time.Now().Unix() > payload.Exp {
		return payload, errors.New("expired token: token has expired")
	}
	signatureValue := token[0] + "." + token[1]
	// verifies if the header and signature is exactly whats in
	// the signature
	if !isValidHash(signatureValue, token[2], secret) {
		return payload, errors.New("invalid token")
	}
	return payload, nil
}
