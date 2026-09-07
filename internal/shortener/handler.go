package shortener

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) HandleShorten(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	saved, err := h.service.Shorten(request.Context(), req.URL)
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusAccepted)
	json.NewEncoder(response).Encode(saved)
	fmt.Printf("url saved and visits updated, original url = %s", saved.OriginalURL)

}

func (h *Handler) HandleResolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract key from path (e.g., /r/{key} or /{key})
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "missing short key", http.StatusBadRequest)
		return
	}
	key := parts[0]

	record, err := h.service.Retrieve(r.Context(), key)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "short url not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(record)
}
