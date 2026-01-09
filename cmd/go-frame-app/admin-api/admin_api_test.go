package adminapi

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
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

	dbPath := filepath.Join(t.TempDir(), "test_admin_api.db")
	storage, err := persistence.NewStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	t.Cleanup(func() {
		storage.Close()
		os.RemoveAll("images")
	})

	// Prepopulate status
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

func setupRouter(storage *persistence.Storage) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	g := r.Group("/admin/api")
	handler := NewHandler(storage)
	handler.RegisterApiEndpoint(g)
	return r
}

func TestUploadImage(t *testing.T) {
	storage := setupTestDB(t)
	r := setupRouter(storage)

	// Create multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test_upload.jpg")
	part.Write([]byte("fake image content"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/admin/api/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify ID returned
	var img model.Image
	json.Unmarshal(w.Body.Bytes(), &img)
	if img.Path != "test_upload.jpg" {
		t.Errorf("Expected path test_upload.jpg, got %s", img.Path)
	}

	// Verify file exists on disk
	if _, err := os.Stat(filepath.Join("images", "test_upload.jpg")); os.IsNotExist(err) {
		t.Error("File not saved to disk")
	}
}

func TestGetImages(t *testing.T) {
	storage := setupTestDB(t)
	storage.SaveImageMetadata("img1.jpg")

	r := setupRouter(storage)
	req, _ := http.NewRequest("GET", "/admin/api/image", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var images []model.Image
	json.Unmarshal(w.Body.Bytes(), &images)
	if len(images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(images))
	}
}

func TestDeleteImage(t *testing.T) {
	storage := setupTestDB(t)
	// img is unused here, renamed to _
	_, _ = storage.SaveImageMetadata("del.jpg")
	// Create file
	os.Create("images/del.jpg")

	r := setupRouter(storage)

	req, _ := http.NewRequest("DELETE", "/admin/api/image/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		// Just fail it if it's not 200, but "1" might be wrong ID.
		// If "1" is wrong, it returns 500 or 400.
		// Assuming 1 since new DB starts sequence at 1.
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Check file deleted
	if _, err := os.Stat("images/del.jpg"); !os.IsNotExist(err) {
		t.Error("File should be deleted")
	}
}

func TestConfiguration(t *testing.T) {
	storage := setupTestDB(t)
	r := setupRouter(storage)

	// Get Config
	req, _ := http.NewRequest("GET", "/admin/api/configuration", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET config failed: %d", w.Code)
	}

	// Update Config
	newConfig := model.Config{ImageDuration: 99, RandomOrder: true}
	body, _ := json.Marshal(newConfig)
	req, _ = http.NewRequest("PUT", "/admin/api/configuration", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("PUT config failed: %d", w.Code)
	}

	config, _ := storage.GetConfiguration()
	if config.ImageDuration != 99 {
		t.Errorf("Config not updated")
	}
}
