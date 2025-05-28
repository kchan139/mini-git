package repository

import (
	"fmt"
	"minigit/internal/index"
	"minigit/internal/objects"
	"minigit/internal/refs"
	"os"
)

func (repo *Repository) AddToIndex(path, hash string, info os.FileInfo) error {
	if repo.index == nil {
		return fmt.Errorf("index not initialized")
	}
	return repo.index.AddEntry(path, hash, info)
}

func (repo *Repository) GetIndex() (*index.Index, error) {
	if repo.index == nil {
		return nil, fmt.Errorf("index init initialized")
	}
	return repo.index, nil
}

func (repo *Repository) GetObjectStore() (*objects.Store, error) {
	if repo.objects == nil {
		return nil, fmt.Errorf("object store not initialized")
	}
	return repo.objects, nil
}

func (repo *Repository) GetRefsManager() (*refs.Manager, error) {
	if repo.refs == nil {
		return nil, fmt.Errorf("refs manager not initialized")
	}
	return repo.refs, nil
}

func (repo *Repository) GetWorkingDirectory() string {
	return repo.workDir
}

func (repo *Repository) GetMinigitDirectory() string {
	return repo.minigitDir
}
