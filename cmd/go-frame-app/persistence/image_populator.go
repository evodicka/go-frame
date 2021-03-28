package persistence

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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
	files, err := ioutil.ReadDir(ImageDir)
	if err != nil {
		return nil, err
	}
	filtered := filter(files, func(info os.FileInfo) bool {
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

func filter(vs []os.FileInfo, f func(info os.FileInfo) bool) []os.FileInfo {
	vsf := make([]os.FileInfo, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
