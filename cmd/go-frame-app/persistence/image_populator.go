package persistence

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	bolt "go.etcd.io/bbolt"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
)

func prepopulateImages(metadataBucket *bolt.Bucket, orderBucket *bolt.Bucket) error {
	sequences, err := persistImagesFromDir(metadataBucket)
	if err != nil {
		return err
	}
	err = persistImageOrder(orderBucket, sequences)
	return err
}

func persistImagesFromDir(metadataBucket *bolt.Bucket) ([]int, error) {
	InfoLogger.Println("Loading images into database")
	// Harden: Ensure directory exists
	if err := os.MkdirAll(ImageDir, 0755); err != nil {
		return nil, err
	}
	files, err := os.ReadDir(ImageDir)
	if err != nil {
		return nil, err
	}
	filtered := filter(files, func(info os.DirEntry) bool {
		return !info.IsDir() && strings.HasSuffix(info.Name(), ".jpg")
	})
	InfoLogger.Println("Found " + strconv.Itoa(len(filtered)) + " images to save")
	var sequences []int
	for _, imageInfo := range filtered {
		sequence, err := metadataBucket.NextSequence()
		if err != nil {
			return nil, err
		}
		image := Image{
			Id:   int(sequence),
			Path: imageInfo.Name(),
			Type: model.ImageType,
		}
		imageJson, _ := json.Marshal(image)
		err = metadataBucket.Put(itob(int(sequence)), imageJson)
		if err != nil {
			return nil, err
		}
		sequences = append(sequences, int(sequence))
	}
	InfoLogger.Println("Persisted all image metadata")
	return sequences, nil
}

func filter(vs []os.DirEntry, f func(info os.DirEntry) bool) []os.DirEntry {
	vsf := make([]os.DirEntry, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
