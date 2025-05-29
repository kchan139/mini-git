# MiniGit

A minimal Git implementation in Go.

## Features
- `init`: Initialize new repository
- `add`: Stage files/directories
- `commit`: Create commits with `-m` flag
- Basic object storage (blobs, trees, commits)
- Simple staging area management

## Usage
```sh
# Build project
make build

# Initialize repository
./mygit init

# Add files
./mygit add <file>  # Add specific file
./mygit add .       # Add all files

# Commit changes
./mygit commit -m "Commit message"

# Clean build artifacts
make clean
```
