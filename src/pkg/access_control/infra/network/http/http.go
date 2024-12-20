package http

import (
	"context"
	// "io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	log      *log.Logger
	ServeMux *http.ServeMux
	Server   *http.Server
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func New(log *log.Logger) Server {
	sm := &http.ServeMux{}

	s := &http.Server{
		Handler: accessControl(sm),
		Addr:    "0.0.0.0:3000",
	}

	// sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("olla222"))
	// })
	// sm.HandleFunc("/epg", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println(io.ReadAll(r.Body))
	// 	w.Write([]byte("olla"))
	// })

	return Server{log: log, ServeMux: sm, Server: s}
}

func (s Server) Serve() {
	go func() {
		err := s.Server.ListenAndServe()
		if err != nil {
			// shut down
			log.Println(err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	s.Server.Shutdown(ctx)
}
