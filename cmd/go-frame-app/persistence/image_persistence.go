package persistence

import (
	"bytes"
	"encoding/json"
	"errors"
	"gitlab.com/go-displays/go-frame/cmd/go-frame-app/model"
	bolt "go.etcd.io/bbolt"
	"os"
)

type Image struct {
	Id       int
	Path     string
	Type     model.Type
	Metadata string
}

const (
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

func LoadImages() ([]Image, error) {
	var images []Image
	err := Db.View(func(tx *bolt.Tx) error {
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

func LoadImage(id int) (Image, error) {
	var image Image
	err := Db.View(func(tx *bolt.Tx) error {
		metadataBucket := tx.Bucket(metadataBucketName)
		returnedImage, err := loadImageByByteId(itob(id), metadataBucket)
		image = returnedImage
		return err
	})
	return image, err
}

func LoadNextImage(id int) (Image, error) {
	var image Image
	err := Db.View(func(tx *bolt.Tx) error {
		orderBucket := tx.Bucket(orderBucketName)
		metadataBucket := tx.Bucket(metadataBucketName)
		cursor := orderBucket.Cursor()
		if id < 0 {
			_, value := cursor.First()
			image, _ = loadImageByByteId(value, metadataBucket)
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

func loadImageByByteId(id []byte, metadataBucket *bolt.Bucket) (Image, error) {
	image := metadataBucket.Get(id)
	if image != nil {
		var imageStruct Image
		err := json.Unmarshal(image, &imageStruct)
		if err != nil {
			return Image{}, err
		} else {
			return imageStruct, nil
		}
	}
	return Image{}, errors.New("Image not found")
}

func ReorderImages(images []Image) error {
	var sequences []int
	for _, image := range images {
		sequences = append(sequences, image.Id)
	}
	return Db.Update(func(tx *bolt.Tx) error {
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

func DeleteImage(id int) error {
	return Db.Update(func(tx *bolt.Tx) error {
		metadataBucket := tx.Bucket(metadataBucketName)
		metadata := metadataBucket.Get(itob(id))
		if metadata == nil {
			return errors.New("Image not found")
		}
		var image Image
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

func SaveImageMetadata(name string) (Image, error) {
	var image Image
	err := Db.Update(func(tx *bolt.Tx) error {
		orderBucket := tx.Bucket(orderBucketName)
		metadataBucket := tx.Bucket(metadataBucketName)

		sequence, err := metadataBucket.NextSequence()
		if err != nil {
			return err
		}
		image = Image{
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
