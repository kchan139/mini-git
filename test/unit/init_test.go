// Organized test structure for test/ directory

// File: test/unit/cli_init_test.go
package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"minigit/internal/cli"
)

func TestInitCurrentDirectory(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	os.Args = []string{"minigit", "init"}
	err := cli.Execute()

	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	if _, err := os.Stat(".minigit"); os.IsNotExist(err) {
		t.Fatal(".minigit directory not created")
	}
}

func TestInitSpecificDirectory(t *testing.T) {
	tempDir := t.TempDir()
	repoPath := filepath.Join(tempDir, "new_repo")

	os.Args = []string{"minigit", "init", repoPath}
	err := cli.Execute()

	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	minigitPath := filepath.Join(repoPath, ".minigit")
	if _, err := os.Stat(minigitPath); os.IsNotExist(err) {
		t.Fatal("repository not created at specified path")
	}
}

func TestInitDuplicateRepository(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	// First init
	os.Args = []string{"minigit", "init"}
	err := cli.Execute()
	if err != nil {
		t.Fatalf("first init failed: %v", err)
	}

	// Duplicate init should fail
	err = cli.Execute()
	if err == nil {
		t.Fatal("duplicate init should fail")
	}
	if !strings.Contains(err.Error(), "already a minigit repository") {
		t.Fatalf("wrong error message: %v", err)
	}
}
