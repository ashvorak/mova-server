package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"
)

type CreateChatRequest struct {
	UserIDs []string `json:"user_ids"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	UserIDs []string `json:"user_ids"`
}

func (h *Handler) chatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.chatsPostHandler(w, r)
	} else if r.Method == http.MethodGet {
		h.chatsGetHandler(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) chatsPostHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	chat := h.chatService.Create(req.UserIDs)

	response := &ChatResponse{
		ID:      chat.ID,
		UserIDs: chat.UserIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) chatsGetHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if len(userID) == 0 {
		http.Error(w, "missing user_id query parameter", http.StatusBadRequest)
		return
	}

	chats := h.chatService.ListByUser(userID)

	chatResponses := make([]ChatResponse, 0, len(chats))

	for _, c := range chats {
		chatResponses = append(chatResponses, ChatResponse{
			ID:      c.ID,
			UserIDs: c.UserIDs,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(chatResponses); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
