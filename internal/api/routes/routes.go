package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"bizarre-vpn-api/internal/api/middleware"
)

func SetupRouter(log *slog.Logger) *gin.Engine {
	customRouter := SetupCustomRouter(log)

	customRouter._router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "method not allowed on this endpoint",
		})
	})

	corsMiddleware := middleware.CORS()
	customRouter._router.Use(corsMiddleware)

	customRouter.AddGroup("/ping", PingRoute)
	customRouter.AddGroup("/users", UserRoutes)
	customRouter.AddGroup("/plans", SubscriptionPlanRoutes)
	customRouter.AddGroup("/docs", DocsRoutes)

	return customRouter._router
}
