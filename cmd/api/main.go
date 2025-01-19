package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	// _ "bizarre-vpn-api/docs"
	"bizarre-vpn-api/internal/api/routes"
	"bizarre-vpn-api/internal/config"
	slogWrapper "bizarre-vpn-api/internal/lib/logger"
	"bizarre-vpn-api/internal/lib/logger/sl"
	cStorage "bizarre-vpn-api/internal/storage"
)

// @title BizarreVPN API
// @version 0.0.1
// @description API for BizarreVPN project.
// @host localhost:8080
// @BasePath /
func main() {
	log := slogWrapper.SetupLogger(config.EnvLocal)

	log.Debug("debug msg")
	log.Info("info msg")
	log.Warn("warning msg")
	log.Error("warning msg")

	if err := godotenv.Load(); err != nil {
		log.Error("config initialization error", sl.Err(err), slog.String("test", "test"))
		return
	}

	portString := os.Getenv("API_PORT")

	log.Info("starting application",
		slog.String("env", "env example"),
		slog.String("port", portString),
	)

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		log.Error("DATABASE_PATH not found")
		return
	}

	log.Info("database initialization")

	storage, err := cStorage.New(dbPath)

	if err != nil {
		log.Error("database initialization error", sl.Err(err))
		return
	}

	log.Info(fmt.Sprintf("Connected to SQLite database at %s", dbPath))

	defer storage.CloseDB()

	log.Info("database initialized successful")

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		log.Error("API_PORT not found")
		return
	}

	swaggerPath := os.Getenv("SWAGGER_PATH")
	if swaggerPath == "" {
		log.Error("SWAGGER_PATH not found")
		return
	}

	r := routes.SetupRouter(swaggerPath)

	log.Info("API successfully started")

	if err := r.Run(apiPort); err != nil {
		log.Error("server listening error", sl.Err(err))
		return
	}
}
