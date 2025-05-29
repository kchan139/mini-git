package objects

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
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
	commit := &Commit{
		Parents: make([]string, 0),
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	// Parse header lines
	for scanner.Scan() {
		line := scanner.Text()

		// Empty line indicates start of commit message
		if line == "" {
			break
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "tree":
			commit.Tree = parts[1]
		case "parent":
			commit.Parents = append(commit.Parents, parts[1])
		case "author":
			commit.Author = parseAuthorLine(parts[1])
		case "committer":
			authorInfo := parseAuthorLine(parts[1])
			commit.Committer = authorInfo
			// Parse timestamp from committer line
			if timestamp := parseTimestamp(parts[1]); !timestamp.IsZero() {
				commit.Timestamp = timestamp
			}
		}
	}

	// Parse commit message (everything after the empty line)
	var messageLines []string
	for scanner.Scan() {
		messageLines = append(messageLines, scanner.Text())
	}

	if len(messageLines) > 0 {
		commit.Message = strings.Join(messageLines, "\n")
		commit.Message = strings.TrimSpace(commit.Message)
	}

	return commit, nil
}

func parseAuthorLine(line string) string {
	// Format: "Name <email> timestamp timezone"
	// Extract just "Name <email>"
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return line
	}

	// Find the last two parts (timestamp and timezone)
	// Everything before that is the author name and email
	authorParts := parts[:len(parts)-2]
	return strings.Join(authorParts, " ")
}

func parseTimestamp(line string) time.Time {
	// Format: "Name <email> timestamp timezone"
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return time.Time{}
	}

	// Get the timestamp (second to last field)
	timestampStr := parts[len(parts)-2]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(timestamp, 0)
}
