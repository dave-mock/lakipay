package socket

import (
	"log"
	"net/http"
)

type Controller struct {
	log *log.Logger
	sm  *http.ServeMux
	
}
