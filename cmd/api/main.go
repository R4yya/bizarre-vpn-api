package main

import (
	"log/slog"
	"strconv"

	"bizarre-vpn-api/internal/api/routes"
	botInternal "bizarre-vpn-api/internal/bot"
	"bizarre-vpn-api/internal/config"
	slogWrapper "bizarre-vpn-api/internal/lib/logger"
	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/services"
	cStorage "bizarre-vpn-api/internal/storage/sqlite"
)

// @title BizarreVPN API
// @version 0.0.1
// @description API for BizarreVPN project.
// @host localhost:5050
// @BasePath /
// @securityDefinitions.apikey token
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345".
// @authorizationurl /users/auth/telegram-init-data

func main() {
	cfg := config.MustLoadConfig()

	log := slogWrapper.SetupLogger(cfg.Env)

	log.Info("starting application",
		slog.String("env", cfg.Env),
		slog.String("host", cfg.Env),
		slog.Int("port", cfg.HttpServer.Port),
		slog.String("baseURL", cfg.HttpServer.BaseURL),
	)

	log.Info("database initialization")

	storage := cStorage.MustInit(cfg.StoragePath, log)

	userService := services.NewUserService(
		log,
		storage.UserStorage,
	)

	userService.CreateDefaultUser()

	//TODO: need to relocate func call to graceful shutdown
	defer func() {
		err := storage.CloseDB()

		if err != nil {
			log.Error("storage.CloseDB error", sl.Err(err))
		}

		log.Info("app is shutdown")
	}()

	log.Info("database initialized successful")

	bot := botInternal.MustInitBot(log, cfg.TelegramBotToken, cfg, storage)

	go bot.Start()
	log.Info("TG bot successfully started")

	r := routes.SetupRouter(log, cfg, storage)

	log.Info("API successfully started")

	serverAddress := ":" + strconv.Itoa(cfg.HttpServer.Port)

	log.Debug("server info", slog.Int("port", cfg.HttpServer.Port))

	if err := r.Run(serverAddress); err != nil {
		log.Error("server listening error", sl.Err(err))
		return
	}
}
