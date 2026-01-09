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

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/persistence"
)

func TestMain(m *testing.M) {
	os.Mkdir("images", 0755)
	code := m.Run()
	if persistence.Db != nil {
		persistence.Db.Close()
	}
	os.Remove("my.db")
	os.RemoveAll("images")
	os.Exit(code)
}

func clearDB() {
	persistence.Db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("images"))
		tx.DeleteBucket([]byte("order"))
		tx.DeleteBucket([]byte("status"))
		return nil
	})
	persistence.Db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("images"))
		tx.CreateBucketIfNotExists([]byte("order"))
		b, _ := tx.CreateBucketIfNotExists([]byte("status"))
		if b.Get([]byte("status")) == nil {
			status := persistence.Status{CurrentImageId: -1}
			d, _ := json.Marshal(status)
			b.Put([]byte("status"), d)
		}

		kb, _ := tx.CreateBucketIfNotExists([]byte("configuration"))
		if kb.Get([]byte("config")) == nil {
			config := persistence.Config{ImageDuration: 60}
			d, _ := json.Marshal(config)
			kb.Put([]byte("config"), d)
		}
		return nil
	})
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	g := r.Group("/admin/api")
	RegisterApiEndpoint(g)
	return r
}

func TestUploadImage(t *testing.T) {
	clearDB()
	r := setupRouter()

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
	var img persistence.Image
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
	clearDB()
	persistence.SaveImageMetadata("img1.jpg")

	r := setupRouter()
	req, _ := http.NewRequest("GET", "/admin/api/image", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var images []persistence.Image
	json.Unmarshal(w.Body.Bytes(), &images)
	if len(images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(images))
	}
}

func TestDeleteImage(t *testing.T) {
	clearDB()
	// img is unused here, renamed to _
	_, _ = persistence.SaveImageMetadata("del.jpg")
	// Create file
	os.Create("images/del.jpg")

	// r is unused here? No we need it for ServeHTTP if we implement it.
	// But since I didn't implement the request yet...

	// Let's implement Delete logic test
	r := setupRouter()

	// We need the ID. Since we reset DB, img ID should be 1.
	req, _ := http.NewRequest("DELETE", "/admin/api/image/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Check file deleted
	if _, err := os.Stat("images/del.jpg"); !os.IsNotExist(err) {
		t.Error("File should be deleted")
	}
}

func TestConfiguration(t *testing.T) {
	clearDB()
	r := setupRouter()

	// Get Config
	req, _ := http.NewRequest("GET", "/admin/api/configuration", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET config failed: %d", w.Code)
	}

	// Update Config
	newConfig := persistence.Config{ImageDuration: 99, RandomOrder: true}
	body, _ := json.Marshal(newConfig)
	req, _ = http.NewRequest("PUT", "/admin/api/configuration", bytes.NewBuffer(body)) // PUT for update?
	// Check api_loader.go: router.PUT("/configuration", updateConfiguration)
	// path is /configuration (relative to group /admin/api -> /admin/api/configuration)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("PUT config failed: %d", w.Code)
	}

	config, _ := persistence.GetConfiguration()
	if config.ImageDuration != 99 {
		t.Errorf("Config not updated")
	}
}
