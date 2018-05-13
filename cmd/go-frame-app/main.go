package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"os"
)
//go:generate go-bindata-assetfs -o web-assets.go ../../web/dist/...

var logger = logging.MustGetLogger("main")

func init() {
	format := logging.MustStringFormatter(
		`%{color}[GO-FRAME] %{time:2006/02/01 - 15:04:05.000} %{level:.5s} %{id:03x} %{shortfunc}%{color:reset} ▶ %{message}`,
	)

	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendErr := logging.NewLogBackend(os.Stderr, "", 0)

	backendFormatted := logging.NewBackendFormatter(backend, format)
	backendLeveld := logging.AddModuleLevel(backendErr)
	backendLeveld.SetLevel(logging.ERROR, "")

	logging.SetBackend(backendFormatted, backendLeveld)
}


func main() {
	logger.Info("Starting up Go-Frame server...")
	logger.Info("Setting up HTTP endpoint")

	router := gin.Default()
	router.Use(static.ServeRoot("/static/images", "images"))
	router.Use(static.Serve("/", BinaryFileSystem("../../web/dist")))

	api := router.Group("/api")

	registerApiEndpoint(api)

	router.Run(":8080")
}

