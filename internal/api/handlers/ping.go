package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PingResponse struct {
	Message string `json:"message"`
}

// PingHandler responds to requests to /ping
// @Summary Checking server availability
// @Description Returns "pong" to check if the API is available
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, PingResponse{Message: "pong"})
}
