package routes

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func DocsRoutes(customRouter *CustomRouter) {
	customRouter._router.Static("/docs", "./docs")
	customRouter._router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/docs/swagger.json"),
	))
}
