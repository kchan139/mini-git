package index

import (
	"os"
	"time"
)

// Represents a file in the staging area
type Entry struct {
	Path    string      `json:"path"`
	Hash    string      `json:"hash"`
	Mode    os.FileMode `json:"mode"`
	Size    int64       `json:"size"`
	ModTime time.Time   `json:"mod_time"`
}
