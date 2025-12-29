package httpserver

import (
	"encoding/json"
	"mova-server/internal/users"
	"net/http"
)

type Handler struct {
	userService *users.Service
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

func New(userService *users.Service) *http.Server {
	h := &Handler{
		userService: userService,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", h.pingHandler)
	mux.HandleFunc("/api/health", h.healthHandler)
	mux.HandleFunc("/api/users", h.usersHandler)

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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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
