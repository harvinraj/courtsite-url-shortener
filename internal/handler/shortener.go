package handler

import (
	"courtsite-url-shortener/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	Url string `json:"url"`
}

type ShortenResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortKey    string `json:"short_key"`
}

type AnalyticRequest struct {
	ShortKey string `json:"short_key"`
}

type AnalyticResponse struct {
	OriginalUrl string `json:"original_url"`
	Visits      int64  `json:"visits"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type HTTPHandler struct {
	svc *service.ShortenerService
}

func NewHTTPHandler(svc *service.ShortenerService) *HTTPHandler {
	return &HTTPHandler{
		svc: svc,
	}
}

func (h *HTTPHandler) HandleShortener(c *gin.Context) {
	var req ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	savedUrl, savedKey, err := h.svc.ShortenUrl(c.Request.Context(), req.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	c.JSON(http.StatusAccepted, ShortenResponse{
		OriginalUrl: savedUrl.OriginalURL,
		ShortKey:    savedKey.ShortKey,
	})
}

func (h *HTTPHandler) HandleAnalytic(c *gin.Context) {

	var req AnalyticRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	record, err := h.svc.GetShortenUrl(c.Request.Context(), req.ShortKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	c.JSON(http.StatusOK, AnalyticResponse{
		OriginalUrl: record.OriginalURL,
		Visits:      record.Visits,
	})

}
