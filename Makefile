BINARY_NAME := gosqlgen
GO_CMD_PATH := ./cmd
BUILD_DIR := bin

TARGET_LINUX_AMD64  := $(BUILD_DIR)/$(BINARY_NAME).linux-amd64
TARGET_WINDOWS_AMD64 := $(BUILD_DIR)/$(BINARY_NAME).windows-amd64.exe
TARGET_DARWIN_AMD64 := $(BUILD_DIR)/$(BINARY_NAME).darwin-amd64
TARGET_DARWIN_ARM64 := $(BUILD_DIR)/$(BINARY_NAME).darwin-arm64
TARGET_LINUX_ARM64  := $(BUILD_DIR)/$(BINARY_NAME).linux-arm64

ALL_TARGETS := \
    $(TARGET_LINUX_AMD64) \
    $(TARGET_WINDOWS_AMD64) \
    $(TARGET_DARWIN_AMD64) \
    $(TARGET_DARWIN_ARM64) \
    $(TARGET_LINUX_ARM64)

.PHONY: all clean local

all: $(ALL_TARGETS)
	@echo "------------------------------------------------------------------------"
	@echo "Build successful! Binaries are located in the '$(BUILD_DIR)/' directory."
	@echo "------------------------------------------------------------------------"


$(TARGET_LINUX_AMD64):
	@echo "Building $< for Linux (amd64) -> $@"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $@ $(GO_CMD_PATH)

$(TARGET_WINDOWS_AMD64):
	@echo "Building $< for Windows (amd64) -> $@"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $@ $(GO_CMD_PATH)

$(TARGET_DARWIN_AMD64):
	@echo "Building $< for macOS (amd64 - Intel) -> $@"
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $@ $(GO_CMD_PATH)

$(TARGET_DARWIN_ARM64):
	@echo "Building $< for macOS (arm64 - Apple Silicon) -> $@"
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $@ $(GO_CMD_PATH)

$(TARGET_LINUX_ARM64):
	@echo "Building $< for Linux (arm64) -> $@"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -o $@ $(GO_CMD_PATH)

local:
	@echo "Building local $(BINARY_NAME) for testing..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(GO_CMD_PATH)

clean:
	@echo "Cleaning up $(BUILD_DIR)/ directory..."
	@rm -rf $(BUILD_DIR)
