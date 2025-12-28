package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/services"
)

func VpnServersRoutes(customRouter *CustomRouter) {
	vpnServersService := services.NewVpnServersService(
		customRouter.log,
		customRouter.storage.VpnServersStorage,
	)

	vpnServersHandler := handlers.VpnServersHandler{
		Log:               customRouter.log,
		VpnServersService: vpnServersService,
	}

	customRouter.routerGroup.GET(
		"",
		// middlewares.AuthRequired(
		// 	customRouter.log,
		// 	customRouter.cfg,
		// ),
		// middlewares.AdminRoleRequired(
		// 	customRouter.log,
		// ),
		vpnServersHandler.GetExpandedList,
	)

	customRouter.routerGroup.GET(
		"/:id",
		// middlewares.AuthRequired(
		// 	customRouter.log,
		// 	customRouter.cfg,
		// ),
		// middlewares.AdminRoleRequired(
		// 	customRouter.log,
		// ),
		vpnServersHandler.GetExpandedItem,
	)

	customRouter.routerGroup.POST(
		"",
		// middlewares.AuthRequired(
		// 	customRouter.log,
		// 	customRouter.cfg,
		// ),
		// middlewares.AdminRoleRequired(
		// 	customRouter.log,
		// ),
		vpnServersHandler.CreateItem,
	)
}
