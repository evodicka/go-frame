package persistence

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"time"
)

type Status struct {
	CurrentImageId int
	LastSwitch     time.Time
	ImageDuration  int
}

const (
	CurrentStatusKey string = "status"
)

var statusBucketName = []byte("status")

func initStatusBuckets(tx *bolt.Tx) error {
	statusBucket, err := tx.CreateBucketIfNotExists(statusBucketName)
	if err != nil {
		return err
	}
	if isBucketEmpty(statusBucket) {
		err = prepopulateStatus(statusBucket)
	}
	return err
}

func prepopulateStatus(statusBucket *bolt.Bucket) error {
	var status = Status{
		CurrentImageId: -1,
		LastSwitch:     time.Unix(0, 0),
		ImageDuration:  300,
	}
	statusBytes, _ := json.Marshal(status)
	return statusBucket.Put([]byte(CurrentStatusKey), statusBytes)
}

func GetCurrentStatus() (Status, error) {
	var status Status
	err := Db.View(func(tx *bolt.Tx) error {
		statusBucket := tx.Bucket(statusBucketName)
		statusBytes := statusBucket.Get([]byte(CurrentStatusKey))
		return json.Unmarshal(statusBytes, &status)
	})
	return status, err
}

func UpdateImageStatus(newId int) error {
	return Db.Update(func(tx *bolt.Tx) error {
		statusBucket := tx.Bucket(statusBucketName)
		statusBytes := statusBucket.Get([]byte(CurrentStatusKey))
		var status Status
		if err := json.Unmarshal(statusBytes, &status); err != nil {
			return err
		}
		status.CurrentImageId = newId
		status.LastSwitch = time.Now()
		newStatus, _ := json.Marshal(status)
		return statusBucket.Put([]byte(CurrentStatusKey), newStatus)
	})
}
