package main

import (
	"log"
	"mova-server/internal/httpserver"
	"mova-server/internal/users"
)

func main() {
	userService := users.NewService()
	server := httpserver.New(userService)
	log.Fatal(server.ListenAndServe())
}
