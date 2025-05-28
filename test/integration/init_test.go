// File: test/integration/full_init_test.go
package integration

import (
	"os"
	"path/filepath"
	"testing"

	"minigit/internal/cli"
)

func TestCompleteInitWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	repoPath := filepath.Join(tempDir, "integration_repo")

	os.Args = []string{"minigit", "init", repoPath}
	err := cli.Execute()

	if err != nil {
		t.Fatalf("integration test failed: %v", err)
	}

	// Verify complete structure
	checkCompleteStructure(t, repoPath)
}

func checkCompleteStructure(t *testing.T, repoPath string) {
	// All required paths with their types
	paths := map[string]bool{
		filepath.Join(repoPath, ".minigit"):                  true,  // directory
		filepath.Join(repoPath, ".minigit", "objects"):       true,  // directory
		filepath.Join(repoPath, ".minigit", "refs"):          true,  // directory
		filepath.Join(repoPath, ".minigit", "refs", "heads"): true,  // directory
		filepath.Join(repoPath, ".minigit", "HEAD"):          false, // file
	}

	for path, shouldBeDir := range paths {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("missing path: %s", path)
		}
		if info.IsDir() != shouldBeDir {
			t.Fatalf("wrong type for %s: isDir=%v, expected=%v",
				path, info.IsDir(), shouldBeDir)
		}
	}
}
