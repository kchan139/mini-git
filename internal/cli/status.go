package cli

import (
	"fmt"
	"minigit/internal/objects"
	"minigit/internal/repository"
	"os"
	"path/filepath"
	"strings"
)

func handleStatus(args []string) error {
	repo, err := findRepository()
	if err != nil {
		return err
	}

	index, err := repo.GetIndex()
	if err != nil {
		return fmt.Errorf("failed to get index: %w", err)
	}

	refsMan, err := repo.GetRefsManager()
	if err != nil {
		return fmt.Errorf("failed to get refs manager: %w", err)
	}

	headRef, err := refsMan.GetHead()
	var branchName string
	if err == nil && strings.HasPrefix(headRef, "refs/heads/") {
		branchName = strings.TrimPrefix(headRef, "refs/heads/")
	} else {
		branchName = "HEAD detached"
	}

	fmt.Printf("On branch %s\n", branchName)

	stagedFiles := index.GetEntries()

	workingFiles, err := getWorkdingDirectory(repo.GetWorkingDirectory())
	if err != nil {
		return fmt.Errorf("failed to scan working directory: %w", err)
	}

	// get last commit files for comparision
	lastCommitFiles, err := getLastCommitFiles(repo)
	if err != nil && !isFirstCommit(err) {
		return fmt.Errorf("failed to get last commit files: %w", err)
	}

	// categorize
	var stagedForCommit,
		modifiedFiles,
		untrackedFiles []string

	for path := range stagedFiles {
		stagedForCommit = append(stagedForCommit, path)
	}

	// check workDir files
	for path, currHash := range workingFiles {
		if _, isStaged := stagedFiles[path]; isStaged {
			continue // already handled as staged
		}

		if lastCommitHash, existsInLastCommit := lastCommitFiles[path]; existsInLastCommit {
			if currHash != lastCommitHash {
				modifiedFiles = append(modifiedFiles, path)
			}
		} else {
			untrackedFiles = append(untrackedFiles, path)
		}
	}

	if len(stagedForCommit) == 0 && len(modifiedFiles) == 0 && len(untrackedFiles) == 0 {
		fmt.Println("nothing to commit, working tree clean")
		return nil
	}

	if len(stagedForCommit) > 0 {
		fmt.Println("Changes to be committed:")
		fmt.Println("  (use \"./mygit restore --staged <file>...\" to unstage)")

		for _, file := range stagedForCommit {
			if _, existsInLastCommit := lastCommitFiles[file]; existsInLastCommit {
				fmt.Printf("\tmodified:   %s\n", file)
			} else {
				fmt.Printf("\tnew file:   %s\n", file)
			}
		}
		fmt.Println()
	}

	if len(modifiedFiles) > 0 {
		fmt.Println("Changes not staged for commit:")
		fmt.Println("  (use \"./mygit add <file>...\" to update what will be committed)")
		fmt.Println("  (use \"./mygit checkout -- <file>...\" to discard changes in working directory)")

		for _, file := range modifiedFiles {
			fmt.Printf("\tmodified:   %s\n", file)
		}
		fmt.Println()
	}

	if len(untrackedFiles) > 0 {
		fmt.Println("Untracked files:")
		fmt.Println("  (use \"./mygit add <file>...\" to include in what will be committed)")

		for _, file := range untrackedFiles {
			fmt.Printf("\t%s\n", file)
		}
		fmt.Println()
	}

	return nil
}

func getWorkdingDirectory(workDir string) (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.Walk(workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip .minigit
		if info.IsDir() || strings.Contains(path, ".minigit") || strings.Contains(path, ".git") {
			return nil
		}

		relPath, err := filepath.Rel(workDir, path)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hash := calculateFileHash(content)
		files[relPath] = hash

		return nil
	})

	return files, err
}

func calculateFileHash(content []byte) string {
	store := &objects.Store{}
	return store.HashContent(objects.BlobObject, content)
}

func getLastCommitFiles(repo *repository.Repository) (map[string]string, error) {
	refsMan, err := repo.GetRefsManager()
	if err != nil {
		return nil, err
	}

	headRef, err := refsMan.GetHead()
	if err != nil {
		return nil, err
	}

	var lastCommitHash string
	if strings.HasPrefix(headRef, "refs/heads/") {
		branchName := strings.TrimPrefix(headRef, "refs/heads/")
		if commit, err := refsMan.GetBranch(branchName); err == nil {
			lastCommitHash = commit
		} else {
			return make(map[string]string), nil // first commit
		}
	} else {
		lastCommitHash = headRef
	}

	if lastCommitHash == "" {
		return make(map[string]string), nil
	}

	store, err := repo.GetObjectStore()
	if err != nil {
		return nil, err
	}

	commitObj, err := store.LoadObject(lastCommitHash)
	if err != nil {
		return nil, err
	}

	commit, err := store.ParseCommit(commitObj.Content)
	if err != nil {
		return nil, err
	}

	files := make(map[string]string)
	err = getFilesFromTree(store, commit.Tree, "", files)
	return files, err
}

func getFilesFromTree(store *objects.Store, treeHash, basePath string, files map[string]string) error {
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
			files[fullPath] = entry.Hash
		} else if entry.Type == objects.TreeObject {
			// recurse
			if err := getFilesFromTree(store, entry.Hash, fullPath, files); err != nil {
				return err
			}
		}
	}

	return nil
}

func isFirstCommit(err error) bool {
	return strings.Contains(err.Error(), "object not found") ||
		strings.Contains(err.Error(), "no such file or directory")
}
