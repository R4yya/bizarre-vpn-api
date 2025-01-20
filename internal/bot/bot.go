package bot

import (
	"fmt"
	"log/slog"
	"time"

	tele "gopkg.in/telebot.v4"

	"bizarre-vpn-api/internal/lib/logger/sl"
)

func MustInitBot(log *slog.Logger, botToken string, webAppUrl string) *tele.Bot {
	log.Info("tg bot starting")

	botSettings := tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(botSettings)
	if err != nil {
		log.Error("tg bot error", sl.Err(err))
		panic(fmt.Errorf("tg bot error: %v", err.Error()))
	}

	registerHandlers(b, webAppUrl)

	return b
}
