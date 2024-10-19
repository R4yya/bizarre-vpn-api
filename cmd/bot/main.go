package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"bizarre-vpn-api/internal/bot"
	"bizarre-vpn-api/pkg/logger"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func main() {
	err := logger.Init("bot")
	if err != nil {
		log.Fatalf("Error initiating logger: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		logger.Error(err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		err := fmt.Errorf("TELEGRAM_BOT_TOKEN not found")
		logger.Error(err)
		return
	}

	webAppUrl := os.Getenv("WEB_APP_URL")
	if webAppUrl == "" {
		err := fmt.Errorf("WEB_APP_URL not found")
		logger.Error(err)
		return
	}

	botSettings := tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(botSettings)
	if err != nil {
		logger.Error(err)
		return
	}

	bot.RegisterHandlers(b, webAppUrl)

	logger.Info("Bot successfully started")

	b.Start()
}
