package persistence

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	bolt "go.etcd.io/bbolt"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
)

const (
	// ImageDir is the directory where image files are stored on disk.
	ImageDir string = "images"
)

var orderBucketName = []byte("order")
var metadataBucketName = []byte("images")

func initImageBuckets(tx *bolt.Tx) error {
	metadataBucket, err := tx.CreateBucketIfNotExists(metadataBucketName)
	if err != nil {
		return err
	}
	orderBucket, err := tx.CreateBucketIfNotExists(orderBucketName)
	if err != nil {
		return err
	}

	if isBucketEmpty(metadataBucket) {
		err = prepopulateImages(metadataBucket, orderBucket)
	}
	return err
}

// LoadImages retrieves all images from the database, ordered by their sequence.
//
// Returns:
//   - []Image: A slice of Image objects.
//   - error: An error if the database read fails.
func (s *Storage) LoadImages() ([]model.Image, error) {
	var images []model.Image
	err := s.Db.View(func(tx *bolt.Tx) error {
		orderBucket := tx.Bucket(orderBucketName)
		metadataBucket := tx.Bucket(metadataBucketName)

		err := orderBucket.ForEach(func(key, value []byte) error {
			image, err := loadImageByByteId(value, metadataBucket)
			if err == nil {
				images = append(images, image)
			}
			return nil
		})
		return err
	})
	return images, err
}

// LoadImage retrieves a specific image by its ID.
//
// Parameters:
//   - id: The ID of the image to retrieve.
//
// Returns:
//   - Image: The requested Image object.
//   - error: An error if the image is not found or database read fails.
func (s *Storage) LoadImage(id int) (model.Image, error) {
	var image model.Image
	err := s.Db.View(func(tx *bolt.Tx) error {
		metadataBucket := tx.Bucket(metadataBucketName)
		returnedImage, err := loadImageByByteId(itob(id), metadataBucket)
		image = returnedImage
		return err
	})
	return image, err
}

// LoadNextImage determines and retrieves the next image to be displayed based on the current image ID.
// It cycles through the images in the defined order.
//
// Parameters:
//   - id: The ID of the currently displayed image.
//
// Returns:
//   - Image: The next Image object to display.
//   - error: An error if the next image cannot be determined or loaded.
func (s *Storage) LoadNextImage(id int) (model.Image, error) {
	var image model.Image
	err := s.Db.View(func(tx *bolt.Tx) error {
		orderBucket := tx.Bucket(orderBucketName)
		metadataBucket := tx.Bucket(metadataBucketName)
		cursor := orderBucket.Cursor()
		if id < 0 {
			_, value := cursor.First()
			if value == nil {
				return errors.New("No images found")
			}
			img, err := loadImageByByteId(value, metadataBucket)
			if err != nil {
				return err
			}
			image = img
		} else {
			for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
				if bytes.Equal(value, itob(id)) {
					k, v := cursor.Next()
					if k == nil {
						k, v = cursor.First()
					}
					image, _ = loadImageByByteId(v, metadataBucket)
					break
				}
			}
		}
		return nil
	})
	return image, err

}

func loadImageByByteId(id []byte, metadataBucket *bolt.Bucket) (model.Image, error) {
	image := metadataBucket.Get(id)
	if image != nil {
		var imageStruct model.Image
		err := json.Unmarshal(image, &imageStruct)
		if err != nil {
			return model.Image{}, err
		} else {
			return imageStruct, nil
		}
	}
	return model.Image{}, errors.New("Image not found")
}

// ReorderImages updates the display order of images in the database.
//
// Parameters:
//   - images: A slice of Image objects in the desired order.
//
// Returns:
//   - error: An error if the database update fails.
func (s *Storage) ReorderImages(images []model.Image) error {
	var sequences []int
	for _, image := range images {
		sequences = append(sequences, image.Id)
	}
	return s.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(orderBucketName)
		return persistImageOrder(bucket, sequences)
	})
}

func persistImageOrder(orderBucket *bolt.Bucket, sequences []int) error {
	err := orderBucket.ForEach(func(key, value []byte) error {
		return orderBucket.Delete(key)
	})
	if err != nil {
		return err
	}

	for i, sequence := range sequences {
		err = orderBucket.Put(itob(i), itob(sequence))
	}
	return err
}

// DeleteImage removals an image from the database and the filesystem.
//
// Parameters:
//   - id: The ID of the image to delete.
//
// Returns:
//   - error: An error if the image is not found or deletion fails.
func (s *Storage) DeleteImage(id int) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		metadataBucket := tx.Bucket(metadataBucketName)
		metadata := metadataBucket.Get(itob(id))
		if metadata == nil {
			return errors.New("Image not found")
		}
		var image model.Image
		if err := json.Unmarshal(metadata, &image); err != nil {
			return err
		}
		if err := deleteImageOnDisk(image.Path); err != nil {
			return err
		}
		return metadataBucket.Delete(itob(id))
	})
}

func deleteImageOnDisk(path string) error {
	var filename = ImageDir + string(os.PathSeparator) + path
	return os.Remove(filename)
}

// SaveImageMetadata creates a new image entry in the database.
//
// Parameters:
//   - name: The filename of the image.
//
// Returns:
//   - Image: The created Image object with assigned ID.
//   - error: An error if the database/metadata update fails.
func (s *Storage) SaveImageMetadata(name string) (model.Image, error) {
	var image model.Image
	err := s.Db.Update(func(tx *bolt.Tx) error {
		orderBucket := tx.Bucket(orderBucketName)
		metadataBucket := tx.Bucket(metadataBucketName)

		sequence, err := metadataBucket.NextSequence()
		if err != nil {
			return err
		}
		image = model.Image{
			Id:       int(sequence),
			Path:     name,
			Type:     model.ImageType,
			Metadata: "",
		}
		imageJson, _ := json.Marshal(image)
		err = metadataBucket.Put(itob(int(sequence)), imageJson)
		if err != nil {
			return err
		}
		order := orderBucket.Stats().KeyN
		return orderBucket.Put(itob(order), itob(int(sequence)))
	})
	return image, err
}
