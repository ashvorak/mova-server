package httpserver

import (
	"mova-server/internal/chats"
	"mova-server/internal/users"
)

type Handler struct {
	userService *users.Service
	chatService *chats.Service
}

func NewHandler(userService *users.Service, chatService *chats.Service) *Handler {
	h := &Handler{
		userService: userService,
		chatService: chatService,
	}

	return h
}
