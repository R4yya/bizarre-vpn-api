package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterPingRoute(router *gin.Engine) {
	router.GET("/ping", handlers.PingHandler)
}
