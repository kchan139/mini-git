package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"minigit/internal/cli"
	"minigit/internal/repository"
)

func TestCommitWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	repoPath := filepath.Join(tempDir, "test_repo")

	os.Args = []string{"minigit", "init", repoPath}
	if err := cli.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	testFiles := map[string]string{
		"README.md": "# Test Repository",
		"main.go":   "package main\n\nfunc main() {}",
		"utils.go":  "package main\n\nfunc helper() {}",
	}

	for file, content := range testFiles {
		fullPath := filepath.Join(repoPath, file)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", file, err)
		}
	}

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(repoPath)

	for file := range testFiles {
		os.Args = []string{"minigit", "add", file}
		if err := cli.Execute(); err != nil {
			t.Fatalf("add %s failed: %v", file, err)
		}
	}

	os.Args = []string{"minigit", "commit", "-m", "Initial commit with multiple files"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	verifyCommitState(t, repoPath, testFiles)
}

func TestMultipleCommits(t *testing.T) {
	tempDir := t.TempDir()
	repoPath := filepath.Join(tempDir, "test_repo")

	os.Args = []string{"minigit", "init", repoPath}
	if err := cli.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(repoPath)

	file1 := "file1.txt"
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}

	os.Args = []string{"minigit", "add", file1}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add file1 failed: %v", err)
	}

	os.Args = []string{"minigit", "commit", "-m", "First commit"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("first commit failed: %v", err)
	}

	file2 := "file2.txt"
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	os.Args = []string{"minigit", "add", file2}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add file2 failed: %v", err)
	}

	os.Args = []string{"minigit", "commit", "-m", "Second commit"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("second commit failed: %v", err)
	}

	verifyMultipleCommitsState(t, repoPath)
}

func TestCommitClearsIndex(t *testing.T) {
	tempDir := t.TempDir()
	repoPath := filepath.Join(tempDir, "test_repo")

	os.Args = []string{"minigit", "init", repoPath}
	if err := cli.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(repoPath)

	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	os.Args = []string{"minigit", "add", testFile}
	if err := cli.Execute(); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	repo, err := repository.NewRepository(repoPath)
	if err != nil {
		t.Fatalf("failed to open repository: %v", err)
	}

	index, err := repo.GetIndex()
	if err != nil {
		t.Fatalf("failed to get index: %v", err)
	}

	entries := index.GetEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry before commit, got %d", len(entries))
	}

	os.Args = []string{"minigit", "commit", "-m", "Test commit"}
	if err := cli.Execute(); err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	repoCheck, err := repository.NewRepository(repoPath)
	if err != nil {
		t.Fatalf("failed to reopen repository: %v", err)
	}

	indexCheck, err := repoCheck.GetIndex()
	if err != nil {
		t.Fatalf("failed to get index after commit: %v", err)
	}

	entriesAfter := indexCheck.GetEntries()
	if len(entriesAfter) != 0 {
		t.Fatalf("expected 0 entries after commit, got %d", len(entriesAfter))
	}
}

func verifyCommitState(t *testing.T, repoPath string, _ map[string]string) {
	repo, err := repository.NewRepository(repoPath)
	if err != nil {
		t.Fatalf("failed to open repository: %v", err)
	}

	refsMan, err := repo.GetRefsManager()
	if err != nil {
		t.Fatalf("failed to get refs manager: %v", err)
	}

	headRef, err := refsMan.GetHead()
	if err != nil {
		t.Fatalf("failed to get HEAD: %v", err)
	}

	if !strings.HasPrefix(headRef, "refs/heads/") {
		t.Fatalf("HEAD should point to a branch, got: %s", headRef)
	}

	branchName := strings.TrimPrefix(headRef, "refs/heads/")
	commitHash, err := refsMan.GetBranch(branchName)
	if err != nil {
		t.Fatalf("failed to get branch commit: %v", err)
	}

	if commitHash == "" {
		t.Fatal("branch should point to a commit")
	}

	store, err := repo.GetObjectStore()
	if err != nil {
		t.Fatalf("failed to get object store: %v", err)
	}

	commitObj, err := store.LoadObject(commitHash)
	if err != nil {
		t.Fatalf("failed to load commit object: %v", err)
	}

	commit, err := store.ParseCommit(commitObj.Content)
	if err != nil {
		t.Fatalf("failed to parse commit: %v", err)
	}

	if commit.Tree == "" {
		t.Fatal("commit should have a tree hash")
	}

	if commit.Message == "" {
		t.Fatal("commit should have a message")
	}

	index, err := repo.GetIndex()
	if err != nil {
		t.Fatalf("failed to get index: %v", err)
	}

	entries := index.GetEntries()
	if len(entries) != 0 {
		t.Fatalf("index should be empty after commit, got %d entries", len(entries))
	}
}

func verifyMultipleCommitsState(t *testing.T, repoPath string) {
	repo, err := repository.NewRepository(repoPath)
	if err != nil {
		t.Fatalf("failed to open repository: %v", err)
	}

	refsMan, err := repo.GetRefsManager()
	if err != nil {
		t.Fatalf("failed to get refs manager: %v", err)
	}

	headRef, err := refsMan.GetHead()
	if err != nil {
		t.Fatalf("failed to get HEAD: %v", err)
	}

	branchName := strings.TrimPrefix(headRef, "refs/heads/")
	commitHash, err := refsMan.GetBranch(branchName)
	if err != nil {
		t.Fatalf("failed to get branch commit: %v", err)
	}

	store, err := repo.GetObjectStore()
	if err != nil {
		t.Fatalf("failed to get object store: %v", err)
	}

	commitObj, err := store.LoadObject(commitHash)
	if err != nil {
		t.Fatalf("failed to load latest commit: %v", err)
	}

	commit, err := store.ParseCommit(commitObj.Content)
	if err != nil {
		t.Fatalf("failed to parse latest commit: %v", err)
	}

	if len(commit.Parents) != 1 {
		t.Fatalf("second commit should have 1 parent, got %d", len(commit.Parents))
	}

	parentObj, err := store.LoadObject(commit.Parents[0])
	if err != nil {
		t.Fatalf("failed to load parent commit: %v", err)
	}

	parentCommit, err := store.ParseCommit(parentObj.Content)
	if err != nil {
		t.Fatalf("failed to parse parent commit: %v", err)
	}

	if len(parentCommit.Parents) != 0 {
		t.Fatalf("first commit should have 0 parents, got %d", len(parentCommit.Parents))
	}
}
