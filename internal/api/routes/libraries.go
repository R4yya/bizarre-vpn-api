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
		customRouter.storage.ProtocolsStorage,
		customRouter.storage.LnkProtocolsBackendTypesStorage,
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

	customRouter.routerGroup.GET(
		"/protocols",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
		librariesHandler.GetProtocolsList,
	)

	customRouter.routerGroup.GET(
		"/protocols-by-backend-type-id",
		//middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		//middlewares.AdminRoleRequired(customRouter.log),
		librariesHandler.GetProtocolsListByBackendTypeId,
	)
}
