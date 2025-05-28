package integration

import (
	"os"
	"path/filepath"
	"testing"

	"minigit/internal/cli"
	"minigit/internal/repository"
)

func TestAddWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	repoPath := filepath.Join(tempDir, "test_repo")

	// Initialize repository
	os.Args = []string{"minigit", "init", repoPath}
	if err := cli.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Create test files
	testFiles := map[string]string{
		"README.md":    "# Test Repository",
		"src/main.go":  "package main\n\nfunc main() {}",
		"src/utils.go": "package main\n\nfunc helper() {}",
		"config.json":  `{"version": "1.0"}`,
	}

	for file, content := range testFiles {
		fullPath := filepath.Join(repoPath, file)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", file, err)
		}
	}

	// Change to repo directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(repoPath)

	// Add individual file
	os.Args = []string{"minigit", "add", "README.md"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add README.md failed: %v", err)
	}

	// Add directory
	os.Args = []string{"minigit", "add", "src"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add src/ failed: %v", err)
	}

	// Verify repository state
	repo, err := repository.NewRepository(repoPath)
	if err != nil {
		t.Fatalf("failed to open repository: %v", err)
	}

	index, err := repo.GetIndex()
	if err != nil {
		t.Fatalf("failed to get index: %v", err)
	}

	entries := index.GetEntries()
	expectedFiles := []string{"README.md", "src/main.go", "src/utils.go"}

	if len(entries) != len(expectedFiles) {
		t.Fatalf("expected %d entries, got %d", len(expectedFiles), len(entries))
	}

	for _, file := range expectedFiles {
		if _, exists := entries[file]; !exists {
			t.Fatalf("file %s not in index", file)
		}
	}

	// Verify objects were created
	store, err := repo.GetObjectStore()
	if err != nil {
		t.Fatalf("failed to get object store: %v", err)
	}

	for file, entry := range entries {
		obj, err := store.LoadObject(entry.Hash)
		if err != nil {
			t.Fatalf("failed to load object for %s: %v", file, err)
		}

		expectedContent := testFiles[file]
		if string(obj.Content) != expectedContent {
			t.Fatalf("wrong content for %s: got %q, want %q",
				file, string(obj.Content), expectedContent)
		}
	}
}
