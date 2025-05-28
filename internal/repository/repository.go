// Repository abstraction layer
package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"minigit/internal/config"
	"minigit/internal/index"
	"minigit/internal/objects"
	"minigit/internal/refs"
)

// Repository represents the main interface to a MiniGit repository
// This is our facade that coordinates all the different subsystems
type Repository struct {
	workDir    string
	minigitDir string
	objects    *objects.Store
	index      *index.Index
	refs       *refs.Manager
	config     *config.Config
}

// Creates or opens a repository at a specified path
func NewRepository(path string) (*Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	minigitDir := filepath.Join(absPath, ".minigit")

	repo := &Repository{
		workDir:    absPath,
		minigitDir: minigitDir,
	}

	// Initialize subsystems
	if repo.objects, err = objects.NewStore(minigitDir); err != nil {
		return nil, fmt.Errorf("failed to initialize object store: %w", err)
	}
	if repo.index, err = index.NewIndex(minigitDir); err != nil {
		return nil, fmt.Errorf("failed to initialize index: %w", err)
	}
	if repo.refs, err = refs.NewManager(minigitDir); err != nil {
		return nil, fmt.Errorf("failed to initialzie manager: %w", err)
	}
	if repo.config, err = config.NewConfig(minigitDir); err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	return repo, nil
}

// Creates a new repository structure
func (r *Repository) Initialize() error {
	// Create .minigit directory structure
	dirs := []string{
		r.minigitDir,
		filepath.Join(r.minigitDir, "objects"),
		filepath.Join(r.minigitDir, "refs"),
		filepath.Join(r.minigitDir, "refs", "heads"),
	}

	for _, dir := range dirs {
		// 0755 ~~ rwxr-xr-x
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Initialize HEAD to point to main branch
	return r.refs.SetHead("refs/heads/main")
}
