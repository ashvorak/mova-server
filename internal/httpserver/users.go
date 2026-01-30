package httpserver

import (
	"encoding/json"
	"net/http"
)

type CreateUserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.usersPostHandler(w, r)
	case http.MethodGet:
		h.usersGetHandler(w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) usersPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Create(req.Name)
	if err != nil {
		http.Error(w, "failed to create user: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := &UserResponse{
		ID:   user.ID.String(),
		Name: user.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) usersGetHandler(w http.ResponseWriter) {
	users := h.userService.List()

	userResponses := make([]UserResponse, 0, len(users))

	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:   u.ID.String(),
			Name: u.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(userResponses); err != nil {
		http.Error(w, "failed to encode responce", http.StatusInternalServerError)
	}
}
