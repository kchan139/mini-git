package integration

import (
	"os"
	"path/filepath"
	"testing"

	"minigit/test/fixtures"
)

func TestCompleteInitWorkflow(t *testing.T) {
	repoPath := fixtures.InitRepo(t)
	checkCompleteStructure(t, repoPath)
}

func checkCompleteStructure(t *testing.T, repoPath string) {
	paths := map[string]bool{
		filepath.Join(repoPath, ".minigit"):                  true,
		filepath.Join(repoPath, ".minigit", "objects"):       true,
		filepath.Join(repoPath, ".minigit", "refs"):          true,
		filepath.Join(repoPath, ".minigit", "refs", "heads"): true,
		filepath.Join(repoPath, ".minigit", "HEAD"):          false,
	}

	for path, isDir := range paths {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("missing: %s", path)
		}
		if info.IsDir() != isDir {
			t.Fatalf("wrong type for %s: isDir=%v", path, isDir)
		}
	}
}
