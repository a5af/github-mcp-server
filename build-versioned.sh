#!/bin/bash
# Auto-versioning build script for github-mcp-server
# Automatically bumps patch version and builds with version in filename

set -e

echo "=== GitHub MCP Server Auto-Versioning Build ==="
echo

# Get the latest version tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.20.0")
echo "Latest version: $LATEST_TAG"

# Extract major.minor.patch
if [[ $LATEST_TAG =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
    MAJOR="${BASH_REMATCH[1]}"
    MINOR="${BASH_REMATCH[2]}"
    PATCH="${BASH_REMATCH[3]}"
else
    echo "Error: Invalid version tag format: $LATEST_TAG"
    exit 1
fi

# Bump patch version
NEW_PATCH=$((PATCH + 1))
NEW_VERSION="v${MAJOR}.${MINOR}.${NEW_PATCH}"

echo "New version: $NEW_VERSION"
echo

# Build the binary
BINARY_NAME="github-mcp-server-${NEW_VERSION}.exe"
echo "Building $BINARY_NAME..."
go build -o "$BINARY_NAME" ./cmd/github-mcp-server

if [ ! -f "$BINARY_NAME" ]; then
    echo "Error: Build failed - $BINARY_NAME not found"
    exit 1
fi

echo "✓ Build successful: $BINARY_NAME ($(du -h "$BINARY_NAME" | cut -f1))"
echo

# Create git tag
echo "Creating git tag $NEW_VERSION..."
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION: Auto-versioned build

Built on: $(date -u '+%Y-%m-%d %H:%M:%S UTC')
Commit: $(git rev-parse --short HEAD)"

echo "✓ Tag created: $NEW_VERSION"
echo

# Copy to bin directory
BIN_DIR="D:/Code/bin"
if [ -d "$BIN_DIR" ]; then
    echo "Copying to $BIN_DIR..."
    cp "$BINARY_NAME" "$BIN_DIR/"
    echo "✓ Copied to $BIN_DIR/$BINARY_NAME"
    echo

    # List all versioned binaries
    echo "Available versions in $BIN_DIR:"
    ls -lh "$BIN_DIR"/github-mcp-server*.exe | awk '{print "  " $9 " (" $5 ")"}'
else
    echo "⚠ Warning: $BIN_DIR not found, binary not copied"
fi

echo
echo "=== Build Complete ==="
echo "Version: $NEW_VERSION"
echo "Binary: $BINARY_NAME"
echo
echo "Next steps:"
echo "1. Push tag to remote: git push origin $NEW_VERSION"
echo "2. Update agent configs to use: D:/Code/bin/$BINARY_NAME"
echo
