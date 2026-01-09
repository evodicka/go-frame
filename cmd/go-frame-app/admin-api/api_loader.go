package adminapi

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
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

// Handler holds dependencies for the admin API.
type Handler struct {
	storage model.AdminStorage
}

// NewHandler creates a new admin API handler with the given storage.
func NewHandler(storage model.AdminStorage) *Handler {
	return &Handler{storage: storage}
}

// RegisterApiEndpoint registers the admin API endpoints on the provided router group.
// It sets up routes for image management (CRUD) and configuration.
//
// Parameters:
//   - router: The Gin router group to attach the endpoints to.
func (h *Handler) RegisterApiEndpoint(router *gin.RouterGroup) {
	router.GET("/image", h.loadAllImageData)
	router.PUT("/image", h.updateImageOrder)
	router.DELETE("/image/:id", h.deleteImage)
	router.POST("/image", h.addImage)
	router.GET("/configuration", h.loadConfiguration)
	router.PUT("/configuration", h.updateConfiguration)
}
