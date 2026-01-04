package bot

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"unicode/utf8"

	tele "gopkg.in/telebot.v4"

	"bizarre-vpn-api/internal/services"
	"bizarre-vpn-api/internal/shared/config"
	"bizarre-vpn-api/internal/shared/coreErrors"
	"bizarre-vpn-api/internal/shared/logger/sl"
	cStorage "bizarre-vpn-api/internal/storage/sqlite"
)

func registerHandlers(b *tele.Bot, cfg *config.Config, log *slog.Logger, storage *cStorage.Storage) {
	b.Handle("/start", func(c tele.Context) error {
		return handleStart(c, cfg.WebAppUrl, log, storage)
	})

	// b.Handle("/auth", func(c tele.Context) error {
	// 	return handleAuth(c, cfg, log, storage)
	// })

	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Извините, я понимаю только команду /start.")
	})
}

func handleStart(c tele.Context, webAppUrl string, log *slog.Logger, storage *cStorage.Storage) error {
	op := "internal.bot.handlers.handleStart"

	log = log.With(slog.String("op", op))

	teleUser := c.Sender()

	if teleUser.IsBot {
		log.Info("bot user request denied")
		return c.Send("Извините, мы не работаем с ботами")
	}

	payload := c.Message().Payload

	if payload == "" || utf8.RuneCountInString(payload) != services.LinkUserCodeLength {
		return c.Send("Запросите инвайт ссылку у представителя bizarre")
	}

	authLinkServiceInstance := services.NewAuthLinksService(log, storage.AuthLinksStorage, storage.LnkUserProviderStorage)

	preparedExternalId := strconv.Itoa(int(teleUser.ID))

	userId, err := authLinkServiceInstance.LinkUserWithTgProviderByCode(payload, preparedExternalId)

	if err != nil {
		if errors.Is(err, coreErrors.ErrorUserAlreadyLinked) {
			return c.Send("Ваш аккаунт уже привязан")
		}

		if errors.Is(err, coreErrors.ErrorNotFound) {
			return c.Send("Ссылка не действительна")
		}

		log.Error("linking tg acc to user error", sl.Err(err))

		return c.Send("Что то пошло не так.")
	}

	// webApp := tele.WebApp{URL: webAppUrl}
	// btn := tele.InlineButton{Text: "Открыть BizarreVPN", WebApp: &webApp}

	// inlineKeyboard := [][]tele.InlineButton{
	// 	{btn},
	// }

	message := fmt.Sprintf("Вы успешно привязали аккаунт. userId: %d", userId)

	return c.Send(message)

	// return c.Send("Нажми на кнопку, чтобы открыть Mini App.", &tele.ReplyMarkup{
	// 	InlineKeyboard: inlineKeyboard,
	// })
}

// func handleAuth(c tele.Context, cfg *config.Config, log *slog.Logger, storage *cStorage.Storage) error {
// 	op := "internal.bot.handlers.handleAuth"

// 	log = log.With(slog.String("op", op))

// 	teleUser := c.Sender()

// 	if teleUser.IsBot {
// 		log.Info("bot user request denied")
// 		_ = c.Send("Извините, мы не работает с ботами")
// 	}

// 	log.Debug("telegram user data",
// 		slog.Int64("teleUser.ID", teleUser.ID),
// 		slog.String("teleUser.Username", teleUser.Username),
// 	)

// 	authService := services.NewAuthService(
// 		log,
// 		storage.UserStorage,
// 		storage.LnkUserProviderStorage,
// 	)
// 	accessToken, refreshToken, err := authService.AuthorizeByTelegram(
// 		teleUser.ID,
// 		teleUser.Username,
// 		[]byte(cfg.JWT.AccessSecretKey),
// 		[]byte(cfg.JWT.RefreshSecretKey),
// 	)

// 	if err != nil {
// 		log.Error(op, sl.Err(err))
// 		_ = c.Send("Произошла непредвиденная ошибка, попробуйте ещё раз")

// 		return nil
// 	}

// 	mes := fmt.Sprintf("Добро пожаловать, %v! \n\n%v \n\n%v", teleUser.FirstName, accessToken, refreshToken)

// 	return c.Send(mes)
// }
