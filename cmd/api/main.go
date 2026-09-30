package main

import (
	handler "courtsite-url-shortener/internal/handler"
	memory "courtsite-url-shortener/internal/repository/memory"
	service "courtsite-url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize dependencies (DI)
	store := memory.NewMemoryStore()
	svc := service.NewService(store)
	h := handler.NewHTTPHandler(svc)

	r := gin.Default()

	r.POST("/shorten", h.HandleShortener)
	r.POST("/analytics", h.HandleAnalytic)

	r.Run(":8080")

}
