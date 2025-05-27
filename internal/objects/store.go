// Object storage system
package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Represents different types of objects
type ObjectType string

const (
	BlobObject   ObjectType = "blob"
	TreeObject   ObjectType = "tree"
	CommitObject ObjectType = "commit"
)

// Represents a minigit object with its metadata
type Object struct {
	Type    ObjectType
	Size    int64
	Content []byte
	Hash    string
}

// Handles all object storage operations
type Store struct {
	objectsDir string
}

func NewStore(minigitDir string) (*Store, error) {
	objectsDir := filepath.Join(minigitDir, "objects")
	return &Store{objectsDir: objectsDir}, nil
}

// Computes the SHA-1 hash of content
func (s *Store) HashContent(objType ObjectType, content []byte) string {
	// Git's object format: "type size\0content"
	contentSize := len(content)
	header := fmt.Sprintf("%s %d\\0", objType, contentSize)
	fullContent := append([]byte(header), content...)

	hash := sha1.Sum(fullContent)
	return fmt.Sprintf("%x", hash)
}

// Saves an object and returns its hash, uses zlib compression (like real Git)
func (s *Store) StoreObject(objType ObjectType, content []byte) (string, error) {
	hash := s.HashContent(objType, content)

	// Create subdir based on first 2 characters of hash
	subDir := filepath.Join(s.objectsDir, hash[:2])
	// 0755 ~~ rwxr-xr-x
	if err := os.MkdirAll(subDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create object directory: %w", err)
	}

	objPath := filepath.Join(subDir, hash[2:])

	// Check if object already exists
	if _, err := os.Stat(objPath); err == nil {
		return hash, nil
	}

	// Prepare content with header
	contentSize := len(content)
	header := fmt.Sprintf("%s %d\\0", objType, contentSize)
	fullContent := append([]byte(header), content...)

	// Compress with zlib
	var compressed bytes.Buffer
	writer := zlib.NewWriter(&compressed)

	if _, err := writer.Write(fullContent); err != nil {
		return "", fmt.Errorf("failed to compress object: %w", err)
	}
	writer.Close()

	// Write compressed content to file
	// 0444 ~ read-only permissions for everyone
	if err := os.WriteFile(objPath, compressed.Bytes(), 0444); err != nil {
		return "", fmt.Errorf("failed to write object file: %w", err)
	}

	return hash, nil
}

// Retrieves an object by its hash
func (s *Store) LoadObject(objHash string) (*Object, error) {
	objPath := filepath.Join(s.objectsDir, objHash[:2], objHash[2:])

	compressedData, err := os.ReadFile(objPath)
	if err != nil {
		return nil, fmt.Errorf("object not found: %w", err)
	}

	// Decompress content
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress object: %w", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %w", err)
	}

	// Parse header
	nullIdx := bytes.IndexByte(content, 0)
	if nullIdx == -1 {
		return nil, fmt.Errorf("invalid object format")
	}

	var objType ObjectType
	var objSize int64
	header := string(content[:nullIdx])
	if _, err := fmt.Sscanf(header, "%s %d", &objType, &objSize); err != nil {
		return nil, fmt.Errorf("failed to parse object header: %w", err)
	}

	objContent := content[nullIdx+1:]

	return &Object{
		Type:    objType,
		Size:    objSize,
		Content: objContent,
		Hash:    objHash,
	}, nil
}
