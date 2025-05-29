package fixtures

import (
	"os"
	"path/filepath"
	"testing"

	"minigit/internal/cli"
)

// InitRepo creates a temp directory, initializes a minigit repo there,
// and returns the repo path.
func InitRepo(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	repoPath := filepath.Join(tempDir, "test_repo")

	// init via CLI
	os.Args = []string{"minigit", "init", repoPath}
	if err := cli.Execute(); err != nil {
		t.Fatalf("fixtures.InitRepo: init failed: %v", err)
	}

	return repoPath
}

// CreateFiles makes all dirs and writes files under basePath.
func CreateFiles(t *testing.T, basePath string, files map[string]string) {
	t.Helper()

	for relPath, content := range files {
		full := filepath.Join(basePath, relPath)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("fixtures.CreateFiles: mkdir %s: %v", full, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("fixtures.CreateFiles: write %s: %v", full, err)
		}
	}
}

// Chdir switches cwd to dir for the duration of the test.
func Chdir(t *testing.T, dir string) func() {
	t.Helper()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("fixtures.Chdir: cwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("fixtures.Chdir: chdir to %s: %v", dir, err)
	}
	return func() {
		_ = os.Chdir(orig) // ignore error on cleanup
	}
}

// RunCLI runs `minigit` with the given args
func RunCLI(t *testing.T, args ...string) {
	t.Helper()

	cmd := append([]string{"minigit"}, args...)
	os.Args = cmd
	if err := cli.Execute(); err != nil {
		t.Fatalf("fixtures.RunCLI %v: %v", args, err)
	}
}

// TryCLI runs `minigit` with the given args and returns the error.
// Use this when we want to assert on the error instead of failing immediately.
func TryCLI(t *testing.T, args ...string) error {
	t.Helper()

	cmd := append([]string{"minigit"}, args...)
	os.Args = cmd
	return cli.Execute()
}
