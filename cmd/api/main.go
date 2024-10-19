package main

import (
	_ "bizarre-vpn-api/docs"
	"bizarre-vpn-api/internal/api/routes"
	"bizarre-vpn-api/pkg/logger"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// @title BizarreVPN API
// @version 0.0.1
// @description API for BizarreVPN project.
// @host localhost:8080
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		logger.Error.Printf("Error loading .env : %v", err)
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		logger.Error.Printf("API_PORT not found")
		return
	}

	swaggerPath := os.Getenv("SWAGGER_PATH")
	if swaggerPath == "" {
		logger.Error.Printf("SWAGGER_PATH not found")
		return
	}

	r := routes.SetupRouter(swaggerPath)

	if err := r.Run(apiPort); err != nil {
		log.Fatalf("Error starting API: %v", err)
	}
}
