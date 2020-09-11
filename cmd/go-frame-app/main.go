package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gitlab.com/go-displays/go-frame/cmd/go-frame-app/api"
	"log"
	"os"
)

//go:generate go-bindata-assetfs -o web-assets.go ../../web/dist/...

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

func main() {
	InfoLogger.Println("Starting up Go-Frame server...")
	InfoLogger.Println("Setting up HTTP endpoint")

	router := gin.Default()
	apiEndpoint := router.Group("/api")
	router.Use(static.ServeRoot("/static/images", "images"))
	router.Use(static.Serve("/", BinaryFileSystem("../../web/dist")))

	api.RegisterApiEndpoint(apiEndpoint)

	router.Run(":8080")
}
