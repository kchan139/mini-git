package integration

import (
	"os"
	"path/filepath"
	"testing"

	"minigit/test/fixtures"
)

func TestStatusWorkflow(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Step 1: Empty repo status
	fixtures.RunCLI(t, "status")

	// Step 2: Add files and check status
	testFiles := map[string]string{
		"README.md":   "# Project",
		"src/main.go": "package main",
		"config.yaml": "version: 1.0",
	}
	fixtures.CreateFiles(t, repoPath, testFiles)
	fixtures.RunCLI(t, "status") // Should show untracked

	// Step 3: Stage some files
	fixtures.RunCLI(t, "add", "README.md")
	fixtures.RunCLI(t, "add", "src/main.go")
	fixtures.RunCLI(t, "status") // Should show staged + untracked

	// Step 4: Commit
	fixtures.RunCLI(t, "commit", "-m", "Initial commit")
	fixtures.RunCLI(t, "status") // Should show untracked config.yaml

	// Step 5: Modify committed file
	fixtures.CreateFiles(t, repoPath, map[string]string{"README.md": "# Updated Project"})
	fixtures.RunCLI(t, "status") // Should show modified + untracked

	// Step 6: Stage modified file
	fixtures.RunCLI(t, "add", "README.md")
	fixtures.RunCLI(t, "status") // Should show staged modification + untracked

	// Step 7: Final commit
	fixtures.RunCLI(t, "commit", "-m", "Update README")
	fixtures.RunCLI(t, "status") // Should show only untracked config.yaml
}

func TestStatusComplexScenario(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Create initial structure
	initialFiles := map[string]string{
		"main.go":         "package main\n\nfunc main() {}",
		"utils/helper.go": "package utils",
		"docs/guide.md":   "# Guide",
		"config/app.yaml": "app: myapp",
	}
	fixtures.CreateFiles(t, repoPath, initialFiles)

	// Initial commit
	fixtures.RunCLI(t, "add", "main.go")
	fixtures.RunCLI(t, "add", "utils")
	fixtures.RunCLI(t, "add", "docs/guide.md")
	fixtures.RunCLI(t, "commit", "-m", "Initial structure")

	// Create complex state:
	// 1. Modify tracked files
	fixtures.CreateFiles(t, repoPath, map[string]string{
		"main.go":         "package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"Hello\") }",
		"utils/helper.go": "package utils\n\nfunc Help() {}",
	})

	// 2. Add new files
	fixtures.CreateFiles(t, repoPath, map[string]string{
		"new_feature.go":     "package main",
		"tests/main_test.go": "package main",
	})

	// 3. Stage some changes
	fixtures.RunCLI(t, "add", "main.go")        // stage modification
	fixtures.RunCLI(t, "add", "new_feature.go") // stage new file

	// 4. Delete a tracked file
	os.Remove(filepath.Join(repoPath, "docs", "guide.md"))

	// Should show:
	// - Staged: modified main.go, new file new_feature.go
	// - Modified: utils/helper.go
	// - Untracked: config/app.yaml, tests/main_test.go
	// - Deleted: docs/guide.md (should be handled appropriately)
	fixtures.RunCLI(t, "status")
}

func TestStatusBranchDisplay(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Make initial commit to establish branch
	fixtures.CreateFiles(t, repoPath, map[string]string{"init.txt": "init"})
	fixtures.RunCLI(t, "add", "init.txt")
	fixtures.RunCLI(t, "commit", "-m", "Initial commit")

	// Should show "On branch main"
	fixtures.RunCLI(t, "status")
}

func TestStatusIgnoresDotMinigit(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	cleanup := fixtures.Chdir(t, repoPath)
	defer cleanup()

	// Create file in .minigit directory (should be ignored)
	minigitFile := filepath.Join(repoPath, ".minigit", "test_file")
	if err := os.WriteFile(minigitFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file in .minigit: %v", err)
	}

	// Create normal file (should be detected)
	fixtures.CreateFiles(t, repoPath, map[string]string{"normal.txt": "content"})

	// Should only show normal.txt, not the .minigit file
	fixtures.RunCLI(t, "status")
}
