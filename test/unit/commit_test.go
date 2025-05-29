package unit

import (
	"strings"
	"testing"

	"minigit/test/fixtures"
)

func TestCommitWithMessage(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// prepare
	fixtures.CreateFiles(t, repoPath, map[string]string{"test.txt": "test content"})
	fixtures.RunCLI(t, "add", "test.txt")

	// commit
	fixtures.RunCLI(t, "commit", "-m", "Initial commit")
}

func TestCommitEmptyMessage(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// prepare
	fixtures.CreateFiles(t, repoPath, map[string]string{"test.txt": "content"})
	fixtures.RunCLI(t, "add", "test.txt")

	// commit with empty -m
	err := fixtures.TryCLI(t, "commit", "-m", "")
	if err == nil {
		t.Fatal("expected error for empty message")
	}
	if !strings.Contains(err.Error(), "switch `m' requires a value") {
		t.Fatalf("wrong error: %v", err)
	}
}

func TestCommitNothingToCommit(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// commit straight away
	err := fixtures.TryCLI(t, "commit", "-m", "Empty commit")
	if err == nil || !strings.Contains(err.Error(), "no changes added to commit") {
		t.Fatalf("expected no changes error, got %v", err)
	}
}

func TestCommitNoChanges(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	fixtures.CreateFiles(t, repoPath, map[string]string{"test.txt": "content"})
	fixtures.RunCLI(t, "add", "test.txt")
	fixtures.RunCLI(t, "commit", "-m", "First commit")

	// second commit without changes
	err := fixtures.TryCLI(t, "commit", "-m", "Second commit")
	if err == nil || !strings.Contains(err.Error(), "no changes added to commit") {
		t.Fatalf("expected no changes error, got %v", err)
	}
}

func TestCommitMissingMessage(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	fixtures.CreateFiles(t, repoPath, map[string]string{"test.txt": "content"})
	fixtures.RunCLI(t, "add", "test.txt")

	// omit -m
	err := fixtures.TryCLI(t, "commit")
	if err == nil {
		t.Fatal("expected missing -m error")
	}
}
