package model

// Type represents the content type of a frame item (e.g. Image or URL).
type Type string

const (
	// ImageType indicates that the item is a local image file.
	ImageType Type = "IMAGE"
	// Url indicates that the item is a remote URL.
	Url Type = "URL"
)
