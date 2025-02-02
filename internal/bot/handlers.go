package bot

import (
	"fmt"
	"log/slog"

	tele "gopkg.in/telebot.v4"

	"bizarre-vpn-api/internal/config"
	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/services"
	cStorage "bizarre-vpn-api/internal/storage/sqlite"
)

func registerHandlers(b *tele.Bot, cfg *config.Config, log *slog.Logger, storage *cStorage.Storage) {
	b.Handle("/start", func(c tele.Context) error {
		return handleStart(c, cfg, log, storage)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Извините, я понимаю только команду /start.")
	})
}

func handleStart(c tele.Context, cfg *config.Config, log *slog.Logger, storage *cStorage.Storage) error {
	op := "internal.bot.handlers.handleStart"
	teleUser := c.Sender()

	if teleUser.IsBot {
		c.Send("Извините, мы не работает с ботами")
	}

	log.Debug("telegram user data",
		slog.Int64("teleUser.ID", teleUser.ID),
		slog.String("teleUser.Username", teleUser.Username),
	)

	authService := services.NewAuthService(
		log,
		storage.UserStorage,
		storage.LnkUserProviderStorage,
	)
	accessToken, refreshToken, err := authService.AuthorizeByTelegram(
		teleUser.ID,
		teleUser.Username,
		[]byte(cfg.JWT.AccessSecretKey),
		[]byte(cfg.JWT.RefreshSecretKey),
	)

	if err != nil {
		log.Error(op, sl.Err(err))
		c.Send("Произошла непредвиденная ошибка, попробуйте ещё раз")

		return nil
	}

	mes := fmt.Sprintf("Добро пожаловать, %v! Ваш аккаунт был создан. %v, %v", teleUser.FirstName, accessToken, refreshToken)

	c.Send(mes)

	return nil
}

// func handleStart(c tele.Context, webAppUrl string) error {
// 	webApp := tele.WebApp{URL: webAppUrl}
// 	btn := tele.InlineButton{Text: "Открыть BizarreVPN", WebApp: &webApp}

// 	inlineKeyboard := [][]tele.InlineButton{
// 		{btn},
// 	}

// 	return c.Send("Нажми на кнопку, чтобы открыть Mini App.", &tele.ReplyMarkup{
// 		InlineKeyboard: inlineKeyboard,
// 	})
// }
