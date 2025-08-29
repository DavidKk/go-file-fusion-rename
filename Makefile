# 项目配置
PROJECT_NAME := file-fusion-rename
BINARY_NAME := file-fusion-rename
MAIN_FILE := main.go

# Go 配置
GO := go
LDFLAGS := -s -w

# 目录配置
TMP_DIR := tmp

# 颜色输出
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m

# 默认目标
.DEFAULT_GOAL := help

# 确保目录存在
$(TMP_DIR):
	@mkdir -p $(TMP_DIR)

# Development mode - hot reload
.PHONY: dev
dev:
	@echo "$(GREEN)Starting hot reload development mode...$(NC)"
	@air

# Build for macOS (both architectures)
.PHONY: build-macos
build-macos:
	@echo "$(GREEN)Building macOS applications...$(NC)"
	@CGO_ENABLED=1 GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-macos-intel $(MAIN_FILE)
	@CGO_ENABLED=1 GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-macos-arm64 $(MAIN_FILE)
	@echo "$(GREEN)macOS Build completed: $(BINARY_NAME)-macos-intel, $(BINARY_NAME)-macos-arm64$(NC)"

# Build for macOS Intel
.PHONY: build-macos-intel
build-macos-intel:
	@echo "$(GREEN)Building macOS Intel application...$(NC)"
	@CGO_ENABLED=1 GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-macos-intel $(MAIN_FILE)
	@echo "$(GREEN)macOS Intel Build completed: $(BINARY_NAME)-macos-intel$(NC)"

# Build for macOS ARM64
.PHONY: build-macos-arm64
build-macos-arm64:
	@echo "$(GREEN)Building macOS ARM64 application...$(NC)"
	@CGO_ENABLED=1 GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-macos-arm64 $(MAIN_FILE)
	@echo "$(GREEN)macOS ARM64 Build completed: $(BINARY_NAME)-macos-arm64$(NC)"

# Build for current platform
.PHONY: build-native
build-native:
	@echo "$(GREEN)Building native application...$(NC)"
	@CGO_ENABLED=1 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-native $(MAIN_FILE)
	@echo "$(GREEN)Native Build completed: $(BINARY_NAME)-native$(NC)"

# Build
.PHONY: build
build:
	@echo "$(GREEN)Building application...$(NC)"
	@$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "$(GREEN)Build completed: $(BINARY_NAME)$(NC)"

# Run
.PHONY: run
run:
	@echo "$(GREEN)Running application...$(NC)"
	@$(GO) run $(MAIN_FILE)

# Test
.PHONY: test
test:
	@echo "$(GREEN)Running tests...$(NC)"
	@$(GO) test -v ./...

# Code formatting
.PHONY: fmt
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	@$(GO) fmt ./...

# Code check
.PHONY: vet
vet:
	@echo "$(GREEN)Running code check...$(NC)"
	@$(GO) vet ./...

# Clean
.PHONY: clean
clean:
	@echo "$(YELLOW)Cleaning files...$(NC)"
	@rm -rf $(TMP_DIR)/
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-macos-intel $(BINARY_NAME)-macos-arm64 $(BINARY_NAME)-native
	@echo "$(GREEN)Clean completed$(NC)"

# Install dependencies
.PHONY: deps
deps:
	@echo "$(GREEN)Installing dependencies...$(NC)"
	@$(GO) mod download
	@$(GO) mod tidy

# Help information
.PHONY: help
help:
	@echo "$(BLUE)================================$(NC)"
	@echo "$(BLUE)    $(PROJECT_NAME) Development Tools$(NC)"
	@echo "$(BLUE)================================$(NC)"
	@echo
	@echo "$(YELLOW)Common Commands:$(NC)"
	@echo "  $(GREEN)make dev$(NC)             - Start hot reload development mode"
	@echo "  $(GREEN)make build$(NC)           - Build application"
	@echo "  $(GREEN)make build-macos$(NC)     - Build for macOS (both architectures)"
	@echo "  $(GREEN)make build-macos-intel$(NC)- Build for macOS Intel"
	@echo "  $(GREEN)make build-macos-arm64$(NC)- Build for macOS ARM64"
	@echo "  $(GREEN)make build-native$(NC)    - Build for current platform"
	@echo "  $(GREEN)make run$(NC)             - Run application"
	@echo "  $(GREEN)make test$(NC)            - Run tests"
	@echo "  $(GREEN)make clean$(NC)           - Clean files"
	@echo
	@echo "$(YELLOW)Code Quality:$(NC)"
	@echo "  $(GREEN)make fmt$(NC)     - Format code"
	@echo "  $(GREEN)make vet$(NC)     - Code check"
	@echo
	@echo "$(YELLOW)Dependency Management:$(NC)"
	@echo "  $(GREEN)make deps$(NC)    - Install dependencies"
	@echo