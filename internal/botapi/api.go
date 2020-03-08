package botapi

import (
	"fmt"
	"log"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/impu1se/movie_to_gif_bot/configs"
)

func NewBotApi(config *configs.Config) (*tgbotapi.BotAPI, error) {

	fmt.Println("Running bot...")
	bot, err := tgbotapi.NewBotAPI(config.ApiToken)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	bot.Debug = config.Debug
	log.Printf("Authorized on account %s", bot.Self.UserName)
	if config.Tls {
		_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert(config.Address+"/"+config.ApiToken, config.CertFile))
		if err != nil {
			return nil, err
		}
	} else {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.Address + "/" + config.ApiToken))
		if err != nil {
			return nil, err
		}
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return bot, nil
}
