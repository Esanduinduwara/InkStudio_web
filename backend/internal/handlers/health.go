package handlers

import (
	"net/http"

	"inkstudio-backend/pkg/response"
)

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check database connection
	if err := h.DB.Ping(); err != nil {
		response.JSON(w, map[string]string{
			"status":   "unhealthy",
			"database": "disconnected",
		}, http.StatusServiceUnavailable)
		return
	}

	response.JSON(w, map[string]string{
		"status":   "healthy",
		"database": "connected",
	}, http.StatusOK)
}
