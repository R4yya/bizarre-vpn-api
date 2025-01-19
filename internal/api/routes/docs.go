package routes

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func DocsRoutes(customRouter *CustomRouter) {
	customRouter.routerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
