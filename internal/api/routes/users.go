package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/api/middlewares"
)

func UserRoutes(customRouter *CustomRouter) {
	userHandler := handlers.UserHandler{
		Log:         customRouter.log,
		UserStorage: customRouter.storage.UserStorage,
	}

	customRouter.AddGroup("/auth", AuthRoutes)

	customRouter.routerGroup.GET("/", middlewares.AuthRequired(customRouter.log, customRouter.cfg), userHandler.GetUserInfoHandler)
	// customRouter.routerGroup.POST("/auth", userHandler.AuthorizeUserHandler)
}
