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

// FOR TESTING !!!!

//const (
//	gifName = "BAADAgADGAYAAl_esEj3k4FBsq-GYhYE"
//	dir     = "138177057"
//	scale   = 0.33
//	delay   = 3
//	quality = "100"
//)
//
//func ClearDir(pattern string) error {
//	files, err := filepath.Glob(pattern)
//	if err != nil {
//		return err
//	}
//
//	for _, f := range files {
//		if err := os.Remove(f); err != nil {
//			return err
//		}
//	}
//	return nil
//}

//func main() {
//
//	if err := ClearDir(fmt.Sprintf("%v/*.jpg", dir)); err != nil {
//		panic(err)
//	}
//
//	err := MakeImagesFromMovie()
//	if err != nil {
//		panic(err)
//	}
//
//	path := "138177057/*.jpg"
//	srcfilenames, err := filepath.Glob(path)
//	if err != nil {
//		panic(err)
//	}
//	if len(srcfilenames) == 0 {
//		log.Fatalf("No source images found via pattern %s", path)
//	}
//	sort.Strings(srcfilenames)
//
//	var frames []*image.Paletted
//
//	for _, filename := range srcfilenames {
//		img, err := imaging.Open(filename)
//		if err != nil {
//			log.Printf("Skipping file %s due to error reading it :%s", filename, err)
//			continue
//		}
//
//		img = ScaleImage(scale, img)
//
//		buf := bytes.Buffer{}
//		if err := gif.Encode(&buf, img, nil); err != nil {
//			log.Printf("Skipping file %s due to error in gif encoding:%s", filename, err)
//			continue
//		}
//
//		tmpimg, err := gif.Decode(&buf)
//		if err != nil {
//			log.Printf("Skipping file %s due to weird error reading the temporary gif :%s", filename, err)
//			continue
//		}
//		frames = append(frames, tmpimg.(*image.Paletted))
//
//	}
//
//	delays := make([]int, len(frames))
//	for j, _ := range delays {
//		delays[j] = delay
//	}
//
//	dest := fmt.Sprintf("%v/%v_%v_%v_%v.gif", dir, gifName, scale, quality, delay)
//	opfile, err := os.Create(dest)
//	if err != nil {
//		log.Fatalf("Error creating the destination file %s : %s", dest, err)
//	}
//
//	if err := gif.EncodeAll(opfile, &gif.GIF{Image: frames, Delay: delays}); err != nil {
//		log.Printf("Error encoding output into animated gif :%s", err)
//	}
//	if err = opfile.Close(); err != nil {
//		panic(err)
//	}
//
//}
//
//func ScaleImage(scale float64, img image.Image) image.Image {
//	newwidth := int(float64(img.Bounds().Dx()) * scale)
//	newheight := int(float64(img.Bounds().Dy()) * scale)
//
//	img = imaging.Resize(img, newwidth, newheight, imaging.Lanczos)
//	return img
//
//}
//
//func MakeImagesFromMovie() error {
//	pwd, err := os.Getwd()
//	if err != nil {
//		return err
//	}
//
//	StartTime := "30"
//	EndTime := "5"
//	mplayer := exec.Command("mplayer", "-vo",
//		fmt.Sprintf("jpeg:outdir=%v/%v:quality="+quality, pwd, dir),
//		"-nosound", "-ss", StartTime, "-endpos", EndTime,
//		fmt.Sprintf("%v/%v.mov", dir, gifName))
//	mplayer.Stderr = os.Stderr
//	mplayer.Stdout = os.Stdout
//	return mplayer.Run()
//}
