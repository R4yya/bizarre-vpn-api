package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(swaggerPath string) *gin.Engine {
	router := gin.Default()

	RegisterPingRoute(router)

	RegisterUserRoutes(router)

	router.GET(swaggerPath+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
