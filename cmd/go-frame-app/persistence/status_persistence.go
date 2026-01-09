package persistence

import (
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Status represents the runtime status of the frame (current image, last switch time).
type Status struct {
	// CurrentImageId is the ID of the currently displayed image.
	CurrentImageId int
	// LastSwitch is the timestamp when the image was last switched.
	LastSwitch time.Time
	// ImageDuration is legacy/unused field? (based on usage it seems unused in update, only loaded).
	ImageDuration int
}

const (
	// CurrentStatusKey is the key used to store the status object in the database.
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
	}
	statusBytes, _ := json.Marshal(status)
	return statusBucket.Put([]byte(CurrentStatusKey), statusBytes)
}

// GetCurrentStatus retrieves the current runtime status from the database.
//
// Returns:
//   - Status: The current status object.
//   - error: An error if retrieval fails.
func GetCurrentStatus() (Status, error) {
	var status Status
	err := Db.View(func(tx *bolt.Tx) error {
		statusBucket := tx.Bucket(statusBucketName)
		statusBytes := statusBucket.Get([]byte(CurrentStatusKey))
		return json.Unmarshal(statusBytes, &status)
	})
	return status, err
}

// UpdateImageStatus updates the current image ID and resets the switch timer.
//
// Parameters:
//   - newId: The ID of the image now being displayed.
//
// Returns:
//   - error: An error if the status update fails.
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
