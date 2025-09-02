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

# Copy binaries
echo "Installing to $INSTALL_DIR..."
cp cxusage "$INSTALL_DIR/"
cp cx "$INSTALL_DIR/"

# Make sure they're executable
chmod +x "$INSTALL_DIR/cxusage"
chmod +x "$INSTALL_DIR/cx"

# Check if directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "⚠️  $INSTALL_DIR is not in your PATH"
    echo "Add this to your shell profile (.bashrc, .zshrc, etc.):"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
fi

# Test installation
if command -v cxusage >/dev/null 2>&1 && command -v cx >/dev/null 2>&1; then
    echo "✅ Installation successful!"
    echo ""
    echo "📦 Installed binaries:"
    echo "  • cxusage (full name)"
    echo "  • cx (short alias)"
    echo ""
    cxusage version
    echo ""
    echo "Try:"
    echo "  • cxusage demo  (or cx demo)"
    echo "  • cx blocks --live"
else
    echo "⚠️  Installation completed but binaries not found in PATH"
    echo "You may need to restart your terminal or add $INSTALL_DIR to PATH"
    echo ""
    echo "Installed:"
    echo "  • $INSTALL_DIR/cxusage"
    echo "  • $INSTALL_DIR/cx"
fi