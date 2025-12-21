package main

import (
	"log"
	"mova-server/internal/httpserver"
)

func main() {
	server := httpserver.New()
	log.Fatal(server.ListenAndServe())
}
