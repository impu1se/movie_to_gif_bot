package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/impu1se/movie_to_gif_bot/configs"
	"github.com/impu1se/movie_to_gif_bot/internal/botapi"
	"github.com/impu1se/movie_to_gif_bot/internal/gif_bot"
	"github.com/impu1se/movie_to_gif_bot/internal/storage"
	"go.uber.org/zap"
)

func main() {

	config := configs.NewConfig()
	if config.Tls {
		go http.ListenAndServeTLS(":"+config.Port, config.CertFile, config.KeyFile, nil)
	} else {
		go http.ListenAndServe(":"+config.Port, nil)
	}

	botApi, err := botapi.NewBotApi(config)
	if err != nil {
		log.Fatalf("can't get new bot api, reason: %v", err)
	}

	db, err := storage.NewDb(config)
	if err != nil {
		log.Fatalf("can't create db, reason: %v", err)
	}

	logger := zap.NewExample()

	system := storage.NewLoader(logger)
	gifBot := gif_bot.NewGifBot(config, botApi.ListenForWebhook("/"+botApi.Token), system, db, logger, *botApi, context.Background())

	fmt.Printf("Start server on %v:%v ", config.Address, config.Port)
	gifBot.Run()
}
