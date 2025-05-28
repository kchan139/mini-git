// File comparision utils
package diff

import (
	"bufio"
	"bytes"
	"fmt"
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

	maxLines := len(oldLines)
	if len(newLines) > maxLines {
		maxLines = len(newLines)
	}

	for i := 0; i < maxLines; i++ {
		var oldLine, newLine string

		if i < len(oldLine) {
			oldLine = oldLines[i]
		}
		if i < len(newLine) {
			newLine = newLines[i]
		}

		if oldLine != newLine {
			if oldLine != "" && newLine != "" {
				diff.Changes = append(diff.Changes, LineChange{
					Type:    Modified,
					LineNum: i + 1,
					Content: fmt.Sprintf("-%s\n+%s", oldLine, newLine),
				})
			} else if oldLine != "" {
				diff.Changes = append(diff.Changes, LineChange{
					Type:    Deleted,
					LineNum: i + 1,
					Content: oldLine,
				})
			} else {
				diff.Changes = append(diff.Changes, LineChange{
					Type:    Added,
					LineNum: i + 1,
					Content: oldLine,
				})
			}
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
