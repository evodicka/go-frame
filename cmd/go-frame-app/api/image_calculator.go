package api

import (
	"github.com/gin-gonic/gin"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/persistence"
	"net/http"
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
	status, err := persistence.GetCurrentStatus()
	if err != nil {
		ErrorLogger.Println("Cannot read current status")
		return "", err
	}
	var image persistence.Image
	if time.Since(status.LastSwitch).Seconds() > float64(status.ImageDuration) {
		image, err = persistence.LoadNextImage(status.CurrentImageId)
		if err == nil {
			err = persistence.UpdateImageStatus(image.Id)
		}
	} else {
		image, err = persistence.LoadImage(status.CurrentImageId)
	}
	if err != nil {
		ErrorLogger.Println("Cannot read Image")
		return "", err
	}
	return image.Path, err
}
