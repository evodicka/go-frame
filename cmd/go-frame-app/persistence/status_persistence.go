package persistence

import (
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
)

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
	var status = model.Status{
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
func (s *Storage) GetCurrentStatus() (model.Status, error) {
	var status model.Status
	err := s.Db.View(func(tx *bolt.Tx) error {
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
func (s *Storage) UpdateImageStatus(newId int) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		statusBucket := tx.Bucket(statusBucketName)
		statusBytes := statusBucket.Get([]byte(CurrentStatusKey))
		var status model.Status
		if err := json.Unmarshal(statusBytes, &status); err != nil {
			return err
		}
		status.CurrentImageId = newId
		status.LastSwitch = time.Now()
		newStatus, _ := json.Marshal(status)
		return statusBucket.Put([]byte(CurrentStatusKey), newStatus)
	})
}
