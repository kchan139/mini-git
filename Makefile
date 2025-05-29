APP_NAME := mygit
CMD_DIR := cmd/minigit
ENTRY := lets.go
BUILD_DIR := bin
MINIGIT_DIR := .minigit

.PHONY: all build run clean

all: build

build:
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/$(ENTRY)
	@ln -sf $(BUILD_DIR)/$(APP_NAME) $(APP_NAME)

clean:
	@rm -f $(APP_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(MINIGIT_DIR)
	@echo "Cleaned build artifacts, symbolic link, and .minigit folder"
