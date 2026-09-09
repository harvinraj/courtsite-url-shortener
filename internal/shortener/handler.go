package shortener

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type AnalyticRequest struct {
	ShortKey string `json:"short_key"`
}

func (h *Handler) HandleShortener(response http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{Error: "Invalid request body"})
		fmt.Println("Request body doest not contain url field : ", err)
		return
	}

	saved, err := h.service.ShortenUrl(request.Context(), req.URL)
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{Error: err.Error()})
		fmt.Println("Error saving url : ", err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusAccepted)
	json.NewEncoder(response).Encode(saved)
	fmt.Println("url saved : " + saved.OriginalURL + " shortKey :" + saved.ShortKey)
}

func (h *Handler) HandleAnalytics(response http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyticRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{Error: "Invalid request body"})
		fmt.Println("Request body doest not contain url field : ", err)
		return
	}

	found, err := h.service.GetShortenURL(request.Context(), req.ShortKey)
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{Error: err.Error()})
		fmt.Println("Error saving url : ", err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusAccepted)
	json.NewEncoder(response).Encode(found)
	fmt.Println("original url :" + found.OriginalURL)
	fmt.Printf("visits : %d", found.Visits)
	fmt.Println("")

}
