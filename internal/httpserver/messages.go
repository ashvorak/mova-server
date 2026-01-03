package httpserver

import (
	"encoding/json"
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

	m := h.messageService.Create(req.ChatID, req.UserID, req.Text)

	response := &MessageResponse{
		ID:        m.ID,
		ChatID:    m.ChatID,
		UserID:    m.UserID,
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

	messages := h.messageService.ListByChat(chatID)

	responses := make([]MessageResponse, 0, len(messages))

	for _, m := range messages {
		responses = append(responses, MessageResponse{
			ID:        m.ID,
			ChatID:    m.ChatID,
			UserID:    m.UserID,
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
