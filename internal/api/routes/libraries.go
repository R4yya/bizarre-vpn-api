package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/api/middlewares"
	"bizarre-vpn-api/internal/services"
)

func LibrariesRoutes(customRouter *CustomRouter) {
	librariesService := services.NewLibrariesService(
		customRouter.log,
		customRouter.storage.BackendTypesStorage,
	)

	librariesHandler := handlers.LibrariesHandler{
		Log:              customRouter.log,
		LibrariesService: librariesService,
	}

	customRouter.routerGroup.GET(
		"/backend-types",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
		librariesHandler.GetBackendTypesList,
	)
}
