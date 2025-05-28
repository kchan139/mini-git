package repository

import (
	"fmt"
	"os"
)

func (repo *Repository) AddToIndex(path, hash string, info os.FileInfo) error {
	if repo.index == nil {
		return fmt.Errorf("index not initialized")
	}
	return repo.index.AddEntry(path, hash, info)
}
