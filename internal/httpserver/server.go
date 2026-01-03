package httpserver

import (
	"mova-server/internal/chats"
	"mova-server/internal/messages"
	"mova-server/internal/users"
	"net/http"
)

func New(userService *users.Service, chatService *chats.Service, messageService *messages.Service) *http.Server {
	h := NewHandler(userService, chatService, messageService)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", h.pingHandler)
	mux.HandleFunc("/api/health", h.healthHandler)
	mux.HandleFunc("/api/users", h.usersHandler)
	mux.HandleFunc("/api/chats", h.chatsHandler)
	mux.HandleFunc("/api/messages", h.messagesHandler)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
