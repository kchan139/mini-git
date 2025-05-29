package cli

import (
	"fmt"
	"minigit/internal/objects"
	"strings"
)

func handleCommit(args []string) error {
	var message string
	argsLen := len(args)

	for i, arg := range args {
		if arg == "-m" && i+1 < argsLen {
			message = args[i+1]
			break
		}
	}

	if message == "" {
		return fmt.Errorf("switch `m' requires a value")
	}

	repo, err := findRepository()
	if err != nil {
		return err
	}

	index, err := repo.GetIndex()
	if err != nil {
		return fmt.Errorf("failed to get index: %w", err)
	}

	entries := index.GetEntries()
	if len(entries) == 0 {
		return fmt.Errorf("no changes added to commit (use \"mygit add\")")
	}

	store, err := repo.GetObjectStore()
	if err != nil {
		return err
	}

	indexEntries := make(map[string]*objects.IndexEntry)
	// Convert indexEntries format to match CreateTreeFromIndex return type
	for path, entry := range entries {
		indexEntries[path] = &objects.IndexEntry{
			Path: entry.Path,
			Hash: entry.Hash,
			Mode: entry.Mode,
		}
	}

	treeHash, err := store.CreateTreeFromIndex(indexEntries)
	if err != nil {
		return fmt.Errorf("failed to create tree: %w", err)
	}

	var parents []string
	refsMan, err := repo.GetRefsManager()
	if err != nil {
		return fmt.Errorf("failed to get refs manager: %w", err)
	}

	// Get current HEAD
	headRef, err := refsMan.GetHead()

	if err == nil && !strings.HasPrefix(headRef, "refs/heads/") {
		// HEAD points to a commit directly
		parents = append(parents, headRef)
	} else {
		// HEAD points to a branch
		branchName := strings.TrimPrefix(headRef, "refs/heads/")
		if currCommit, err := refsMan.GetBranch(branchName); err == nil {
			parents = append(parents, currCommit)
		}
		// If branch doesn't exist yet, parents will be empty (first commit)
	}

	commitHash, err := store.CreateCommit(treeHash, parents, "", message)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Update current branch
	headRef, err = refsMan.GetHead()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	if strings.HasPrefix(headRef, "refs/heads/") {
		// Update branch
		branchName := strings.TrimPrefix(headRef, "refs/heads/")
		if err := refsMan.SetBranch(branchName, commitHash); err != nil {
			return fmt.Errorf("failed to update branch: %w", err)
		}
	} else {
		// Update HEAD directly
		if err := refsMan.SetHead(commitHash); err != nil {
			return fmt.Errorf("failed to update HEAD: %w", err)
		}
	}

	// Clear the idx after successful commit
	if err := index.Clear(); err != nil {
		return fmt.Errorf("failed to clear index: %w", err)
	}

	// Print commit information (mimic real Git)
	shortHash := commitHash
	if len(commitHash) > 7 {
		shortHash = commitHash[:7]
	}

	fmt.Printf("[%s] %s\n", shortHash, message)
	fmt.Printf("%d file(s) changed\n", len(entries))

	return nil
}
