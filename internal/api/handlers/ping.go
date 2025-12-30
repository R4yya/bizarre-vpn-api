package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler responds to requests to /ping
// @Summary Checking server availability
// @Description Returns "pong" to check if the API is available
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} MessageResponse
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, MessageResponse{Message: "pong"})
}
