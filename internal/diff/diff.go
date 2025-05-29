// File comparision utils
package diff

import (
	"bufio"
	"bytes"
	"fmt"
	"minigit/internal/index"
	"minigit/internal/objects"
	"minigit/internal/repository"
	"os"
	"path/filepath"
	"strings"
)

type LineChange struct {
	Type    ChangeType
	LineNum int
	Content string
}

type ChangeType int

const (
	Added    ChangeType = iota // 0
	Deleted                    // 1
	Modified                   // 2
)

// Represents the differences between 2 files
type FileDiff struct {
	OldPath string
	NewPath string
	Changes []LineChange
}

// Compares 2 files and returns the differences
func CompareFiles(oldContent, newContent []byte) *FileDiff {
	oldLines := splitLines(oldContent)
	newLines := splitLines(newContent)

	diff := &FileDiff{
		Changes: make([]LineChange, 0),
	}

	// Simple line-by-line comparison
	oldLen := len(oldLines)
	newLen := len(newLines)
	maxLen := max(oldLen, newLen)

	for i := range maxLen {
		var oldLine, newLine string
		hasOld := i < oldLen
		hasNew := i < newLen

		if hasOld {
			oldLine = oldLines[i]
		}
		if hasNew {
			newLine = newLines[i]
		}

		if hasOld && hasNew {
			if oldLine != newLine {
				diff.Changes = append(diff.Changes, LineChange{
					Type:    Modified,
					LineNum: i + 1,
					Content: fmt.Sprintf("-%s\n+%s", oldLine, newLine),
				})
			}
			// If lines are equal, no change to record
		} else if hasOld && !hasNew {
			// Line was deleted
			diff.Changes = append(diff.Changes, LineChange{
				Type:    Deleted,
				LineNum: i + 1,
				Content: oldLine,
			})
		} else if !hasOld && hasNew {
			// Line was added
			diff.Changes = append(diff.Changes, LineChange{
				Type:    Added,
				LineNum: i + 1,
				Content: newLine,
			})
		}
	}

	return diff
}

func splitLines(content []byte) []string {
	if len(content) == 0 {
		return []string{}
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	lines := make([]string, 0)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

type DiffStats struct {
	Insertions int
	Deletions  int
}

func CalculateDiffStats(repo *repository.Repository, entries map[string]*index.Entry, lastTreeHash string) (*DiffStats, error) {
	stats := &DiffStats{}

	// If this is the first commit, count all lines as insertions
	if lastTreeHash == "" {
		for _, entry := range entries {
			absPath := filepath.Join(repo.GetWorkingDirectory(), entry.Path)
			if content, err := os.ReadFile(absPath); err == nil {
				lines := strings.Count(string(content), "\n")
				if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
					lines++ // Count last line if file doesn't end with newline
				}
				stats.Insertions += lines
			}
		}
		return stats, nil
	}

	// Get previous commit's file contents
	store, err := repo.GetObjectStore()
	if err != nil {
		return stats, err
	}

	previousFiles := make(map[string][]byte)
	if err := getFilesFromTree(store, lastTreeHash, "", previousFiles); err != nil {
		return stats, err
	}

	// Compare each file
	for _, entry := range entries {
		absPath := filepath.Join(repo.GetWorkingDirectory(), entry.Path)
		newContent, err := os.ReadFile(absPath)
		if err != nil {
			continue
		}

		oldContent, existed := previousFiles[entry.Path]
		if !existed {
			// New file - count all lines as insertions
			lines := strings.Count(string(newContent), "\n")
			if len(newContent) > 0 && !strings.HasSuffix(string(newContent), "\n") {
				lines++
			}
			stats.Insertions += lines
		} else {
			// Compare files
			fileDiff := CompareFiles(oldContent, newContent)
			for _, change := range fileDiff.Changes {
				switch change.Type {
				case Added:
					stats.Insertions++
				case Deleted:
					stats.Deletions++
				case Modified:
					// For modified lines, count as both deletion and insertion
					stats.Deletions++
					stats.Insertions++
				}
			}
		}
	}

	return stats, nil
}

func getFilesFromTree(store *objects.Store, treeHash, basePath string, files map[string][]byte) error {
	if treeHash == "" {
		return nil
	}

	treeObj, err := store.LoadObject(treeHash)
	if err != nil {
		return err
	}

	tree, err := store.ParseTree(treeObj.Content)
	if err != nil {
		return err
	}

	for _, entry := range tree.Entries {
		fullPath := entry.Name
		if basePath != "" {
			fullPath = filepath.Join(basePath, entry.Name)
		}

		if entry.Type == objects.BlobObject {
			// It's a file
			blobObj, err := store.LoadObject(entry.Hash)
			if err != nil {
				continue
			}
			files[fullPath] = blobObj.Content
		} else if entry.Type == objects.TreeObject {
			// It's a directory, recurse
			if err := getFilesFromTree(store, entry.Hash, fullPath, files); err != nil {
				return err
			}
		}
	}

	return nil
}
