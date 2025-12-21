package httpserver

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func New() *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/api/health", healthHandler)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return s
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
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
