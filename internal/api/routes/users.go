package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/api/middlewares"
	"bizarre-vpn-api/internal/services"
)

func UserRoutes(customRouter *CustomRouter) {
	userService := services.NewUserService(
		customRouter.log,
		customRouter.storage.UserStorage,
	)

	authLinkService := services.NewAuthLinksService(
		customRouter.log,
		customRouter.storage.AuthLinksStorage,
		customRouter.storage.LnkUserProviderStorage,
	)

	userHandler := handlers.UserHandler{
		Log:             customRouter.log,
		BotSharedData:   customRouter.botSharedData,
		UserService:     userService,
		AuthLinkService: authLinkService,
	}

	customRouter.AddGroup("/auth", AuthRoutes)

	customRouter.routerGroup.GET("/self",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		userHandler.GetUserDataHandler,
	)

	customRouter.routerGroup.GET("/",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
		userHandler.GetUsersListHandler,
	)

	customRouter.routerGroup.GET("/:id/auth-link",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
		userHandler.CreateUserAuthLink,
	)

	customRouter.routerGroup.POST("/",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
		userHandler.CreateUser,
	)
}
