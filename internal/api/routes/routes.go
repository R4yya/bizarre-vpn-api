package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"bizarre-vpn-api/internal/api/middlewares"
	"bizarre-vpn-api/internal/config"
	"bizarre-vpn-api/internal/storage/sqlite"
)

func SetupRouter(log *slog.Logger, cfg *config.Config, storage *sqlite.Storage) *gin.Engine {
	customRouter := SetupCustomRouter(log, cfg, storage)
	gin.SetMode(gin.ReleaseMode)

	customRouter._router.Use(middlewares.RequestsLogger(log))
	customRouter._router.Use(middlewares.CORS(cfg.HttpServer.AllowOrigins))

	customRouter._router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "method not allowed on this endpoint",
		})
	})

	customRouter.AddGroup("/ping", PingRoute)
	customRouter.AddGroup("/users", UserRoutes)
	//customRouter.AddGroup("/plans", SubscriptionPlanRoutes)
	DocsRoutes(customRouter)

	return customRouter._router
}
