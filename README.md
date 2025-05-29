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

## Implementation Notes
- Stores data in `.minigit` directory
- Uses SHA-1 hashing for objects
- Compression with zlib
- JSON-based index for staging
- First-parent-only commit history

## Limitations
- No branching/checkout (yet)
- No diff/log/status (yet)
- No remote operations
- Minimal error handling

> Note: Educational project - not for production use.