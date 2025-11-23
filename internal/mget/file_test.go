package mget

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetOutputDirectory(t *testing.T) {
	dir, err := GetOutputDirectory()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if dir == "" {
		t.Error("Expected non-empty directory path")
	}

	// Should end with mget-downloads
	expectedSuffix := "mget-downloads"
	if !filepath.IsAbs(dir) {
		t.Error("Expected absolute path")
	}

	if filepath.Base(dir) != expectedSuffix {
		t.Errorf("Expected directory to end with %q, got %q", expectedSuffix, filepath.Base(dir))
	}
}

func TestGenerateFilename(t *testing.T) {
	// Generate multiple filenames to ensure uniqueness
	filenames := make(map[string]bool)

	for i := 0; i < 100; i++ {
		filename := GenerateFilename()

		if filename == "" {
			t.Error("Expected non-empty filename")
		}

		// Check for uniqueness
		if filenames[filename] {
			t.Errorf("Generated duplicate filename: %s", filename)
		}
		filenames[filename] = true

		// Filename should be a valid UUID format (36 characters with hyphens)
		if len(filename) != 36 {
			t.Errorf("Expected UUID format (36 chars), got %d chars: %s", len(filename), filename)
		}
	}
}

func TestGetOutputDirectory_Integration(t *testing.T) {
	// Test that the directory can be created
	dir, err := GetOutputDirectory()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Create a temporary directory to test
	testDir := filepath.Join(dir, "test")
	defer os.RemoveAll(testDir)

	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Verify directory exists
	info, err := os.Stat(testDir)
	if err != nil {
		t.Fatalf("Directory should exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("Expected directory, got file")
	}
}
