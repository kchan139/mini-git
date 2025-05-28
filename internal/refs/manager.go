// Reference management (branches, HEAD, ...)
package refs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Manager struct {
	refsDir    string
	minigitDir string
}

func NewManager(minigitDir string) (*Manager, error) {
	return &Manager{
		refsDir:    filepath.Join(minigitDir, "refs"),
		minigitDir: minigitDir,
	}, nil
}

// Updates HEAD
func (m *Manager) SetHead(ref string) error {
	headPath := filepath.Join(m.minigitDir, "HEAD")
	content := fmt.Sprintf("ref: %s\n", ref)
	// 0644 ~ owners can read and write, others can only read
	return os.WriteFile(headPath, []byte(content), 0644)
}

// Returns current HEAD reference
func (m *Manager) GetHead() (string, error) {
	headPath := filepath.Join(m.minigitDir, "HEAD")
	content, err := os.ReadFile(headPath)

	if err != nil {
		return "", err
	}

	headStr := strings.TrimSpace(string(content))
	if strings.HasPrefix(headStr, "ref: ") {
		return strings.TrimPrefix(headStr, "ref: "), nil
	}

	// Direct hash reference
	return headStr, nil
}

// Updates a branch to point to a specific commit
func (m *Manager) SetBranch(branch, commit string) error {
	branchPath := filepath.Join(m.refsDir, "heads", branch)
	// 0644 ~ owners can read and write, others can only read
	return os.WriteFile(branchPath, []byte(commit+"\n"), 0644)
}

// Returns the commit hash that a branch points to
func (m *Manager) GetBranch(branch string) (string, error) {
	branchPath := filepath.Join(m.refsDir, "heads", branch)
	content, err := os.ReadFile(branchPath)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}
