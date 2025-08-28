#!/bin/bash

# Simplified build script - for personal project use
set -e

# Project configuration
PROJECT_NAME="file-fusion-rename"
BUILD_DIR="dist"

# Color output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Show help
show_help() {
  echo "Simplified build script - for personal project use"
  echo ""
  echo "Usage:"
  echo "  ./build.sh           # Build current platform"
  echo "  ./build.sh all       # Build all platforms"
  echo "  ./build.sh clean     # Clean build files"
  echo ""
}

# Clean build files
clean() {
  echo -e "${YELLOW}Cleaning build files...${NC}"
  rm -rf "${BUILD_DIR}"
  rm -f "${PROJECT_NAME}"
  echo -e "${GREEN}Clean completed${NC}"
}

# Build current platform
build_local() {
  echo -e "${GREEN}Building current platform...${NC}"
  go build -ldflags "-s -w" -o "${PROJECT_NAME}" main.go
  echo -e "${GREEN}Build completed: ${PROJECT_NAME}${NC}"
}

# Build all platforms
build_all() {
  echo -e "${GREEN}Building all platforms...${NC}"
  mkdir -p "${BUILD_DIR}"
  
  # macOS
  echo -e "${YELLOW}Building macOS...${NC}"
  GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o "${BUILD_DIR}/${PROJECT_NAME}-darwin-amd64" main.go
  GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o "${BUILD_DIR}/${PROJECT_NAME}-darwin-arm64" main.go
  
  # Linux (disable CGO to avoid dependency issues)
  echo -e "${YELLOW}Building Linux...${NC}"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "${BUILD_DIR}/${PROJECT_NAME}-linux-amd64" main.go
  
  # Windows (disable CGO)
  echo -e "${YELLOW}Building Windows...${NC}"
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o "${BUILD_DIR}/${PROJECT_NAME}-windows-amd64.exe" main.go
  
  echo -e "${GREEN}All platforms build completed!${NC}"
  ls -la "${BUILD_DIR}/"
}

# 主逻辑
case "${1:-local}" in
  "all")
    build_all
    ;;
  "clean")
    clean
    ;;
  "help"|"--help"|"-h")
    show_help
    ;;
  *)
    build_local
    ;;
esac