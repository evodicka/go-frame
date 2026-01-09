package persistence

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

// Config represents the application configuration stored in the database.
type Config struct {
	// ImageDuration is the time in seconds each image is displayed.
	ImageDuration int
	// RandomOrder toggles random image shuffling.
	RandomOrder bool
}

const (
	// ConfigKey is the key, used to store the configuration object in the database.
	ConfigKey = "config"
)

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

// GetConfiguration retrieves the current application configuration.
//
// Returns:
//   - Config: The current configuration object.
//   - error: An error if retrieval fails.
func GetConfiguration() (Config, error) {
	var config Config
	err := Db.View(func(tx *bolt.Tx) error {
		configBucket := tx.Bucket(configBucketName)
		configBytes := configBucket.Get([]byte(ConfigKey))
		return json.Unmarshal(configBytes, &config)
	})
	return config, err
}

// UpdateConfiguration persists a new configuration to the database.
//
// Parameters:
//   - config: The new configuration object to save.
//
// Returns:
//   - error: An error if the update fails.
func UpdateConfiguration(config Config) error {
	return Db.Update(func(tx *bolt.Tx) error {
		configBucket := tx.Bucket(configBucketName)
		conigBytes, _ := json.Marshal(config)
		return configBucket.Put([]byte(ConfigKey), conigBytes)
	})
}
