APP_NAME := mygit
CMD_DIR := cmd/minigit
ENTRY := lets.go
BUILD_DIR := bin
MINIGIT_DIR := .minigit

.PHONY: all build run clean

all: build

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/$(ENTRY)

rm_git:
	@rm -rf $(MINIGIT_DIR)
	@echo "Removed .minigit"

clean:
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned build artifacts"
