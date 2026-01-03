package httpserver

import (
	"mova-server/internal/chats"
	"mova-server/internal/messages"
	"mova-server/internal/users"
)

type Handler struct {
	userService    *users.Service
	chatService    *chats.Service
	messageService *messages.Service
}

func NewHandler(userService *users.Service, chatService *chats.Service, messageService *messages.Service) *Handler {
	h := &Handler{
		userService:    userService,
		chatService:    chatService,
		messageService: messageService,
	}

	return h
}
