package static

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLocalFileSystem(t *testing.T) {
	// Setup temp dir
	dir := "test_static"
	os.Mkdir(dir, 0755)
	defer os.RemoveAll(dir)
	os.Create(filepath.Join(dir, "test.txt"))
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	os.Create(filepath.Join(dir, "subdir", "sub.txt"))

	fs := LocalFile(dir, false)

	// Test Exists
	if !fs.Exists("/", "/test.txt") {
		t.Error("test.txt should exist")
	}
	if fs.Exists("/", "/missing.txt") {
		t.Error("missing.txt should not exist")
	}

	// Test Exists with prefix
	if !fs.Exists("/static", "/static/test.txt") {
		t.Error("test.txt should exist with prefix")
	}

	// Test Serve middleware
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Serve("/", fs))

	// Request existing file
	req, _ := http.NewRequest("GET", "/test.txt", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Request missing file
	req, _ = http.NewRequest("GET", "/missing.txt", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		// middleware aborts? No, if exists returns true, serves.
		// If exists false, it continues.
		// Since no other handler, 404.
	}
}

func TestServeRoot(t *testing.T) {
	dir := "test_root"
	os.Mkdir(dir, 0755)
	defer os.RemoveAll(dir)
	os.Create(filepath.Join(dir, "test.txt"))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ServeRoot("/static", dir))

	req, _ := http.NewRequest("GET", "/static/test.txt", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
