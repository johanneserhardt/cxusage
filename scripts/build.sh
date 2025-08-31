#!/bin/bash

# Build script for cxusage with version information

set -e

# Get version info
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")}
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "Building cxusage ${VERSION}"
echo "Build time: ${BUILD_TIME}"
echo "Git commit: ${GIT_COMMIT}"

# Build with version info
go build -ldflags "
    -X 'github.com/johanneserhardt/cxusage/internal/commands.Version=${VERSION}' 
    -X 'github.com/johanneserhardt/cxusage/internal/commands.BuildTime=${BUILD_TIME}' 
    -X 'github.com/johanneserhardt/cxusage/internal/commands.GitCommit=${GIT_COMMIT}'
" -o cxusage ./cmd/cxusage

echo "Build completed successfully!"