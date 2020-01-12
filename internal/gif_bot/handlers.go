package gif_bot

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/impu1se/movie_to_gif_bot/internal/storage"
)

const videoUrl = "https://api.telegram.org/file/bot"

func (bot *GifBot) handlerCommands(update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		bot.handleStart(update)
	}
}

func (bot *GifBot) handlerMessages(update *tgbotapi.Update) {
	switch update.Message.Text {
	case "Новая Gif":
		bot.handleNewGif(update)
	case "Очистить время начала и конца":
		bot.handleNewGif(update)
	default:
		bot.handleTimes(update)
	}
}

func (bot *GifBot) handlerVideo(update *tgbotapi.Update) {
	chatId := update.Message.Chat.ID

	video, err := bot.api.GetFile(tgbotapi.FileConfig{update.Message.Video.FileID})
	if err != nil {
		bot.logger.Error(fmt.Sprintf("can't get file from chat id: %v, reason: %v", chatId, err))
		if err := bot.NewMessage(chatId, "Не получилось загрузить видео, попробуйте позднее", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		}
		return
	}

	err = bot.system.Download(fmt.Sprintf("%v/%v.mov", chatId, video.FileID),
		fmt.Sprintf("%v%v/%v", videoUrl, bot.Config.ApiToken, video.FilePath))
	if err != nil {
		bot.logger.Error(fmt.Sprintf("can't download video, reason %v", err))
		return
	}
	if err := bot.db.UpdateLastVideo(bot.ctx, chatId, video.FileID); err != nil {
		bot.logger.Error(fmt.Sprintf("can't update last video, reason %v", err))
		return
	}
	if err := bot.db.ClearTime(bot.ctx, chatId); err != nil {
		bot.logger.Error(fmt.Sprintf("can't clear time, reason %v", err))
		return
	}
	if err := bot.NewMessage(chatId, "Видео успешно загружено, укажите с какой секунды начать делать gif", &Clear); err != nil {
		bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		return
	}
}

func (bot *GifBot) handleStart(update *tgbotapi.Update) {
	user := &storage.User{
		ChatId:   update.Message.Chat.ID,
		UserName: update.Message.Chat.UserName,
	}

	if err := bot.db.CreateUser(bot.ctx, user); err != nil {
		bot.logger.Error("can't crete user, reason:", zap.Field{String: err.Error()})
		return
	}

	text, err := bot.db.GetText(bot.ctx, update.Message.Command())
	if err != nil {
		bot.logger.Error("can't get text, reason:", zap.Field{String: err.Error()})
		return
	}

	if err := bot.system.CreateNewDir(user.ChatId); err != nil {
		bot.logger.Error(fmt.Sprintf("can't create new dir for user with chat %v, reason %v", user.UserName, err))
	}

	if err := bot.NewMessage(user.ChatId, text, &NewGif); err != nil {
		bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		return
	}
}

func (bot *GifBot) handleNewGif(update *tgbotapi.Update) {
	chatId := update.Message.Chat.ID

	if err := bot.db.ClearTime(bot.ctx, chatId); err != nil {
		bot.logger.Error(fmt.Sprintf("can't clear time for user with id %v, reason: %v", chatId, err))
		return
	}

	text, err := bot.db.GetText(bot.ctx, update.Message.Text)
	if err != nil {
		bot.logger.Error("can't get text, reason:", zap.Field{String: err.Error()})
		return
	}

	if err := bot.NewMessage(chatId, text, &Clear); err != nil {
		bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		return
	}
}

func (bot *GifBot) handleTimes(update *tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	time, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		bot.logger.Error("can't parse time from message")
		return
	}

	user, err := bot.db.GetUser(bot.ctx, chatId)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("can't get user by chat id: %v, reason: %v", chatId, err))
		return
	}

	if user.StartTime == nil {
		if err := bot.db.UpdateStartTime(bot.ctx, chatId, time); err != nil {
			bot.logger.Error(fmt.Sprintf("can't update start time by chat id: %v, reason: %v", chatId, err))
			return
		}
		if err := bot.NewMessage(chatId, "Ок теперь введите секунду окончания видео", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}
	} else {
		if message, valid := checkValidTimes(time, *user.StartTime); !valid {
			if err := bot.NewMessage(chatId, message, nil); err != nil {
				bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			}
			return
		}

		if err := bot.db.UpdateEndTime(bot.ctx, chatId, time); err != nil {
			bot.logger.Error(fmt.Sprintf("can't update end time by chat id: %v, reason: %v", chatId, err))
			return
		}
		endTime := time - *user.StartTime
		user.EndTime = &endTime
		text, err := bot.db.GetText(bot.ctx, "create video")
		if err != nil {
			bot.logger.Error(fmt.Sprintf("can't get message from db %v", err))
			return
		}

		if err := bot.NewMessage(chatId, fmt.Sprintf(text, *user.StartTime, time), nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}

		err = bot.system.MakeImagesFromMovie(user)
		if err != nil {
			bot.logger.Error(fmt.Sprintf("can't make image from movie, reason: %v", err))
			return
		}
		if err := bot.NewMessage(chatId, "Обработка видео завершена...\nНачалось создание gif...", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}

		gifPath := fmt.Sprintf("%v/%v.gif", chatId, user.LastVideo)
		err = bot.system.MakeGif(chatId, gifPath)
		if err != nil {
			bot.logger.Error(fmt.Sprintf("can't make gif from movie, reason: %v", err))
			return
		}
		if err := bot.NewMessage(chatId, "Создание gif завершено\nЗагружаем gif в чат...", &NewGif); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}

		gif := tgbotapi.NewAnimationUpload(chatId, gifPath)
		if _, err := bot.api.Send(gif); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}
	}
}

func checkValidTimes(endTime, startTime int) (string, bool) {
	if endTime <= startTime {
		return "Конец не должен быть меньше начала", false
	} else if endTime-startTime > 10 {
		return "Продолжительность gif не должно превышать 10 сек", false
	}
	return "", true
}
