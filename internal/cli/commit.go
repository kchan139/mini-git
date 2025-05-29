package cli

import (
	"fmt"
	"minigit/internal/diff"
	"minigit/internal/objects"
	"strings"
)

func handleCommit(args []string) error {
	var message string
	argsLen := len(args)

	if argsLen == 0 {
		return fmt.Errorf("interactive commit (opening editor) not implemented. please use `-m` to pass a message")
	}

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

	// Create tree from current index
	indexEntries := make(map[string]*objects.IndexEntry)
	for path, entry := range entries {
		indexEntries[path] = &objects.IndexEntry{
			Path: entry.Path,
			Hash: entry.Hash,
			Mode: entry.Mode,
		}
	}

	newTreeHash, err := store.CreateTreeFromIndex(indexEntries)
	if err != nil {
		return fmt.Errorf("failed to create tree: %w", err)
	}

	// Check if we have changes compared to the last commit
	refsMan, err := repo.GetRefsManager()
	if err != nil {
		return fmt.Errorf("failed to get refs manager: %w", err)
	}

	headRef, err := refsMan.GetHead()
	var lastCommitTreeHash string

	if err == nil {
		var lastCommitHash string
		if strings.HasPrefix(headRef, "refs/heads/") {
			// HEAD points to a branch
			branchName := strings.TrimPrefix(headRef, "refs/heads/")
			if commit, err := refsMan.GetBranch(branchName); err == nil {
				lastCommitHash = commit
			}
		} else {
			// HEAD points to a commit directly
			lastCommitHash = headRef
		}

		// Get the tree hash from the last commit
		if lastCommitHash != "" {
			if lastCommitObj, err := store.LoadObject(lastCommitHash); err == nil {
				if lastCommit, err := store.ParseCommit(lastCommitObj.Content); err == nil {
					lastCommitTreeHash = lastCommit.Tree
				}
			}
		}
	}

	// If the tree hasn't changed, don't create a new commit
	if lastCommitTreeHash == newTreeHash {
		fmt.Println("On branch main")
		fmt.Println("nothing to commit, working tree clean")
		return nil
	}

	var parents []string
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

	commitHash, err := store.CreateCommit(newTreeHash, parents, "", message)
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

	// Calculate diff statistics
	stats, err := diff.CalculateDiffStats(repo, entries, lastCommitTreeHash)
	if err != nil {
		return fmt.Errorf("failed to calculate diff stats: %w", err)
	}

	// Print commit information (mimic real Git)
	shortHash := commitHash
	if len(commitHash) > 7 {
		shortHash = commitHash[:7]
	}

	fmt.Printf("[%s] %s\n", shortHash, message)

	// Print diff statistics
	filesChanged := len(entries)
	if filesChanged == 1 {
		fmt.Printf(" %d file changed", filesChanged)
	} else {
		fmt.Printf(" %d files changed", filesChanged)
	}

	if stats.Insertions > 0 {
		fmt.Printf(", %d insertion(+)", stats.Insertions)
		if stats.Insertions > 1 {
			fmt.Print("s")
		}
	}

	if stats.Deletions > 0 {
		fmt.Printf(", %d deletion(-)", stats.Deletions)
		if stats.Deletions > 1 {
			fmt.Print("s")
		}
	}

	fmt.Println()

	return nil
}
