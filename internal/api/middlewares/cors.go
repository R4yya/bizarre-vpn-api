package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS loads the CORS configuration from environment variables and returns middleware
func CORS(originsDomains []string) gin.HandlerFunc {
	allowOrigins := originsDomains
	allowMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	allowHeaders := []string{"Origin", "Content-Type", "Authorization"}
	allowCredentials := true

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: allowCredentials,
	})
}
