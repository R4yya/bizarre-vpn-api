package main

import (
	_ "bizarre-vpn-api/docs"
	"bizarre-vpn-api/internal/api/routes"
	"bizarre-vpn-api/pkg/logger"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// @title BizarreVPN API
// @version 0.0.1
// @description API for BizarreVPN project.
// @host localhost:8080
// @BasePath /
func main() {
	logger.Init("api")
	defer logger.Close()

	if err := godotenv.Load(); err != nil {
		logger.Error(err)
		return
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		err := fmt.Errorf("API_PORT not found")
		logger.Error(err)
		return
	}

	swaggerPath := os.Getenv("SWAGGER_PATH")
	if swaggerPath == "" {
		err := fmt.Errorf("SWAGGER_PATH not found")
		logger.Error(err)
		return
	}

	r := routes.SetupRouter(swaggerPath)

	logger.Info("API successfully started")

	if err := r.Run(apiPort); err != nil {
		logger.Error(err)
		return
	}
}
