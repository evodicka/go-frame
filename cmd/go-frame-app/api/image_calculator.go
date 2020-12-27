package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.com/go-displays/go-frame/cmd/go-frame-app/model"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type ImageRef struct {
	Path     string     `json:"path" binding:"required"`
	Type     model.Type `json:"type" binding:"required"`
	Metadata string     `json:"metadata"`
}

func getCurrentImageData(context *gin.Context) {
	fileName, err := calculateCurrentImage()
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ref := ImageRef{Path: fileName, Type: model.ImageType}
	context.JSON(http.StatusOK, ref)
}

func calculateCurrentImage() (path string, err error) {
	files, err := ioutil.ReadDir("images")
	if err != nil {
		ErrorLogger.Println("Cannot access images directory")
		return "", err
	}

	filtered := filter(files, func(info os.FileInfo) bool {
		return !info.IsDir() && strings.HasSuffix(info.Name(), ".jpg")
	})

	if len(filtered) == 0 {
		return "", errors.New("Images directory is empty")
	}

	_, min, _ := time.Now().Clock()
	duration := 60 / len(filtered)
	index := min / duration

	return filtered[index].Name(), nil
}

func filter(vs []os.FileInfo, f func(info os.FileInfo) bool) []os.FileInfo {
	vsf := make([]os.FileInfo, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
