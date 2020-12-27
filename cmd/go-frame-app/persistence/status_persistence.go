package persistence

import bolt "go.etcd.io/bbolt"

const (
	CurrentImageKey    string = "currentImage"
	LastImageSwitchKey string = "lastImageSwitch"
	ImageDurationKey   string = "imageDuration"
)

func initStatusBuckets(tx *bolt.Tx) error {
	statusBucket, err := tx.CreateBucketIfNotExists([]byte("status"))
	if err != nil {
		return err
	}
	if isBucketEmpty(statusBucket) {
		err = prepopulateStatus(statusBucket)
	}
	return err
}

func prepopulateStatus(statusBucket *bolt.Bucket) error {
	err := statusBucket.Put([]byte(CurrentImageKey), itob(-1))
	err = statusBucket.Put([]byte(LastImageSwitchKey), i64tob(0))
	err = statusBucket.Put([]byte(ImageDurationKey), itob(300))
	return err
}
