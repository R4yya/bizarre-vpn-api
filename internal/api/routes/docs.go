package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func DocsRoutes(customRouter *CustomRouter) {
	baseUrl := customRouter.cfg.HttpServer.BaseURL

	redirectUrl := baseUrl + "/swagger/index.html"
	customRouter.log.Debug("redirect swagger url", slog.String("redirectUrl", redirectUrl))

	customRouter._router.GET("", func(c *gin.Context) {

		c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
	})

	customRouter._router.Static("/docs", "./docs")
	customRouter._router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL(baseUrl+"/docs/swagger.json"),
	))
}
