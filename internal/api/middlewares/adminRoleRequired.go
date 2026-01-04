package middlewares

import (
	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/api/helpers"
	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/shared/logger/sl"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRoleRequired(log *slog.Logger) gin.HandlerFunc {
	const op = "internal.api.middlewares.AdminRoleRequired"

	return func(c *gin.Context) {
		log = log.With(slog.String("op", op))

		tokenInfo, err := helpers.GetTokenInfo(c)

		if err != nil {
			log.Info("getting token info error", sl.Err(err))
			c.JSON(http.StatusForbidden, handlers.ErrorResponse{Error: err.Error()})
			c.Abort()
			return
		}

		if tokenInfo.Role != models.UserRoleAdmin {
			msg := fmt.Sprintf("role matching error, needed role is %v", models.UserRoleAdmin)
			log.Info(msg, slog.String("role", tokenInfo.Role))
			c.JSON(http.StatusForbidden, handlers.ErrorResponse{Error: msg})
			c.Abort()
			return
		}

		c.Next()
	}
}
