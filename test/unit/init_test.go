package unit

import (
	"os"
	"path/filepath"
	"testing"

	"minigit/test/fixtures"
)

func TestInitCurrentDirectory(t *testing.T) {
	cwd := fixtures.InitRepo(t) // inits inside a temp dir
	if _, err := os.Stat(filepath.Join(cwd, ".minigit")); os.IsNotExist(err) {
		t.Fatal(".minigit not created")
	}
}

func TestInitSpecificDirectory(t *testing.T) {
	// mimic CLI: init new_repo under temp
	repoPath := filepath.Join(t.TempDir(), "new_repo")
	fixtures.RunCLI(t, "init", repoPath)

	if _, err := os.Stat(filepath.Join(repoPath, ".minigit")); os.IsNotExist(err) {
		t.Fatal("repository not created at specified path")
	}
}

func TestInitDuplicateRepository(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// second init should succeed (reinitialize)
	err := fixtures.TryCLI(t, "init")
	if err != nil {
		t.Fatalf("expected reinit to succeed, got error: %v", err)
	}
}
