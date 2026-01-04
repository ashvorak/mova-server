package httpserver

import (
	"encoding/json"
	"mova-server/internal/chats"
	"mova-server/internal/users"
	"net/http"
	"strings"
	"time"
)

type CreateMessageRequest struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

type MessageResponse struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	UserID    string    `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) messagesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.messagesPostHandler(w, r)
	case http.MethodGet:
		h.messagesGetHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) messagesPostHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	// TODo: move parsing to a different layer
	chatID, err := chats.ParseID(req.ChatID)
	if err != nil {
		http.Error(w, "failed to parse chat ID", http.StatusBadRequest)
		return
	}

	userID, err := users.ParseID(req.UserID)
	if err != nil {
		http.Error(w, "failed to parse user ID", http.StatusBadRequest)
		return
	}

	m := h.messageService.Create(chatID, userID, req.Text)

	response := &MessageResponse{
		// TODO: move it to a different layer
		ID:        m.ID.String(),
		ChatID:    m.ChatID.String(),
		UserID:    m.UserID.String(),
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) messagesGetHandler(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat_id")
	if len(chatID) == 0 {
		http.Error(w, "missing chat_id query parameter", http.StatusBadRequest)
		return
	}

	chatIDParsed, err := chats.ParseID(chatID)
	if err != nil {
		http.Error(w, "invalid chat_id query parameter", http.StatusBadRequest)
		return
	}
	messages := h.messageService.ListByChat(chatIDParsed)

	responses := make([]MessageResponse, 0, len(messages))

	for _, m := range messages {
		responses = append(responses, MessageResponse{
			// TODO: move it to a different layer
			ID:        m.ID.String(),
			ChatID:    m.ChatID.String(),
			UserID:    m.UserID.String(),
			Text:      m.Text,
			CreatedAt: m.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "failed to encode responses", http.StatusInternalServerError)
	}
}
