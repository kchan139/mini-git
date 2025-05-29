package objects

import (
	"fmt"
	"time"
)

type Commit struct {
	Tree      string    `json:"tree"`
	Parents   []string  `json:"parents"`
	Author    string    `json:"author"`
	Committer string    `json:"committer"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"time"`
}

// Creates and stores commit objects
func (store *Store) CreateCommit(treeHash string, parents []string, author, message string) (string, error) {
	if treeHash == "" {
		return "", fmt.Errorf("tree hash cannot be empty")
	}
	if message == "" {
		return "", fmt.Errorf("aborting commit due to empty commit message")
	}
	if author == "" {
		author = "MiniGit User <user@minigit.local>"
	}

	commit := &Commit{
		Tree:      treeHash,
		Parents:   parents,
		Author:    author,
		Committer: author, // For simplicity
		Message:   message,
		Timestamp: time.Now(),
	}

	content := store.serializeCommit(commit)

	return store.StoreObject(CommitObject, content)
}

// Serialize to Git's commit format
func (store *Store) serializeCommit(commit *Commit) []byte {
	var content string

	content += fmt.Sprintf("tree %s\n", commit.Tree)

	for _, parent := range commit.Parents {
		content += fmt.Sprintf("parent %s\n", parent)
	}

	timestamp := commit.Timestamp.Unix()
	timezone := commit.Timestamp.Format("-0700")

	content += fmt.Sprintf("author %s %d %s\n", commit.Author, timestamp, timezone)
	content += fmt.Sprintf("committer %s %d %s\n", commit.Committer, timestamp, timezone)

	content += "\n" + commit.Message + "\n"

	return []byte(content)
}

func (store *Store) ParseCommit(content []byte) (*Commit, error) {
	return nil, fmt.Errorf("ParseCommit not implemented")
}
