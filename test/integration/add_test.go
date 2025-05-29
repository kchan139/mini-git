// test/integration/add_test.go
package integration

import (
	"testing"

	"minigit/internal/repository"
	"minigit/test/fixtures"
)

func TestAddWorkflow(t *testing.T) {
	// 1) init repo
	repoPath := fixtures.InitRepo(t)

	// 2) create files under repo
	testFiles := map[string]string{
		"README.md":    "# Test Repository",
		"src/main.go":  "package main\n\nfunc main() {}",
		"src/utils.go": "package main\n\nfunc helper() {}",
		"config.json":  `{"version": "1.0"}`,
	}
	fixtures.CreateFiles(t, repoPath, testFiles)

	// 3) switch into it
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// 4) run add commands
	fixtures.RunCLI(t, "add", "README.md")
	fixtures.RunCLI(t, "add", "src")

	// 5) verify index
	repo, err := repository.NewRepository(repoPath)
	if err != nil {
		t.Fatalf("open repo: %v", err)
	}
	index, err := repo.GetIndex()
	if err != nil {
		t.Fatalf("get index: %v", err)
	}
	entries := index.GetEntries()

	// expect 3 entries
	expected := []string{"README.md", "src/main.go", "src/utils.go"}
	if len(entries) != len(expected) {
		t.Fatalf("got %d entries, want %d", len(entries), len(expected))
	}
	for _, want := range expected {
		if _, ok := entries[want]; !ok {
			t.Fatalf("%s not in index", want)
		}
	}
}
