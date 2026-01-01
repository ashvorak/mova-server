package httpserver

import (
	"encoding/json"
	"mova-server/internal/chats"
	"mova-server/internal/users"
	"net/http"
)

type Handler struct {
	userService *users.Service
	chatService *chats.Service
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateChatRequest struct {
	UserIDs []string `json:"user_ids"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	UserIDs []string `json:"user_ids"`
}

func New(userService *users.Service, chatService *chats.Service) *http.Server {
	h := &Handler{
		userService: userService,
		chatService: chatService,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", h.pingHandler)
	mux.HandleFunc("/api/health", h.healthHandler)
	mux.HandleFunc("/api/users", h.usersHandler)
	mux.HandleFunc("/api/chats", h.chatsHandler)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func (s *Handler) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (s *Handler) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response := &HealthResponse{
		Status:  "ok",
		Service: "mova-server",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Handler) usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.usersPostHandler(w, r)
	} else if r.Method == http.MethodGet {
		s.usersGetHandler(w)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (s *Handler) usersPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	user := s.userService.Create(req.Name)

	response := &UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Handler) usersGetHandler(w http.ResponseWriter) {
	users := s.userService.List()

	userResponses := make([]UserResponse, 0, len(users))

	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:   u.ID,
			Name: u.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(userResponses); err != nil {
		http.Error(w, "failed to encode responce", http.StatusInternalServerError)
	}
}

func (s *Handler) chatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.chatsPostHandler(w, r)
	} else if r.Method == http.MethodGet {
		s.chatsGetHandler(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (s *Handler) chatsPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	chat := s.chatService.Create(req.UserIDs)

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

func (s *Handler) chatsGetHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if len(userID) == 0 {
		http.Error(w, "missing user_id query parameter", http.StatusBadRequest)
		return
	}

	chats := s.chatService.ListByUser(userID)

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
		http.Error(w, "failed to encode responce", http.StatusInternalServerError)
	}
}
