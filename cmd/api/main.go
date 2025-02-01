package main

import (
	"log/slog"
	"strconv"

	"bizarre-vpn-api/internal/api/routes"
	botInternal "bizarre-vpn-api/internal/bot"
	"bizarre-vpn-api/internal/config"
	slogWrapper "bizarre-vpn-api/internal/lib/logger"
	"bizarre-vpn-api/internal/lib/logger/sl"
	cStorage "bizarre-vpn-api/internal/storage/sqlite"
)

// @title BizarreVPN API
// @version 0.0.1
// @description API for BizarreVPN project.
// @host 127.0.0.1:5050
// @BasePath /
func main() {
	cfg := config.MustLoadConfig()

	log := slogWrapper.SetupLogger(cfg.Env)

	log.Info("starting application",
		slog.String("env", cfg.Env),
		slog.String("host", cfg.Env),
		slog.Int("port", cfg.HttpServer.Port),
	)

	log.Info("database initialization")

	storage := cStorage.MustInit(cfg.StoragePath, log)

	//TODO: need to relocate func call to graceful shutdown
	defer storage.CloseDB()

	log.Info("database initialized successful")

	bot := botInternal.MustInitBot(log, cfg.TelegramBotToken, cfg.WebAppUrl, storage)

	go bot.Start()
	log.Info("TG bot successfully started")

	r := routes.SetupRouter(log)

	log.Info("API successfully started")

	serverAddress := cfg.HttpServer.Host + ":" + strconv.Itoa(cfg.HttpServer.Port)

	log.Debug("server address", slog.String("host", cfg.HttpServer.Host), slog.Int("port", cfg.HttpServer.Port))

	if err := r.Run(serverAddress); err != nil {
		log.Error("server listening error", sl.Err(err))
		return
	}
}
