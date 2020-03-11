package gif_bot

import tgbotapi "github.com/Syfaro/telegram-bot-api"

const (
	clearTimes = "Очистить время начала и конца"
	newGif     = "Gif из нового видео"
	oldGif     = "Gif из того же видео"
)

var Clear = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(clearTimes),
	),
)

var NewGif = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(newGif),
	),
)

var Else = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(oldGif),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(newGif),
	),
)
