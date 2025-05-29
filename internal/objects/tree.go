package objects

import (
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Represents a file or a subdirectory
type TreeEntry struct {
	Mode os.FileMode `json:"mode"`
	Name string      `json:"name"`
	Hash string      `json:"hash"`
	Type ObjectType  `json:"type"`
}

// Represents a tree object containing files and subdirectories
type Tree struct {
	Entries []TreeEntry `json:"entries"`
}

func (store *Store) CreateTreeFromIndex(entries map[string]*IndexEntry) (string, error) {
	if len(entries) == 0 {
		return "", fmt.Errorf("no entries to create tree from")
	}

	// Build dir structure
	root := &TreeNode{
		name:     "",
		isDir:    true,
		children: make(map[string]*TreeNode),
	}

	// Add all entries to tree
	for path, entry := range entries {
		parts := strings.Split(path, "/")
		current := root

		partsLength := len(parts)
		// Navigate/create dir structure
		for _, part := range parts[:partsLength-1] {
			if current.children[part] == nil {
				current.children[part] = &TreeNode{
					name:     part,
					isDir:    true,
					children: make(map[string]*TreeNode),
				}
			}
			current = current.children[part]
		}

		fileName := parts[partsLength-1]
		current.children[fileName] = &TreeNode{
			name:  fileName,
			isDir: false,
			mode:  entry.Mode,
			hash:  entry.Hash,
		}
	}

	return store.storeTreeNode(root)
}

type TreeNode struct {
	name     string
	isDir    bool
	children map[string]*TreeNode
	mode     os.FileMode
	hash     string
}

type IndexEntry struct {
	Path string      `json:"path"`
	Hash string      `json:"hash"`
	Mode os.FileMode `json:"mode"`
}

func (store *Store) storeTreeNode(node *TreeNode) (string, error) {
	var entries []TreeEntry

	var names []string
	for name := range node.children {
		names = append(names, name)
	}
	sort.Strings(names)

	// Process child
	for _, name := range names {
		child := node.children[name]

		if child.isDir {
			childHash, err := store.storeTreeNode(child)
			if err != nil {
				return "", fmt.Errorf("failed to create subtree for %s: %w", name, err)
			}

			entries = append(entries, TreeEntry{
				Mode: 0755, // Dir mode
				Name: name,
				Hash: childHash,
				Type: TreeObject,
			})
		} else {
			// File entry
			entries = append(entries, TreeEntry{
				Mode: child.mode,
				Name: name,
				Hash: child.hash,
				Type: BlobObject,
			})
		}
	}

	content := store.serializeTree(entries)

	return store.StoreObject(TreeObject, content)
}

func (store *Store) serializeTree(entries []TreeEntry) []byte {
	var content []byte

	for _, entry := range entries {
		// Git tree format: "mode name\0hash"
		// Mode is octal without leading zeros for files, padded for dirs
		var modeStr string
		if entry.Type == TreeObject {
			modeStr = "40000" // Directory
		} else {
			modeStr = fmt.Sprintf("%o", entry.Mode)
		}

		line := fmt.Sprintf("%s %s\x00", modeStr, entry.Name)
		content = append(content, []byte(line)...)

		// Git object hashes (SHA-1) are 20 bytes
		hashBytes, err := hex.DecodeString(entry.Hash)
		if err != nil {
			panic(fmt.Sprintf("invalid hash string %q: %v", entry.Hash, err))
		}

		content = append(content, hashBytes...)
	}

	return content
}
