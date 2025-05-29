package unit

import (
	"os"
	"path/filepath"
	"strings"
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

	// second init should error
	err := fixtures.TryCLI(t, "init")
	if err == nil || !strings.Contains(err.Error(), "already a minigit repository") {
		t.Fatalf("expected duplicate-repo error, got %v", err)
	}
}
