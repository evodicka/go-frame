package adminapi

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/persistence"
)

// ImageRef represents an image object for the admin API.
type ImageRef struct {
	// Id is the unique identifier of the image.
	Id int `json:"id" binding:"required"`
	// Path is the filename of the image in the storage.
	Path string `json:"path" binding:"required"`
	// Type indicates the content type (e.g. IMAGE).
	Type model.Type `json:"type" binding:"required"`
	// Metadata stores optional metadata about the image.
	Metadata string `json:"metadata"`
}

func (h *Handler) loadAllImageData(context *gin.Context) {
	var images []ImageRef

	loadedImages, err := h.storage.LoadImages()
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

func (h *Handler) updateImageOrder(context *gin.Context) {
	var images []ImageRef
	if err := context.ShouldBindJSON(&images); err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var dbImages []model.Image
	for _, image := range images {
		var dbImage = model.Image{
			Id:       image.Id,
			Path:     image.Path,
			Type:     image.Type,
			Metadata: image.Metadata,
		}
		dbImages = append(dbImages, dbImage)
	}
	if err := h.storage.ReorderImages(dbImages); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	h.loadAllImageData(context)
}

func (h *Handler) deleteImage(context *gin.Context) {
	id := context.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = h.storage.DeleteImage(intId)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	context.Status(http.StatusOK)
}

func (h *Handler) addImage(context *gin.Context) {
	form, err := context.FormFile("image")
	if err != nil || form == nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = context.SaveUploadedFile(form, persistence.ImageDir+string(os.PathSeparator)+form.Filename)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	loadedImage, err := h.storage.SaveImageMetadata(form.Filename)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var image = ImageRef{
		Id:       loadedImage.Id,
		Path:     loadedImage.Path,
		Type:     loadedImage.Type,
		Metadata: loadedImage.Metadata,
	}

	context.JSON(http.StatusOK, image)
}
