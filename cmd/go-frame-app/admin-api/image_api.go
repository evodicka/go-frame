package adminapi

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/go-displays/go-frame/cmd/go-frame-app/model"
	"gitlab.com/go-displays/go-frame/cmd/go-frame-app/persistence"
	"net/http"
)

type ImageRef struct {
	Id       int        `json:"id" binding:"required"`
	Path     string     `json:"path" binding:"required"`
	Type     model.Type `json:"type" binding:"required"`
	Metadata string     `json:"metadata"`
}

func loadAllImageData(context *gin.Context) {
	var images []ImageRef

	loadedImages, err := persistence.LoadImages()
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, loadedImage := range loadedImages {
		var image = ImageRef{
			Id:       loadedImage.Id,
			Path:     loadedImage.Path,
			Type:     loadedImage.Type,
			Metadata: loadedImage.Metadata,
		}
		images = append(images, image)
	}
	context.JSON(http.StatusOK, images)
}

func updateImageOrder(context *gin.Context) {
	var images []ImageRef
	if err := context.ShouldBindJSON(&images); err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var dbImages []persistence.Image
	for _, image := range images {
		var dbImage = persistence.Image{
			Id:       image.Id,
			Path:     image.Path,
			Type:     image.Type,
			Metadata: image.Metadata,
		}
		dbImages = append(dbImages, dbImage)
	}
	if err := persistence.ReorderImages(dbImages); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	loadAllImageData(context)
}
