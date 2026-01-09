package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/persistence"
)

// ImageRef represents a reference to an image or content to be displayed.
type ImageRef struct {
	// Path is the filename or URL of the content.
	Path string `json:"path" binding:"required"`
	// Type indicates the type of content (e.g. IMAGE or URL).
	Type model.Type `json:"type" binding:"required"`
	// Metadata contains additional information about the image (optional).
	Metadata string `json:"metadata"`
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
	config, err := persistence.GetConfiguration()
	if err != nil {
		ErrorLogger.Println("Cannot read configuration")
		return "", err
	}
	var image persistence.Image
	if time.Since(status.LastSwitch).Seconds() > float64(config.ImageDuration) {
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
