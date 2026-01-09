package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/persistence"
)

func setupTestDB(t *testing.T) *persistence.Storage {
	_ = os.MkdirAll("images", 0755)

	dbPath := filepath.Join(t.TempDir(), "test_api.db")
	storage, err := persistence.NewStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	// Helper to clear DB if needed or just return fresh one.
	// Since we use TempDir, we get fresh one each time setupTestDB is called if we call it per test.

	t.Cleanup(func() {
		storage.Close()
		os.RemoveAll("images")
	})

	// Prepopulate status for tests
	err = storage.Db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte("status"))
		if b.Get([]byte("status")) == nil {
			status := model.Status{
				CurrentImageId: -1,
				LastSwitch:     time.Unix(0, 0),
			}
			bytes, _ := json.Marshal(status)
			b.Put([]byte("status"), bytes)
		}

		kb, _ := tx.CreateBucketIfNotExists([]byte("configuration"))
		if kb.Get([]byte("config")) == nil {
			config := model.Config{ImageDuration: 60, RandomOrder: false}
			bytes, _ := json.Marshal(config)
			kb.Put([]byte("config"), bytes)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to prepopulate DB: %v", err)
	}

	return storage
}

func TestCalculateCurrentImage(t *testing.T) {
	storage := setupTestDB(t)
	handler := NewHandler(storage)

	// 1. Initial state, no images.
	path, err := handler.calculateCurrentImage()
	if err != nil {
		// If it errors, that's acceptable too in case of fix, but currently it doesn't.
	}
	if path != "" {
		t.Errorf("Expected empty path when no images, got %s", path)
	}

	// 2. Add an image
	img1, err := storage.SaveImageMetadata("img1.jpg")
	if err != nil {
		t.Fatalf("Failed to save image: %v", err)
	}

	// 3. Current ID is -1 (from reset status), so should load first image
	path, err = handler.calculateCurrentImage()
	if err != nil {
		t.Fatalf("calculateCurrentImage failed: %v", err)
	}
	if path != "img1.jpg" {
		t.Errorf("Expected img1.jpg, got %s", path)
	}

	// 4. Verify status updated
	status, _ := storage.GetCurrentStatus()
	if status.CurrentImageId != img1.Id {
		t.Errorf("Expected status ID %d, got %d", img1.Id, status.CurrentImageId)
	}

	// 5. Add second image
	_, _ = storage.SaveImageMetadata("img2.jpg")

	// 6. Call again immediately. Duration is 60s. Should still be img1.
	path, err = handler.calculateCurrentImage()
	if err != nil {
		t.Fatalf("calculateCurrentImage failed: %v", err)
	}
	if path != "img1.jpg" {
		t.Errorf("Expected img1.jpg (not switched yet), got %s", path)
	}

	// 7. Manipulate LastSwitch to force switch
	storage.UpdateImageStatus(img1.Id) // Resets time to Now
	// Use explicit DB update to set timeback
	storage.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("status"))
		status := model.Status{
			CurrentImageId: img1.Id,
			LastSwitch:     time.Now().Add(-61 * time.Second),
		}
		bytes, _ := json.Marshal(status)
		return b.Put([]byte("status"), bytes)
	})

	path, err = handler.calculateCurrentImage()
	if err != nil {
		t.Fatalf("calculateCurrentImage failed: %v", err)
	}
	if path != "img2.jpg" {
		t.Errorf("Expected img2.jpg, got %s", path)
	}
}

func TestGetCurrentImageData(t *testing.T) {
	storage := setupTestDB(t)
	storage.SaveImageMetadata("test.jpg")

	handler := NewHandler(storage)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler.RegisterApiEndpoint(r.Group("/"))
	// RegisterApiEndpoint registers /image/current on the group.
	// if group is /, then /image/current.

	// Wait, RegisterApiEndpoint calls router.GET("/image/current", ...)
	// So we request /image/current

	req, _ := http.NewRequest("GET", "/image/current", nil)
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
