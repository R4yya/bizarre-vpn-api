package main

import (
	"os"
	"time"

	"bizarre-vpn-api/internal/bot"
	"bizarre-vpn-api/pkg/logger"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func main() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		logger.Error.Printf("Error loading .env : %v", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		logger.Error.Printf("TELEGRAM_BOT_TOKEN not found")
	}

	webAppUrl := os.Getenv("WEB_APP_URL")
	if webAppUrl == "" {
		logger.Error.Printf("WEB_APP_URL not found")
	}

	botSettings := tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(botSettings)
	if err != nil {
		logger.Error.Printf("Error creating bot: %v", err)
		return
	}

	bot.RegisterHandlers(b, webAppUrl)

	logger.Info.Println("Bot successfully started")

	b.Start()
}
