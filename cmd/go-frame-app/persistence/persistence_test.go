package persistence_test

import (
	"os"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/model"
	"go.evodicka.dev/go-frame/cmd/go-frame-app/persistence"
)

const tempDbName = "test_persistence.db"

func TestMain(m *testing.M) {
	// Setup
	os.Mkdir("images", 0755)

	db, err := bolt.Open(tempDbName, 0600, nil)
	if err != nil {
		panic(err)
	}
	persistence.Db = db

	// Run tests
	code := m.Run()

	// Teardown
	db.Close()
	os.RemoveAll("images")
	os.Remove(tempDbName)
	os.Exit(code)
}

func setupTestDB() error {
	// Create images dir required for InitBuckets -> prepopulateImages
	_ = os.MkdirAll("images", 0755)

	// Clear existing buckets to ensure clean state
	err := persistence.Db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("images"))
		tx.DeleteBucket([]byte("order"))
		tx.DeleteBucket([]byte("configuration"))
		tx.DeleteBucket([]byte("status"))
		return nil
	})
	if err != nil {
		// Ignore error if buckets don't exist
	}
	// Re-init buckets
	persistence.InitBuckets()
	return nil
}

func TestImageCRUD(t *testing.T) {
	setupTestDB()

	// Test Save
	img, err := persistence.SaveImageMetadata("test_image.jpg")
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
	loadedImg, err := persistence.LoadImage(img.Id)
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}
	if loadedImg != img {
		t.Errorf("Loaded image mismatch: got %v, want %v", loadedImg, img)
	}

	// Test LoadAll
	images, err := persistence.LoadImages()
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

	err = persistence.DeleteImage(img.Id)
	if err != nil {
		t.Fatalf("Failed to delete image: %v", err)
	}

	_, err = persistence.LoadImage(img.Id)
	if err == nil {
		t.Error("Expected error loading deleted image, got nil")
	}
}

func TestConfigPersistence(t *testing.T) {
	setupTestDB()

	// Initial config should use defaults if prepopulated, check prepopulate logic
	config, err := persistence.GetConfiguration()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	// Default is 60
	if config.ImageDuration != 60 {
		t.Errorf("Expected default duration 60, got %d", config.ImageDuration)
	}

	// Update
	newConfig := persistence.Config{
		ImageDuration: 120,
		RandomOrder:   true,
	}
	err = persistence.UpdateConfiguration(newConfig)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	updatedConfig, err := persistence.GetConfiguration()
	if err != nil {
		t.Fatalf("Failed to get updated config: %v", err)
	}
	if updatedConfig != newConfig {
		t.Errorf("Config mismatch: got %v, want %v", updatedConfig, newConfig)
	}
}

func TestStatusPersistence(t *testing.T) {
	setupTestDB()

	// Get initial status
	_, err := persistence.GetCurrentStatus()
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	// Update status
	err = persistence.UpdateImageStatus(10)
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	updatedStatus, err := persistence.GetCurrentStatus()

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
	setupTestDB()

	// Create 3 images
	img1, _ := persistence.SaveImageMetadata("img1.jpg")
	img2, _ := persistence.SaveImageMetadata("img2.jpg")
	img3, _ := persistence.SaveImageMetadata("img3.jpg")

	// Order: img1, img2, img3 (default order is creation order based on sequences)

	// Next of img1 should be img2
	next, err := persistence.LoadNextImage(img1.Id)
	if err != nil {
		t.Fatalf("LoadNextImage failed: %v", err)
	}
	if next.Id != img2.Id {
		t.Errorf("Expected next of img1 to be img2, got img%d", next.Id)
	}

	// Next of img3 (last) should be img1 (loop)
	next, err = persistence.LoadNextImage(img3.Id)
	if err != nil {
		t.Fatalf("LoadNextImage failed: %v", err)
	}
	if next.Id != img1.Id {
		t.Errorf("Expected loop back to img1, got img%d", next.Id)
	}

	// Test Reorder
	err = persistence.ReorderImages([]persistence.Image{img3, img2, img1})
	if err != nil {
		t.Fatalf("Reorder failed: %v", err)
	}

	// Now order: img3, img2, img1
	next, err = persistence.LoadNextImage(img3.Id)
	if err != nil {
		t.Fatalf("LoadNextImage failed: %v", err)
	}
	if next.Id != img2.Id {
		t.Errorf("After reorder, expected next of img3 to be img2, got img%d", next.Id)
	}
}
