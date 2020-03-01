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
	case "Gif из нового видео":
		bot.handleNewGif(update)
	case "Очистить время начала и конца", "Gif из того же видео":
		bot.handleNewGif(update)
	default:
		bot.handleTimes(update)
	}
}

func (bot *GifBot) handlerVideo(update *tgbotapi.Update) {
	chatId := update.Message.Chat.ID

	if err := bot.system.ClearDir(fmt.Sprintf("%v/*.mov", chatId)); err != nil {
		bot.logger.Error("can't clear dir for new video")
		return
	}

	video, err := bot.api.GetFile(tgbotapi.FileConfig{update.Message.Video.FileID}) // TODO: make check file size
	if err != nil {
		bot.logger.Error(fmt.Sprintf("can't get file from chat id: %v, reason: %v", chatId, err))
		if err := bot.NewMessage(chatId, "download error", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		}
		return
	} else {
		if err := bot.NewMessage(chatId, "save video", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		}
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
	if err := bot.NewMessage(chatId, "successful download", &Clear); err != nil {
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

	if err := bot.system.CreateNewDir(user.ChatId); err != nil {
		bot.logger.Error(fmt.Sprintf("can't create new dir for user with chat %v, reason %v", user.UserName, err))
	}

	if err := bot.NewMessage(user.ChatId, update.Message.Command(), &NewGif); err != nil {
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

	if err := bot.NewMessage(chatId, update.Message.Text, &Clear); err != nil {
		bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		return
	}
}

func (bot *GifBot) handleTimes(update *tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	time, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		bot.logger.Error("can't parse time from message")
		if err := bot.NewMessage(chatId, "invalid message", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
		}
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
		if err := bot.NewMessage(chatId, "end second", nil); err != nil {
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
		if err := bot.NewMessage(chatId, "create video", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}

		err = bot.system.MakeImagesFromMovie(user)
		if err != nil {
			bot.logger.Error(fmt.Sprintf("can't make image from movie, reason: %v", err))
			return
		}
		if err := bot.NewMessage(chatId, "start create video", nil); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}

		gifPath := fmt.Sprintf("%v/%v.gif", chatId, user.LastVideo)
		err = bot.system.MakeGif(chatId, gifPath)
		if err != nil {
			bot.logger.Error(fmt.Sprintf("can't make gif from movie, reason: %v", err))
			return
		}
		if err := bot.NewMessage(chatId, "loading gif", &NewGif); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}

		gif := tgbotapi.NewAnimationUpload(chatId, "user_data/"+gifPath) // TODO: make user_data from system
		gif.ReplyMarkup = Else
		if _, err := bot.api.Send(gif); err != nil {
			bot.logger.Error(fmt.Sprintf("can't send message, reason: %v", err))
			return
		}
	}
}

func checkValidTimes(endTime, startTime int) (string, bool) {
	if endTime <= startTime {
		return "end more start", false
	} else if endTime-startTime > 10 {
		return "video more 10s", false
	}
	return "", true
}
