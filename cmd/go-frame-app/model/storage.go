package model

type ImageStorage interface {
	// Status Operations
	GetCurrentStatus() (Status, error)
	UpdateImageStatus(newId int) error

	// Configuration Operations
	GetConfiguration() (Config, error)
	UpdateConfiguration(config Config) error

	// Image Operations
	LoadImage(id int) (Image, error)
	LoadNextImage(id int) (Image, error)
}

type ConfigurationAdminStorage interface {
	// Configuration Operations
	GetConfiguration() (Config, error)
	UpdateConfiguration(config Config) error
}

type ImageAdminStorage interface {
	// Image Operations
	LoadImages() ([]Image, error)
	ReorderImages(images []Image) error
	DeleteImage(id int) error
	SaveImageMetadata(name string) (Image, error)
}

type AdminStorage interface {
	ConfigurationAdminStorage
	ImageAdminStorage
}
