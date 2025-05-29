package unit

import (
	"os"
	"path/filepath"
	"testing"

	"minigit/internal/repository"
)

func TestRepositoryInitialization(t *testing.T) {
	tempDir := t.TempDir()

	repo, err := repository.NewRepository(tempDir)
	if err != nil {
		t.Fatalf("NewRepository failed: %v", err)
	}

	err = repo.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Check directory structure
	dirs := []string{
		".minigit",
		".minigit/objects",
		".minigit/refs",
		".minigit/refs/heads",
	}

	for _, dir := range dirs {
		path := filepath.Join(tempDir, dir)
		if info, err := os.Stat(path); os.IsNotExist(err) {
			t.Fatalf("directory missing: %s", dir)
		} else if !info.IsDir() {
			t.Fatalf("not a directory: %s", dir)
		}
	}

	// Check HEAD file
	headPath := filepath.Join(tempDir, ".minigit", "HEAD")
	content, err := os.ReadFile(headPath)
	if err != nil {
		t.Fatalf("HEAD file missing: %v", err)
	}

	expected := "ref: refs/heads/main\n"
	if string(content) != expected {
		t.Fatalf("wrong HEAD content: got %q, want %q", string(content), expected)
	}
}
