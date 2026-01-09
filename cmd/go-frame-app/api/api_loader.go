package api

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

// RegisterApiEndpoint registers the public API endpoints on the provided router group.
// It sets up the route for getting the current image.
//
// Parameters:
//   - router: The Gin router group to attach the endpoints to.
func RegisterApiEndpoint(router *gin.RouterGroup) {
	router.GET("/image/current", getCurrentImageData)
}
