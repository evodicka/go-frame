package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	adminapi "go.evodicka.dev/go-frame/cmd/go-frame-app/admin-api"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/api"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/persistence"
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

func main() {
	InfoLogger.Println("Starting up Go-Frame server...")
	InfoLogger.Println("Setting up HTTP endpoint")

	router := gin.Default()
	apiEndpoint := router.Group("/api")
	adminEndpoint := router.Group("/admin/api")
	router.Use(static.ServeRoot("/static/images", "images"))
	router.Use(static.Serve("/", EmbeddedFileSystem("web")))

	api.RegisterApiEndpoint(apiEndpoint)
	adminapi.RegisterApiEndpoint(adminEndpoint)

	router.Run(":8080")

	defer persistence.Close()
}
