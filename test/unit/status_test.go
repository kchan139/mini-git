package unit

import (
	"strings"
	"testing"

	"minigit/test/fixtures"
)

func TestStatusCleanWorkingTree(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Create, add, and commit a file
	fixtures.CreateFiles(t, repoPath, map[string]string{"test.txt": "content"})
	fixtures.RunCLI(t, "add", "test.txt")
	fixtures.RunCLI(t, "commit", "-m", "Initial commit")

	// Should show clean working tree
	err := fixtures.TryCLI(t, "status")
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}
}

func TestStatusStagedFiles(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Create and add files
	testFiles := map[string]string{
		"new.txt":     "new file content",
		"another.txt": "another file",
	}
	fixtures.CreateFiles(t, repoPath, testFiles)
	fixtures.RunCLI(t, "add", "new.txt")
	fixtures.RunCLI(t, "add", "another.txt")

	// Should show staged files
	err := fixtures.TryCLI(t, "status")
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}
}

func TestStatusModifiedFiles(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Create, add, and commit
	fixtures.CreateFiles(t, repoPath, map[string]string{"test.txt": "original content"})
	fixtures.RunCLI(t, "add", "test.txt")
	fixtures.RunCLI(t, "commit", "-m", "Initial commit")

	// Modify the file
	fixtures.CreateFiles(t, repoPath, map[string]string{"test.txt": "modified content"})

	// Should show modified file
	err := fixtures.TryCLI(t, "status")
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}
}

func TestStatusUntrackedFiles(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Create files without adding them
	testFiles := map[string]string{
		"untracked1.txt": "untracked content 1",
		"untracked2.txt": "untracked content 2",
		"dir/file.txt":   "nested untracked file",
	}
	fixtures.CreateFiles(t, repoPath, testFiles)

	// Should show untracked files
	err := fixtures.TryCLI(t, "status")
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}
}

func TestStatusMixedState(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Initial commit with one file
	fixtures.CreateFiles(t, repoPath, map[string]string{"committed.txt": "committed content"})
	fixtures.RunCLI(t, "add", "committed.txt")
	fixtures.RunCLI(t, "commit", "-m", "Initial commit")

	// Create mixed state:
	// 1. Modify existing file
	fixtures.CreateFiles(t, repoPath, map[string]string{"committed.txt": "modified content"})

	// 2. Add a new file to staging
	fixtures.CreateFiles(t, repoPath, map[string]string{"staged.txt": "staged content"})
	fixtures.RunCLI(t, "add", "staged.txt")

	// 3. Create untracked file
	fixtures.CreateFiles(t, repoPath, map[string]string{"untracked.txt": "untracked content"})

	// Should show all categories
	err := fixtures.TryCLI(t, "status")
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}
}

func TestStatusFirstCommitState(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Add files without committing (first commit scenario)
	testFiles := map[string]string{
		"file1.txt": "content 1",
		"file2.txt": "content 2",
	}
	fixtures.CreateFiles(t, repoPath, testFiles)
	fixtures.RunCLI(t, "add", "file1.txt")
	fixtures.RunCLI(t, "add", "file2.txt")

	// Should show staged files for first commit
	err := fixtures.TryCLI(t, "status")
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}
}

func TestStatusEmptyRepository(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Should work on empty repo
	err := fixtures.TryCLI(t, "status")
	if err != nil {
		t.Fatalf("status command failed on empty repo: %v", err)
	}
}

func TestStatusOutsideRepository(t *testing.T) {
	tempDir := t.TempDir()
	cleanup := fixtures.Chdir(t, tempDir)
	defer cleanup()

	// Should fail gracefully on empty repo
	err := fixtures.TryCLI(t, "status")
	if err == nil {
		t.Fatal("expected error when running status outside repository")
	}
	if !strings.Contains(err.Error(), "not a minigit repository") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestStatusWithDirectories(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Create nested directory structure
	testFiles := map[string]string{
		"src/main.go":         "package main",
		"src/utils/helper.go": "package utils",
		"docs/README.md":      "# Documentation",
		"config.json":         "{}",
	}
	fixtures.CreateFiles(t, repoPath, testFiles)

	// Add some files
	fixtures.RunCLI(t, "add", "src/main.go")
	fixtures.RunCLI(t, "add", "config.json")

	// Should handle nested files correctly
	err := fixtures.TryCLI(t, "status")
	if err != nil {
		t.Fatalf("status command failed with directories: %v", err)
	}
}
