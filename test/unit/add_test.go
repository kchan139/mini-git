package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"minigit/internal/repository"
	"minigit/test/fixtures"
)

func TestAddSingleFile(t *testing.T) {
	// init repo
	repoPath := fixtures.InitRepo(t)

	// create file
	rel := "test.txt"
	content := "Hello, MiniGit!"
	fixtures.CreateFiles(t, repoPath, map[string]string{rel: content})

	// cd
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// add
	fixtures.RunCLI(t, "add", rel)

	// verify index
	repo, _ := repository.NewRepository(repoPath)
	index, _ := repo.GetIndex()
	entries := index.GetEntries()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	entry, ok := entries[rel]
	if !ok {
		t.Fatalf("%s not in index", rel)
	}
	if entry.Size != int64(len(content)) {
		t.Fatalf("size mismatch: got %d", entry.Size)
	}
}

func TestAddDirectory(t *testing.T) {
	repoPath := fixtures.InitRepo(t)

	// create subdir files
	files := map[string]string{
		"subdir/file1.txt": "Content 1",
		"subdir/file2.txt": "Content 2",
	}
	fixtures.CreateFiles(t, repoPath, files)

	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// add dir
	fixtures.RunCLI(t, "add", "subdir")

	repo, _ := repository.NewRepository(repoPath)
	index, _ := repo.GetIndex()
	entries := index.GetEntries()

	if len(entries) != len(files) {
		t.Fatalf("expected %d entries, got %d", len(files), len(entries))
	}
	for path := range files {
		rel := filepath.ToSlash(path)
		if _, ok := entries[rel]; !ok {
			t.Fatalf("%s not found", rel)
		}
	}
}

func TestAddNonExistentFile(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// try to add missing
	os.Args = []string{"minigit", "add", "nope.txt"}
	err := fixtures.TryCLI(t, "add", "nope.txt")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
	if !strings.Contains(err.Error(), "did not match any files") {
		t.Fatalf("wrong error: %v", err)
	}
}
