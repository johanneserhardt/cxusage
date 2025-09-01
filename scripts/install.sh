#!/bin/bash

# Install script for cxusage

set -e

echo "🚀 Installing cxusage..."

# Build first
echo "Building cxusage..."
./scripts/build.sh

# Determine install location
if [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
elif [ -d "$HOME/bin" ]; then
    INSTALL_DIR="$HOME/bin"
elif [ -d "$HOME/.local/bin" ]; then
    INSTALL_DIR="$HOME/.local/bin"
else
    echo "Creating $HOME/.local/bin directory..."
    mkdir -p "$HOME/.local/bin"
    INSTALL_DIR="$HOME/.local/bin"
fi

# Copy binary
echo "Installing to $INSTALL_DIR..."
cp cxusage "$INSTALL_DIR/"

# Make sure it's executable
chmod +x "$INSTALL_DIR/cxusage"

# Check if directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "⚠️  $INSTALL_DIR is not in your PATH"
    echo "Add this to your shell profile (.bashrc, .zshrc, etc.):"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
fi

# Test installation
if command -v cxusage >/dev/null 2>&1; then
    echo "✅ Installation successful!"
    echo ""
    cxusage version
    echo ""
    echo "Try: cxusage demo"
else
    echo "⚠️  Installation completed but cxusage not found in PATH"
    echo "You may need to restart your terminal or add $INSTALL_DIR to PATH"
fi