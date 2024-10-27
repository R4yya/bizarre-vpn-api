package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterSubscriptionPlanRoutes(router *gin.Engine) {
	subscriptionGroup := router.Group("/plans")

	subscriptionGroup.POST("/", handlers.CreatePlanHandler)
	subscriptionGroup.GET("/", handlers.GetAllPlansHandler)
	subscriptionGroup.GET("/:id", handlers.GetPlanHandler)
	subscriptionGroup.PUT("/:id", handlers.UpdatePlanHandler)
	subscriptionGroup.DELETE("/:id", handlers.DeletePlanHandler)
}
