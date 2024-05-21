package rest

import (
	"net/http"
	"time"
)

func (controller Controller) Upload(w http.ResponseWriter, r *http.Request) {

	// AuthN

	// AuthZ

	// Parse
	err := r.ParseMultipartForm(0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}

	defer r.Body.Close()

	file, fh, err := r.FormFile("file")
	controller.log.Println(time.Now())
	controller.log.Println(file)
	controller.log.Println(fh.Filename)
	controller.log.Println(fh.Header)
	controller.log.Println(fh.Size)
	controller.log.Println(err)

	w.Write([]byte("Welcome to storage service"))
}
