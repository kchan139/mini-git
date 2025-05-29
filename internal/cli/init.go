package cli

import (
	"fmt"
	"minigit/internal/repository"
	"os"
	"path/filepath"
)

func handleInit(args []string) error {
	var targetDir string

	if len(args) > 0 {
		targetDir = args[0]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		targetDir = cwd
	}

	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		// 0755 ~ owners can read, write, execute; others can read and execute
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	minigitDir := filepath.Join(absPath, ".minigit")
	if _, err := os.Stat(minigitDir); err == nil {
		fmt.Printf("Reinitialized existing MiniGit repository in %s\n", absPath)
		return nil
	}

	repo, err := repository.NewRepository(absPath)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	if err := repo.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	fmt.Printf("Initialized empty MiniGit repository in %s\n", minigitDir)
	return nil
}
