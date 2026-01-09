package model

import "time"

// Type represents the content type of a frame item (e.g. Image or URL).
type Type string

const (
	// ImageType indicates that the item is a local image file.
	ImageType Type = "IMAGE"
	// Url indicates that the item is a remote URL.
	Url Type = "URL"
)

// Image represents the metadata of an image stored in the database.
type Image struct {
	// Id is the unique identifier of the image.
	Id int
	// Path is the filename of the image.
	Path string
	// Type indicates the media type (e.g. IMAGE).
	Type Type
	// Metadata contains additional info about the image.
	Metadata string
}

// Config represents the application configuration stored in the database.
type Config struct {
	// ImageDuration is the time in seconds each image is displayed.
	ImageDuration int
	// RandomOrder toggles random image shuffling.
	RandomOrder bool
}

// Status represents the runtime status of the frame (current image, last switch time).
type Status struct {
	// CurrentImageId is the ID of the currently displayed image.
	CurrentImageId int
	// LastSwitch is the timestamp when the image was last switched.
	LastSwitch time.Time
}
