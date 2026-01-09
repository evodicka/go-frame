package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
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

func (h *Handler) getCurrentImageData(context *gin.Context) {
	fileName, err := h.calculateCurrentImage()
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ref := ImageRef{Path: fileName, Type: model.ImageType}
	context.JSON(http.StatusOK, ref)
}

func (h *Handler) calculateCurrentImage() (path string, err error) {
	status, err := h.storage.GetCurrentStatus()
	if err != nil {
		ErrorLogger.Println("Cannot read current status")
		return "", err
	}
	config, err := h.storage.GetConfiguration()
	if err != nil {
		ErrorLogger.Println("Cannot read configuration")
		return "", err
	}
	var image model.Image
	if time.Since(status.LastSwitch).Seconds() > float64(config.ImageDuration) {
		image, err = h.storage.LoadNextImage(status.CurrentImageId)
		if err == nil {
			// Update status with new image ID
			err = h.storage.UpdateImageStatus(image.Id)
		}
	} else {
		image, err = h.storage.LoadImage(status.CurrentImageId)
	}
	if err != nil {
		ErrorLogger.Println("Cannot read Image")
		return "", err
	}
	return image.Path, err
}
