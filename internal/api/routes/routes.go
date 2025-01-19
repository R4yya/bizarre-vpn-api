package routes

import (
	"bizarre-vpn-api/internal/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const swaggerPath = "/docs"

func SetupRouter(swaggerPath string) *gin.Engine {
	router := gin.Default()

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "method not allowed on this endpoint",
		})
	})

	corsMiddleware := middleware.CORS()
	router.Use(corsMiddleware)

	RegisterPingRoute(router)
	RegisterUserRoutes(router)
	RegisterSubscriptionPlanRoutes(router)

	router.GET(swaggerPath+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
