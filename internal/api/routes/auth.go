package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
)

func AuthRoutes(customRouter *CustomRouter) {
	authHandler := handlers.AuthHandler{
		Log:                    customRouter.log,
		CFG:                    customRouter.cfg,
		UserStorage:            customRouter.storage.UserStorage,
		LnkUserProviderStorage: customRouter.storage.LnkUserProviderStorage,
	}

	customRouter.routerGroup.POST("/refresh-tokens", authHandler.RefreshTokens)
}
