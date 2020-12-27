package persistence

import (
	"encoding/json"
	"gitlab.com/go-displays/go-frame/cmd/go-frame-app/model"
	bolt "go.etcd.io/bbolt"
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
