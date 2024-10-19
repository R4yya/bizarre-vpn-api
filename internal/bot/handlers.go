package bot

import (
	tele "gopkg.in/telebot.v4"
)

func handleStart(c tele.Context, webAppUrl string) error {
	webApp := tele.WebApp{URL: webAppUrl}
	btn := tele.InlineButton{Text: "Открыть BizarreVPN", WebApp: &webApp}

	inlineKeyboard := [][]tele.InlineButton{
		{btn},
	}

	return c.Send("Нажми на кнопку, чтобы открыть Mini App.", &tele.ReplyMarkup{
		InlineKeyboard: inlineKeyboard,
	})
}

func RegisterHandlers(b *tele.Bot, webAppUrl string) {
	b.Handle("/start", func(c tele.Context) error {
		return handleStart(c, webAppUrl)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Извините, я понимаю только команду /start.")
	})
}
