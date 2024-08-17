package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthCheck responds with the status of the application
// @Summary Check the service health
// @Description Check the service health
// @Tags service
// @Accept  json
// @Produce  json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status: "UP",
	}
	c.JSON(http.StatusOK, response)
}
