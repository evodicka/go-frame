package persistence

import (
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

func initImageBuckets(tx *bolt.Tx) error {
	metadataBucket, err := tx.CreateBucketIfNotExists([]byte("images"))
	if err != nil {
		return err
	}
	orderBucket, err := tx.CreateBucketIfNotExists([]byte("order"))
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
		orderBucket := tx.Bucket([]byte("order"))
		metadataBucket := tx.Bucket([]byte("images"))

		err := orderBucket.ForEach(func(key, value []byte) error {
			image := metadataBucket.Get(value)
			if image != nil {
				var imageStruct Image
				err := json.Unmarshal(image, &imageStruct)
				if err == nil {
					images = append(images, imageStruct)
				} else {
					return err
				}
			}
			return nil
		})

		return err
	})
	return images, err
}

func ReorderImages(images []Image) error {
	var sequences []int
	for _, image := range images {
		sequences = append(sequences, image.Id)
	}
	return Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("order"))
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
		metadataBucket := tx.Bucket([]byte("images"))
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
		orderBucket := tx.Bucket([]byte("order"))
		metadataBucket := tx.Bucket([]byte("images"))

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
