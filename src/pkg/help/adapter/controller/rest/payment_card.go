package rest

import (
	"log"
	"net/http"
)

func (controller Controller) GetAddPaymentCard(w http.ResponseWriter, r *http.Request) {
	log.Println("Payment Cards POST")
	controller.log.Println(r.Body)
	w.Write([]byte("Success"))
}
