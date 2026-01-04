package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/services"
)

func SubscriptionPlanRoutes(customRouter *CustomRouter) {
	subscriptionPlanService := services.NewSubscriptionPlanService(
		customRouter.log,
		customRouter.storage.SubscriptionPlanStorage,
	)

	subscriptionPlanHandler := handlers.SubscriptionPlanHandler{
		Log:                     customRouter.log,
		SubscriptionPlanService: subscriptionPlanService,
	}

	customRouter.routerGroup.GET("/", subscriptionPlanHandler.GetAllPlansHandler)
	customRouter.routerGroup.GET("/:id", subscriptionPlanHandler.GetPlanHandler)
	customRouter.routerGroup.POST("/", subscriptionPlanHandler.CreatePlanHandler)
	customRouter.routerGroup.PUT("/:id", subscriptionPlanHandler.UpdatePlanHandler)
	customRouter.routerGroup.DELETE("/:id", subscriptionPlanHandler.DeletePlanHandler)
}
