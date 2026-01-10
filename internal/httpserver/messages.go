package httpserver

import (
	"encoding/json"
	"mova-server/internal/chats"
	"mova-server/internal/messages"
	"mova-server/internal/users"
	"net/http"
	"strconv"
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
	chatIDStr := r.URL.Query().Get("chat_id")
	if len(chatIDStr) == 0 {
		http.Error(w, "missing chat_id query parameter", http.StatusBadRequest)
		return
	}

	chatID, err := chats.ParseID(chatIDStr)
	if err != nil {
		http.Error(w, "invalid chat_id query parameter", http.StatusBadRequest)
		return
	}

	var afterID messages.ID
	afterIDStr := r.URL.Query().Get("after")
	if len(afterIDStr) != 0 {
		afterID, err = messages.ParseID(afterIDStr)
		if err != nil {
			http.Error(w, "invalid after parameter", http.StatusBadRequest)
			return
		}
	}

	limitStr := r.URL.Query().Get("limit")
	var limit int
	if len(limitStr) != 0 {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	msgs := h.messageService.ListByChatAfter(chatID, afterID, limit)

	responses := make([]MessageResponse, 0, len(msgs))

	for _, m := range msgs {
		responses = append(responses, MessageResponse{
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
