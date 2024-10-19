package main

import (
	"bizarre-vpn-api/internal/api/routes"
	"bizarre-vpn-api/pkg/logger"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Error.Printf("Error loading .env : %v", err)
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		logger.Error.Printf("API_PORT not found")
		return
	}

	r := routes.SetupRouter()

	if err := r.Run(apiPort); err != nil {
		log.Fatalf("Error starting API: %v", err)
	}
}
