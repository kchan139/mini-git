package cli

import (
	"fmt"
	"minigit/internal/objects"
	"minigit/internal/repository"
	"os"
	"path/filepath"
	"strings"
)

func handleAdd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("nothing specified, nothing added")
	}

	repo, err := findRepository()
	if err != nil {
		return err
	}

	for _, arg := range args {
		// Special case for "."
		if arg == "." {
			repoRoot := repo.GetWorkingDirectory()
			err := addFile(repo, repoRoot)
			if err != nil {
				return fmt.Errorf("failed to add directory: %w", err)
			}
			continue
		}

		if err := addFile(repo, arg); err != nil {
			return fmt.Errorf("failed to add '%s': %w", arg, err)
		}
	}

	return nil
}

func findRepository() (*repository.Repository, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	dir := cwd
	for {
		minigitDir := filepath.Join(dir, ".minigit")
		if _, err := os.Stat(minigitDir); err == nil {
			return repository.NewRepository(dir)
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return nil, fmt.Errorf("fatal: not a minigit repository (or any of the parent directories): .minigit")
}

func addFile(repo *repository.Repository, filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("fatal: pathspec '%s' did not match any files", filePath)
	}

	if info.IsDir() {
		return addDirectory(repo, absPath)
	}

	return addSingleFile(repo, absPath, info)
}

func addDirectory(repo *repository.Repository, dirPath string) error {
	repoRoot := repo.GetWorkingDirectory()
	if repoRoot == "" {
		return fmt.Errorf("not in a repository")
	}

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and .*git/
		if info.IsDir() || strings.Contains(path, ".minigit") || strings.Contains(path, ".git") {
			return nil
		}

		return addSingleFile(repo, path, info)
	})
}

func addSingleFile(repo *repository.Repository, absPath string, info os.FileInfo) error {
	content, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Store as blob object and get hash
	store, err := repo.GetObjectStore()
	if err != nil {
		return fmt.Errorf("failed to access object store: %w", err)
	}

	hash, err := store.StoreObject(objects.BlobObject, content)
	if err != nil {
		return fmt.Errorf("failed to store object: %w", err)
	}

	repoRoot := repo.GetWorkingDirectory()
	relPath, err := filepath.Rel(repoRoot, absPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	if err := repo.AddToIndex(relPath, hash, info); err != nil {
		return fmt.Errorf("failed to add to index: %w", err)
	}

	return nil
}
