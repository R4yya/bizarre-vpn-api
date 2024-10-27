package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"time"
)

// CORS loads the CORS configuration from environment variables and returns middleware
func CORS() gin.HandlerFunc {
	_ = godotenv.Load()

	allowOrigins := strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ",")
	allowMethods := strings.Split(os.Getenv("CORS_ALLOW_METHODS"), ",")
	allowHeaders := strings.Split(os.Getenv("CORS_ALLOW_HEADERS"), ",")
	exposeHeaders := strings.Split(os.Getenv("CORS_EXPOSE_HEADERS"), ",")
	allowCredentials := os.Getenv("CORS_ALLOW_CREDENTIALS") == "true"

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		ExposeHeaders:    exposeHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	})
}
