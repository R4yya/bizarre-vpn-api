package main

import (
	"fmt"
	"os"
	"time"

	"bizarre-vpn-api/internal/bot"
	"bizarre-vpn-api/pkg/logger"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func main() {
	logger.Init("bot")

	if err := godotenv.Load(); err != nil {
		logger.Error(err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		err := fmt.Errorf("TELEGRAM_BOT_TOKEN not found")
		logger.Error(err)
	}

	webAppUrl := os.Getenv("WEB_APP_URL")
	if webAppUrl == "" {
		err := fmt.Errorf("WEB_APP_URL not found")
		logger.Error(err)
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
