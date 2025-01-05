# Name of the application
APP_NAME := shooter

# Go file to build
SRC_FILE := client/desktop/main.go

# Output directory
OUTPUT_DIR := bin

# OS and architecture combinations
PLATFORMS := windows/amd64 linux/amd64

# Default target
.PHONY: all
all: clean build

# Clean output directory
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(OUTPUT_DIR)

# Build for all platforms
.PHONY: build
build: clean
	@echo "Building for all platforms..."
	@mkdir -p $(OUTPUT_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d/ -f1); \
		ARCH=$$(echo $$platform | cut -d/ -f2); \
		EXT=$$(if [ "$$OS" = "windows" ]; then echo ".exe"; else echo ""; fi); \
		OUTPUT=$(OUTPUT_DIR)/$(APP_NAME)_$$OS_$$ARCH$$EXT; \
		echo "Building $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -o $$OUTPUT $(SRC_FILE); \
	done
	@echo "Build completed!"

# Build for current OS/ARCH
.PHONY: build-local
build-local: clean
	@echo "Building for local platform..."
	@mkdir -p $(OUTPUT_DIR)
	@go build -o $(OUTPUT_DIR)/$(APP_NAME) $(SRC_FILE)
	@echo "Local build completed: $(OUTPUT_DIR)/$(APP_NAME)"
