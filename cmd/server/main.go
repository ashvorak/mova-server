package main

import (
	"log"
	"mova-server/internal/chats"
	"mova-server/internal/httpserver"
	"mova-server/internal/messages"
	"mova-server/internal/users"
)

func main() {
	userService := users.NewService()
	chatService := chats.NewService()
	messageService := messages.NewService()

	server := httpserver.New(userService, chatService, messageService)
	log.Fatal(server.ListenAndServe())
}
