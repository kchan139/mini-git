// Staging area management
package index

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
)

// Manages the staging area
type Index struct {
	indexPath string
	entries   map[string]*Entry
}

func NewIndex(minigitDir string) (*Index, error) {
	indexPath := filepath.Join(minigitDir, "index")

	index := &Index{
		indexPath: indexPath,
		entries:   make(map[string]*Entry),
	}

	if err := index.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load index: %w", err)
	}

	return index, nil
}

// Adds or updates a file in the staging area
func (idx *Index) AddEntry(path, hash string, info os.FileInfo) error {
	idx.entries[path] = &Entry{
		Path:    path,
		Hash:    hash,
		Mode:    info.Mode(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	return idx.save()
}

// Removes a file from the staging area
func (idx *Index) RemoveEntry(path string) error {
	delete(idx.entries, path)
	return idx.save()
}

// Returns all staged entries
func (idx *Index) GetEntries() map[string]*Entry {
	result := make(map[string]*Entry)
	maps.Copy(result, idx.entries)
	return result
}

// Removes all entries from the staging area
func (idx *Index) Clear() error {
	idx.entries = make(map[string]*Entry)
	return idx.save()
}

// Check if index is empty
func (idx *Index) IsEmpty() bool {
	return len(idx.entries) == 0
}

// Get number of entries
func (idx *Index) Count() int {
	return len(idx.entries)
}

// Utils functions
func (idx *Index) load() error {
	data, err := os.ReadFile(idx.indexPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &idx.entries)
}

func (idx *Index) save() error {
	data, err := json.MarshalIndent(idx.entries, "", "  ")
	if err != nil {
		return err
	}
	// 0644 ~ Owner can read and write, group and others can only read.
	return os.WriteFile(idx.indexPath, data, 0644)
}
