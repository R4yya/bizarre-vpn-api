package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/api/middlewares"
	"bizarre-vpn-api/internal/services/authLinkService"
	"bizarre-vpn-api/internal/services/userService"
)

func UserRoutes(customRouter *CustomRouter) {
	userService := userService.NewUserService(
		customRouter.log,
		customRouter.storage.UserStorage,
	)

	authLinkService := authLinkService.NewAuthLinksService(
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

	customRouter.routerGroup.GET("",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		userHandler.GetUserDataHandler,
	)

	customRouter.routerGroup.GET("/list",
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
