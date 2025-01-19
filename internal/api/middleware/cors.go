package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS loads the CORS configuration from environment variables and returns middleware
func CORS() gin.HandlerFunc {
	allowOrigins := []string{"*"}
	allowMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	allowHeaders := []string{"Origin", "Content-Type", "Authorization"}
	exposeHeaders := []string{"Content-Length"}
	allowCredentials := true

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		ExposeHeaders:    exposeHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	})
}
