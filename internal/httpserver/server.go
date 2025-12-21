package httpserver

import (
	"fmt"
	"net/http"
)

func New() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return s
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
