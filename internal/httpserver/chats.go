package httpserver

import (
	"encoding/json"
	"mova-server/internal/users"
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
	switch r.Method {
	case http.MethodPost:
		h.chatsPostHandler(w, r)
	case http.MethodGet:
		h.chatsGetHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
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

	var userIDs []users.ID
	for _, uidStr := range req.UserIDs {
		uid, err := users.ParseID(uidStr)
		if err != nil {
			http.Error(w, "failed to parse user ID", http.StatusBadRequest)
			return
		}
		userIDs = append(userIDs, uid)
	}
	chat := h.chatService.Create(userIDs)

	var userIDsStr []string
	for _, uid := range chat.UserIDs {
		userIDsStr = append(userIDsStr, uid.String())
	}

	response := &ChatResponse{
		ID:      chat.ID.String(),
		UserIDs: userIDsStr,
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
		var userIDs []string
		for _, uid := range c.UserIDs {
			userIDs = append(userIDs, uid.String())
		}

		chatResponses = append(chatResponses, ChatResponse{
			ID:      c.ID.String(),
			UserIDs: userIDs,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(chatResponses); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
