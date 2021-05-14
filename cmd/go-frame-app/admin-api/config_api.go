package adminapi

import (
	"github.com/gin-gonic/gin"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/persistence"
	"net/http"
)

type ConfigRef struct {
	ImageDuration int  `json:"imageDuration" binding:"required"`
	RandomOrder   bool `json:"randomOrder"`
}

func loadConfiguration(context *gin.Context) {
	loadedConfig, err := persistence.GetConfiguration()
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var config = ConfigRef{
		ImageDuration: loadedConfig.ImageDuration,
		RandomOrder:   loadedConfig.RandomOrder,
	}
	context.JSON(http.StatusOK, config)
}

func updateConfiguration(context *gin.Context) {
	var config ConfigRef
	if err := context.ShouldBindJSON(&config); err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if config.ImageDuration < 0 {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var dbConfig = persistence.Config{
		ImageDuration: config.ImageDuration,
		RandomOrder:   config.RandomOrder,
	}

	if err := persistence.UpdateConfiguration(dbConfig); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	loadConfiguration(context)
}
