package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adrinicomartin/keystore-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AHBPayload struct {
	Sub string  `json:"sub"`
	ISS string  `json:"iss"`
	Aud string  `json:"aud"`
	Exp float64 `json:"exp"`
	Iat float64 `json:"iat"`
	Jti string  `json:"jti"`
}

var jwtToken string
var expTime int = 30
var nowMillis int64 = time.Now().Unix()
var expMillis int64 = (time.Now().Unix()) + int64(expTime*60*1000)

func main() {
	ks := readKeyStore("./IntellectARXSJWT.jks", []byte("Intellect01"))

	var key = ks["intellect"].(*keystore.PrivateKeyEntry).PrivKey

	var t = jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"sub": "intellect",
			"iss": "intellect",
			"aud": "default",
			"exp": expMillis / 1000,
			"iat": nowMillis / 1000,
			"jti": uuid.NewString(),
		})

	t.Header = map[string]interface{}{
		"alg": "RS256",
		"typ": "JWT",
	}

	s, err := t.SignedString(key)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("TOKEN")
	fmt.Println(s)

}

func readKeyStore(filename string, password []byte) keystore.KeyStore {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	keyStore, err := keystore.Decode(f, password)
	if err != nil {
		log.Fatal(err)
	}
	return keyStore
}

// getPrivateKey returns a private key from a key store file

// func getPrivateKey() (privateKey crypto.PrivateKey, err error) {
// 	// initialize the private key to nil
// 	privateKey = nil
// 	// define the key store path, password, and alias
// 	keyStorePath := "./IntellectARXSJWT.jks"
// 	keyPassword := "Intellect01"
// 	keyAlias := "intellect"

// 	// open the key store file
// 	file, err := os.Open(keyStorePath)
// 	if err != nil {
// 		// return the error if the file cannot be opened
// 		return
// 	}
// 	// defer closing the file
// 	defer file.Close()
// 	// create a key store object
// 	ks := keystore.New()
// 	// load the key store from the file
// 	err = ks.Load(file, []byte(keyPassword))
// 	if err != nil {
// 		// return the error if the key store cannot be loaded
// 		return
// 	}
// 	// get the key entry from the key store by the alias and the password
// 	keyEntry, ok := ks.GetPrivateKeyEntry(keyAlias, []byte(keyPassword))
// 	if !ok {
// 		// return an error if the key entry cannot be found
// 		err = errors.New("key entry not found")
// 		return
// 	}
// 	// get the private key from the key entry
// 	privateKey = keyEntry.PrivKey
// 	// return the private key and nil error
// 	return
// }
