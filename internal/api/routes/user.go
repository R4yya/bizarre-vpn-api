package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
)

func UserRoutes(customRouter *CustomRouter) {
	userHandler := handlers.UserHandler{Log: customRouter.log}

	customRouter.routerGroup.POST("/auth", userHandler.AuthorizeUserHandler)
}
