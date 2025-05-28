// File: test/unit/add_test.go
package unit

import (
	"os"
	"path/filepath"
	"testing"

	"minigit/internal/cli"
	"minigit/internal/repository"
)

func TestAddSingleFile(t *testing.T) {
	tempDir := t.TempDir()

	// Initialize repository
	repo, err := repository.NewRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	if err := repo.Initialize(); err != nil {
		t.Fatalf("failed to initialize repository: %v", err)
	}

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, MiniGit!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Change to repository directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	// Run add command
	os.Args = []string{"minigit", "add", "test.txt"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	// Create a fresh repository instance to check the persisted index
	repoCheck, err := repository.NewRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to create check repository: %v", err)
	}

	index, err := repoCheck.GetIndex()
	if err != nil {
		t.Fatalf("failed to get index: %v", err)
	}

	entries := index.GetEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry in index, got %d", len(entries))
	}

	entry, exists := entries["test.txt"]
	if !exists {
		t.Fatal("test.txt not found in index")
	}

	if entry.Size != int64(len(testContent)) {
		t.Fatalf("wrong file size in index: got %d, want %d", entry.Size, len(testContent))
	}
}

func TestAddDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Initialize repository
	repo, err := repository.NewRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	if err := repo.Initialize(); err != nil {
		t.Fatalf("failed to initialize repository: %v", err)
	}

	// Create test files in subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	os.MkdirAll(subDir, 0755)

	files := map[string]string{
		"file1.txt": "Content 1",
		"file2.txt": "Content 2",
	}

	for filename, content := range files {
		path := filepath.Join(subDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", filename, err)
		}
	}

	// Change to repository directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	// Run add command on directory
	os.Args = []string{"minigit", "add", "subdir"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	// Create a fresh repository instance to check the persisted index
	repoCheck, err := repository.NewRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to create check repository: %v", err)
	}

	index, err := repoCheck.GetIndex()
	if err != nil {
		t.Fatalf("failed to get index: %v", err)
	}

	entries := index.GetEntries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries in index, got %d", len(entries))
	}

	for filename := range files {
		expectedPath := filepath.Join("subdir", filename)
		if _, exists := entries[expectedPath]; !exists {
			t.Fatalf("file %s not found in index", expectedPath)
		}
	}
}

func TestAddNonExistentFile(t *testing.T) {
	tempDir := t.TempDir()

	// Initialize repository
	repo, err := repository.NewRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	if err := repo.Initialize(); err != nil {
		t.Fatalf("failed to initialize repository: %v", err)
	}

	// Change to repository directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	// Try to add non-existent file
	os.Args = []string{"minigit", "add", "nonexistent.txt"}
	err = cli.Execute()

	if err == nil {
		t.Fatal("expected error when adding non-existent file")
	}

	if !contains(err.Error(), "did not match any files") {
		t.Fatalf("wrong error message: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
