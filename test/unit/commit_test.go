package unit

import (
	"os"
	"strings"
	"testing"

	"minigit/internal/cli"
)

func TestCommitWithMessage(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	os.Args = []string{"minigit", "init"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	testFile := "test.txt"
	testContent := "test content"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	os.Args = []string{"minigit", "add", testFile}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	os.Args = []string{"minigit", "commit", "-m", "Initial commit"}
	err := cli.Execute()

	if err != nil {
		t.Fatalf("commit failed: %v", err)
	}
}

func TestCommitEmptyMessage(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	os.Args = []string{"minigit", "init"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	os.Args = []string{"minigit", "add", testFile}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	os.Args = []string{"minigit", "commit", "-m", ""}
	err := cli.Execute()

	if err == nil {
		t.Fatal("expected error for empty commit message")
	}

	if !strings.Contains(err.Error(), "switch `m' requires a value") {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestCommitNothingToCommit(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	os.Args = []string{"minigit", "init"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	os.Args = []string{"minigit", "commit", "-m", "Empty commit"}
	err := cli.Execute()

	if !strings.Contains(err.Error(), "no changes added to commit") {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestCommitNoChanges(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	os.Args = []string{"minigit", "init"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	os.Args = []string{"minigit", "add", testFile}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	os.Args = []string{"minigit", "commit", "-m", "First commit"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("first commit failed: %v", err)
	}

	os.Args = []string{"minigit", "commit", "-m", "Second commit"}
	err := cli.Execute()

	if !strings.Contains(err.Error(), "no changes added to commit") {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestCommitMissingMessage(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	os.Args = []string{"minigit", "init"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	os.Args = []string{"minigit", "add", testFile}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	os.Args = []string{"minigit", "commit"}
	err := cli.Execute()

	if err == nil {
		t.Fatal("expected error for missing -m flag")
	}
}
