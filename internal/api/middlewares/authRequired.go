package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/api/helpers"
	"bizarre-vpn-api/internal/config"
)

func AuthRequired(log *slog.Logger, cfg *config.Config) gin.HandlerFunc {
	op := "internal.api.middleware"

	return func(c *gin.Context) {
		log = log.With(slog.String("op", op))

		tokenInfo, _, err := helpers.ParseTokenFromHeader(c, log, "Authorization", cfg.JWT.AccessSecretKey)

		if err != nil {
			c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Error: err.Error()})
			c.Abort()
			return
		}

		// Set example variable
		c.Set("tokenInfo", tokenInfo)

		// before request

		c.Next()
	}
}
