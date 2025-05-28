APP_NAME := mygit
CMD_DIR := cmd/minigit
ENTRY := lets.go
BUILD_DIR := bin

.PHONY: all build run clean

all: build

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/$(ENTRY)

clean:
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned build artifacts"
