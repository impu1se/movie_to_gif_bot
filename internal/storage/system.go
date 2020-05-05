package storage

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"go.uber.org/zap"

	"github.com/disintegration/imaging"
)

const dataDir = "user_data/"

type System struct {
	logger *zap.Logger
}

func NewLoader(logger *zap.Logger) *System {
	err := os.Mkdir(dataDir, os.FileMode(0777))
	if err != nil {
		logger.Error("can't create data dir")
	}
	return &System{logger: logger}
}

func (l *System) Download(filepath, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dataDir + filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (l *System) MakeGif(chatId int64, dest string) error {

	path := fmt.Sprintf(dataDir+"%v/*.jpg", chatId)

	srcfilenames, err := filepath.Glob(path)
	if err != nil {
		l.logger.Error(fmt.Sprintf("error in globbing source file pattern %v : %v", path, err))
		return err
	}
	if len(srcfilenames) == 0 {
		log.Fatalf("No source images found via pattern %s", path)
	}
	sort.Strings(srcfilenames)

	var frames []*image.Paletted

	for _, filename := range srcfilenames {
		img, err := imaging.Open(filename)
		if err != nil {
			log.Printf("Skipping file %s due to error reading it :%s", filename, err)
			continue
		}

		img = ScaleImage(0.4, img)

		buf := bytes.Buffer{}
		if err := gif.Encode(&buf, img, nil); err != nil {
			log.Printf("Skipping file %s due to error in gif encoding:%s", filename, err)
			continue
		}

		tmpimg, err := gif.Decode(&buf)
		if err != nil {
			log.Printf("Skipping file %s due to weird error reading the temporary gif :%s", filename, err)
			continue
		}
		frames = append(frames, tmpimg.(*image.Paletted))

	}

	delays := make([]int, len(frames))
	for j, _ := range delays {
		delays[j] = 3
	}

	opfile, err := os.Create(dataDir + dest)
	if err != nil {
		log.Fatalf("Error creating the destination file %s : %s", dest, err)
	}

	if err := gif.EncodeAll(opfile, &gif.GIF{Image: frames, Delay: delays}); err != nil {
		log.Printf("Error encoding output into animated gif :%s", err)
	}
	if err = opfile.Close(); err != nil {
		panic(err)
	}
	return nil
}

func ScaleImage(scale float64, img image.Image) image.Image {
	newwidth := int(float64(img.Bounds().Dx()) * scale)
	newheight := int(float64(img.Bounds().Dy()) * scale)

	img = imaging.Resize(img, newwidth, newheight, imaging.Lanczos)
	return img

}

func (l *System) ClearDir(pattern string) error {
	files, err := filepath.Glob(dataDir + pattern)
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func (l *System) CreateNewDir(chatId int64) error {
	return os.Mkdir(dataDir+fmt.Sprint(chatId), os.FileMode(0777))
}

func (l *System) MakeImagesFromMovie(user *User) error {

	path := fmt.Sprintf("%v/*.jpg", user.ChatId)

	if err := l.ClearDir(path); err != nil {
		l.logger.Error(fmt.Sprintf("can't clear dir %v , err: %v", path, err))
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	mplayer := exec.Command("/usr/bin/mplayer", "-vo",
		fmt.Sprintf("jpeg:outdir=%v/%v%v:quality=100", pwd, dataDir, user.ChatId),
		"-nosound", "-ss", fmt.Sprint(*user.StartTime), "-endpos", fmt.Sprint(*user.EndTime),
		fmt.Sprintf(dataDir+"%v/%v.mov", user.ChatId, user.LastVideo))
	mplayer.Stderr = os.Stderr
	mplayer.Stdout = os.Stdout
	return mplayer.Run()
}
