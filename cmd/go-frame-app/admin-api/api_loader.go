package adminapi

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
}

func RegisterApiEndpoint(router *gin.RouterGroup) {
	router.GET("/image", loadAllImageData)
	router.PUT("/image", updateImageOrder)
	router.DELETE("/image/:id", deleteImage)
	router.POST("/image", addImage)
}
