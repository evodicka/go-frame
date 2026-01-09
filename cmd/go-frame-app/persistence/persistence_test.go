package persistence_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/persistence"
)

func setupTestDB(t *testing.T) *persistence.Storage {
	// Create images dir required for InitBuckets -> prepopulateImages
	// prepopulate uses "images" directory relative to CWD.
	// We might need to ensure CWD is correct or mock it?
	// For simplicity, we create "images" in CWD.
	_ = os.MkdirAll("images", 0755)

	dbPath := filepath.Join(t.TempDir(), "test_persistence.db")
	storage, err := persistence.NewStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	// Cleanup
	t.Cleanup(func() {
		storage.Close()
		os.RemoveAll("images")
	})

	return storage
}

func TestImageCRUD(t *testing.T) {
	storage := setupTestDB(t)

	// Test Save
	img, err := storage.SaveImageMetadata("test_image.jpg")
	if err != nil {
		t.Fatalf("Failed to save image: %v", err)
	}
	if img.Path != "test_image.jpg" {
		t.Errorf("Expected path test_image.jpg, got %s", img.Path)
	}
	if img.Type != model.ImageType {
		t.Errorf("Expected type IMAGE, got %s", img.Type)
	}

	// Test Load
	loadedImg, err := storage.LoadImage(img.Id)
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}
	if loadedImg != img {
		t.Errorf("Loaded image mismatch: got %v, want %v", loadedImg, img)
	}

	// Test LoadAll
	images, err := storage.LoadImages()
	if err != nil {
		t.Fatalf("Failed to load images: %v", err)
	}
	if len(images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(images))
	}

	// Test Delete
	// Create dummy file for deletion
	os.Mkdir("images", 0755)
	f, _ := os.Create("images/test_image.jpg")
	f.Close()
	defer os.RemoveAll("images")

	err = storage.DeleteImage(img.Id)
	if err != nil {
		t.Fatalf("Failed to delete image: %v", err)
	}

	_, err = storage.LoadImage(img.Id)
	if err == nil {
		t.Error("Expected error loading deleted image, got nil")
	}
}

func TestConfigPersistence(t *testing.T) {
	storage := setupTestDB(t)

	// Initial config should use defaults if prepopulated, check prepopulate logic
	config, err := storage.GetConfiguration()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	// Default is 60
	if config.ImageDuration != 60 {
		t.Errorf("Expected default duration 60, got %d", config.ImageDuration)
	}

	// Update
	newConfig := model.Config{
		ImageDuration: 120,
		RandomOrder:   true,
	}
	err = storage.UpdateConfiguration(newConfig)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	updatedConfig, err := storage.GetConfiguration()
	if err != nil {
		t.Fatalf("Failed to get updated config: %v", err)
	}
	if updatedConfig != newConfig {
		t.Errorf("Config mismatch: got %v, want %v", updatedConfig, newConfig)
	}
}

func TestStatusPersistence(t *testing.T) {
	storage := setupTestDB(t)

	// Get initial status
	_, err := storage.GetCurrentStatus()
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	// Update status
	err = storage.UpdateImageStatus(10)
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	updatedStatus, err := storage.GetCurrentStatus()

	if err != nil {
		t.Fatalf("Failed to get updated status: %v", err)
	}
	if updatedStatus.CurrentImageId != 10 {
		t.Errorf("Expected current image ID 10, got %d", updatedStatus.CurrentImageId)
	}
	if time.Since(updatedStatus.LastSwitch) > time.Second {
		t.Error("LastSwitch time seems too old")
	}
}

func TestLoadNextImage(t *testing.T) {
	storage := setupTestDB(t)

	// Create 3 images
	img1, _ := storage.SaveImageMetadata("img1.jpg")
	img2, _ := storage.SaveImageMetadata("img2.jpg")
	img3, _ := storage.SaveImageMetadata("img3.jpg")

	// Order: img1, img2, img3 (default order is creation order based on sequences)

	// Next of img1 should be img2
	next, err := storage.LoadNextImage(img1.Id)
	if err != nil {
		t.Fatalf("LoadNextImage failed: %v", err)
	}
	if next.Id != img2.Id {
		t.Errorf("Expected next of img1 to be img2, got img%d", next.Id)
	}

	// Next of img3 (last) should be img1 (loop)
	next, err = storage.LoadNextImage(img3.Id)
	if err != nil {
		t.Fatalf("LoadNextImage failed: %v", err)
	}
	if next.Id != img1.Id {
		t.Errorf("Expected loop back to img1, got img%d", next.Id)
	}

	// Test Reorder
	err = storage.ReorderImages([]model.Image{img3, img2, img1})
	if err != nil {
		t.Fatalf("Reorder failed: %v", err)
	}

	// Now order: img3, img2, img1
	next, err = storage.LoadNextImage(img3.Id)
	if err != nil {
		t.Fatalf("LoadNextImage failed: %v", err)
	}
	if next.Id != img2.Id {
		t.Errorf("After reorder, expected next of img3 to be img2, got img%d", next.Id)
	}
}
