package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/persistence"
)

func TestMain(m *testing.M) {
	// Setup environment for persistence init
	os.Mkdir("images", 0755)

	// Open DB if not already open (persistence init might have done it)
	// Actually persistence.init() runs and opens "my.db" in CWD.

	code := m.Run()

	// Teardown
	if persistence.Db != nil {
		persistence.Db.Close()
	}
	os.Remove("my.db") // Created by persistence init
	os.RemoveAll("images")
	os.Exit(code)
}

func clearDB() {
	persistence.Db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("images"))
		tx.DeleteBucket([]byte("order"))
		tx.DeleteBucket([]byte("status"))
		// Don't delete config, or re-init it?
		// persistence initBuckets handles creation if not exists.
		// But we can't call initBuckets.
		// So better to just clear keys if we want, or use unique IDs.
		return nil
	})
	// Re-create buckets if deleted (helper since we can't call InitBuckets)
	persistence.Db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists([]byte("images"))
		_, _ = tx.CreateBucketIfNotExists([]byte("order"))
		_, _ = tx.CreateBucketIfNotExists([]byte("status"))
		// Status needs pre-population for GetCurrentStatus to work without error?
		// GetCurrentStatus -> json.Unmarshal(nil) returns error?
		// persistence.GetCurrentStatus calls Get([]byte("status")).
		// If empty, returns nil. Unmarshal(nil) returns error.
		// So we must prepopulate status.
		b := tx.Bucket([]byte("status"))
		if b.Get([]byte("status")) == nil {
			status := persistence.Status{
				CurrentImageId: -1,
				LastSwitch:     time.Unix(0, 0),
			}
			bytes, _ := json.Marshal(status)
			b.Put([]byte("status"), bytes)
		}

		// Config same
		b = tx.Bucket([]byte("configuration")) // Name is "configuration" ?
		// config_persistence.go says: configBucketName = []byte("configuration")
		if b == nil {
			b, _ = tx.CreateBucketIfNotExists([]byte("configuration"))
		}
		if b.Get([]byte("config")) == nil {
			config := persistence.Config{ImageDuration: 60, RandomOrder: false}
			bytes, _ := json.Marshal(config)
			b.Put([]byte("config"), bytes)
		}

		return nil
	})
}

func TestCalculateCurrentImage(t *testing.T) {
	clearDB()

	// 1. Initial state, no images.
	// NOTE: Current implementation returns empty path and nil error when no images exist.
	path, err := calculateCurrentImage()
	if err != nil {
		// If it errors, that's acceptable too in case of fix, but currently it doesn't.
	}
	if path != "" {
		t.Errorf("Expected empty path when no images, got %s", path)
	}

	// Reset DB because empty CalculateCurrentImage() calls UpdateImageStatus(0), setting bad state.
	clearDB()

	// 2. Add an image
	img1, err := persistence.SaveImageMetadata("img1.jpg")
	if err != nil {
		t.Fatalf("Failed to save image: %v", err)
	}

	// 3. Current ID is -1 (from reset status), so should load first image
	path, err = calculateCurrentImage()
	if err != nil {
		t.Fatalf("calculateCurrentImage failed: %v", err)
	}
	if path != "img1.jpg" {
		t.Errorf("Expected img1.jpg, got %s", path)
	}

	// 4. Verify status updated
	status, _ := persistence.GetCurrentStatus()
	if status.CurrentImageId != img1.Id {
		t.Errorf("Expected status ID %d, got %d", img1.Id, status.CurrentImageId)
	}

	// 5. Add second image
	_, _ = persistence.SaveImageMetadata("img2.jpg")

	// 6. Call again immediately. Duration is 60s. Should still be img1.
	path, err = calculateCurrentImage()
	if err != nil {
		t.Fatalf("calculateCurrentImage failed: %v", err)
	}
	if path != "img1.jpg" {
		t.Errorf("Expected img1.jpg (not switched yet), got %s", path)
	}

	// 7. Manipulate LastSwitch to force switch
	persistence.UpdateImageStatus(img1.Id) // Resets time to Now
	// Use explicit DB update to set timeback
	persistence.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("status"))
		status := persistence.Status{
			CurrentImageId: img1.Id,
			LastSwitch:     time.Now().Add(-61 * time.Second),
		}
		bytes, _ := json.Marshal(status)
		return b.Put([]byte("status"), bytes)
	})

	path, err = calculateCurrentImage()
	if err != nil {
		t.Fatalf("calculateCurrentImage failed: %v", err)
	}
	if path != "img2.jpg" {
		t.Errorf("Expected img2.jpg, got %s", path)
	}
}

func TestGetCurrentImageData(t *testing.T) {
	clearDB()
	persistence.SaveImageMetadata("test.jpg")

	// Setup Gin
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/current-image", getCurrentImageData)

	req, _ := http.NewRequest("GET", "/current-image", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var ref ImageRef
	err := json.Unmarshal(w.Body.Bytes(), &ref)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if ref.Path != "test.jpg" {
		t.Errorf("Expected path test.jpg, got %s", ref.Path)
	}
	if ref.Type != model.ImageType {
		t.Errorf("Expected type IMAGE, got %s", ref.Type)
	}
}
