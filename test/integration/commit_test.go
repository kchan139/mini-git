package integration

import (
	"strings"
	"testing"

	"minigit/internal/repository"
	"minigit/test/fixtures"
)

func TestCommitWorkflow(t *testing.T) {
	// init
	repoPath := fixtures.InitRepo(t)

	// create files
	testFiles := map[string]string{
		"README.md": "# Test Repository",
		"main.go":   "package main\n\nfunc main() {}",
		"utils.go":  "package main\n\nfunc helper() {}",
	}
	fixtures.CreateFiles(t, repoPath, testFiles)

	// cd
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// add all
	for file := range testFiles {
		fixtures.RunCLI(t, "add", file)
	}

	// commit
	fixtures.RunCLI(t, "commit", "-m", "Initial commit with multiple files")

	verifyCommitState(t, repoPath, testFiles)
}

func TestMultipleCommits(t *testing.T) {
	// init
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// first file & commit
	fixtures.CreateFiles(t, repoPath, map[string]string{"file1.txt": "content1"})
	fixtures.RunCLI(t, "add", "file1.txt")
	fixtures.RunCLI(t, "commit", "-m", "First commit")

	// second file & commit
	fixtures.CreateFiles(t, repoPath, map[string]string{"file2.txt": "content2"})
	fixtures.RunCLI(t, "add", "file2.txt")
	fixtures.RunCLI(t, "commit", "-m", "Second commit")

	verifyMultipleCommitsState(t, repoPath)
}

func TestCommitClearsIndex(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// add & pre-check
	fixtures.CreateFiles(t, repoPath, map[string]string{"test.txt": "content"})
	fixtures.RunCLI(t, "add", "test.txt")

	repo, _ := repository.NewRepository(repoPath)
	index, _ := repo.GetIndex()
	if len(index.GetEntries()) != 1 {
		t.Fatalf("expected 1 entry pre-commit")
	}

	// commit
	fixtures.RunCLI(t, "commit", "-m", "Test commit")

	// post-check
	repo2, _ := repository.NewRepository(repoPath)
	index2, _ := repo2.GetIndex()
	if len(index2.GetEntries()) != 0 {
		t.Fatalf("expected 0 entries post-commit")
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
