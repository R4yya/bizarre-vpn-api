package main

import (
	"bizarre-vpn-api/pkg/logger"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Error.Printf("Error loading .env : %v", err)
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		logger.Error.Printf("API_PORT not found")
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(apiPort); err != nil {
		log.Fatalf("Error starting API: %v", err)
	}
}
