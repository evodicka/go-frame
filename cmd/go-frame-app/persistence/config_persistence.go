package persistence

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
)

type Config struct {
	ImageDuration int
	RandomOrder   bool
}

const ConfigKey = "config"

var configBucketName = []byte("configuration")

func initConfigBucket(tx *bolt.Tx) error {
	configBucket, err := tx.CreateBucketIfNotExists(configBucketName)
	if err != nil {
		return err
	}
	if isBucketEmpty(configBucket) {
		err = prepopulateConfiguration(configBucket)
	}
	return err
}

func prepopulateConfiguration(bucket *bolt.Bucket) error {
	var config = Config{
		ImageDuration: 60,
		RandomOrder:   false,
	}
	configBytes, _ := json.Marshal(config)
	return bucket.Put([]byte(ConfigKey), configBytes)
}

func GetConfiguration() (Config, error) {
	var config Config
	err := Db.View(func(tx *bolt.Tx) error {
		configBucket := tx.Bucket(configBucketName)
		configBytes := configBucket.Get([]byte(ConfigKey))
		return json.Unmarshal(configBytes, &config)
	})
	return config, err
}

func UpdateConfiguration(config Config) error {
	return Db.Update(func(tx *bolt.Tx) error {
		configBucket := tx.Bucket(configBucketName)
		conigBytes, _ := json.Marshal(config)
		return configBucket.Put([]byte(ConfigKey), conigBytes)
	})
}
