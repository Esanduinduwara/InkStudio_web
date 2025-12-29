package handlers

import (
	"net/http"
	"time"

	"inkstudio-backend/pkg/response"
)

// Health returns the health status of the server
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}, http.StatusOK)
}
