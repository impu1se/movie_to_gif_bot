package gif_bot

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/impu1se/movie_to_gif_bot/internal/storage"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/impu1se/movie_to_gif_bot/configs"
)

type System interface {
	Download(filepath, url string) error
	CreateNewDir(chatId int64) error
	MakeGif(chatId int64, dest string) error
	MakeImagesFromMovie(user *storage.User) error
	ClearDir(pattern string) error
}

type GifBot struct {
	Config  *configs.Config
	Updates tgbotapi.UpdatesChannel
	db      *storage.Database
	system  System
	logger  *zap.Logger
	ctx     context.Context
	api     tgbotapi.BotAPI
}

func NewGifBot(
	config *configs.Config,
	updates tgbotapi.UpdatesChannel,
	system System,
	db *storage.Database,
	logger *zap.Logger,
	api tgbotapi.BotAPI,
	ctx context.Context,
) *GifBot {
	return &GifBot{
		Config:  config,
		Updates: updates,
		system:  system,
		db:      db,
		logger:  logger,
		ctx:     ctx,
		api:     api,
	}
}

func (bot *GifBot) Run() {
	bot.logger.Info("in run...")
	for update := range bot.Updates {
		bot.logger.Info("got update...")
		if update.Message == nil {
			continue
		}

		if update.Message.Video != nil {
			bot.handlerVideo(&update)
			continue
		}

		if update.Message.IsCommand() {
			bot.handlerCommands(&update)
			continue
		}

		if update.Message != nil {
			bot.handlerMessages(&update)
			continue
		}
	}
}

func (bot *GifBot) NewMessage(chatId int64, message string, button *tgbotapi.ReplyKeyboardMarkup) error {

	if message == "" {
		return nil
	}
	text, err := bot.db.GetText(bot.ctx, message)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("can't get text from db for message: %v", message))
		return err
	}
	msg := tgbotapi.NewMessage(chatId, text)
	if button != nil {
		msg.ReplyMarkup = button
	}
	if _, err := bot.api.Send(msg); err != nil {
		return err
	}
	return nil
}
