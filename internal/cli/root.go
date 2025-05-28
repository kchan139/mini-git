// Command structure using Cobra pattern
package cli

import (
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Description string
	Handler     func(args []string) error
}

var commands = map[string]Command{
	"init":   {"init", "Initialize a new repository", handleInit},
	"add":    {"add", "Add files to staging area", handleAdd},
	"commit": {"commit", "Create a new commit", handleCommit},
	"status": {"status", "Show repository status", handleStatus},
	"log":    {"log", "Show commit history", handleLog},
}

func Execute() error {
	if len(os.Args) < 2 {
		return showHelp()
	}

	cmdName := os.Args[1]
	cmd, exists := commands[cmdName]

	if !exists {
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	return cmd.Handler(os.Args[2:])
}

func showHelp() error {
	fmt.Println("Usage: minigit <command>")
	fmt.Println("Available commands:")
	for _, cmd := range commands {
		fmt.Printf("	%-7s %s\n", cmd.Name, cmd.Description)
	}
	return nil
}
