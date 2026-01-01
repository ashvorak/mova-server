package main

import (
	"log"
	"mova-server/internal/chats"
	"mova-server/internal/httpserver"
	"mova-server/internal/users"
)

func main() {
	userService := users.NewService()
	chatService := chats.NewService()

	server := httpserver.New(userService, chatService)
	log.Fatal(server.ListenAndServe())
}
