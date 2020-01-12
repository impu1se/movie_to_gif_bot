package gif_bot

import tgbotapi "github.com/Syfaro/telegram-bot-api"

var Clear = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Очистить время начала и конца"),
	),
)

var NewGif = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Новая Gif"),
	),
)
