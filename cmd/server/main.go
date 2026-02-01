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

	repo := messages.NewMemoryRepository()
	messageService := messages.NewService(repo)

	server := httpserver.New(userService, chatService, messageService)
	log.Fatal(server.ListenAndServe())
}
