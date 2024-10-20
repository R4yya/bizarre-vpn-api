package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/user")

	userGroup.POST("/auth", handlers.AuthorizeUserHandler)
}
