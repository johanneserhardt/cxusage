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

echo "Installing to $INSTALL_DIR..."
# Install single binary as 'cx'
if [ -L "$INSTALL_DIR/cx" ] || [ -f "$INSTALL_DIR/cx" ]; then
  rm -f "$INSTALL_DIR/cx"
fi
cp cx "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/cx"

# Clean up old 'cxusage' if present to avoid confusion
if [ -e "$INSTALL_DIR/cxusage" ]; then
  echo "Removing legacy binary: $INSTALL_DIR/cxusage"
  rm -f "$INSTALL_DIR/cxusage"
fi

# Check if directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "⚠️  $INSTALL_DIR is not in your PATH"
    echo "Add this to your shell profile (.bashrc, .zshrc, etc.):"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
fi

# Test installation
if command -v cx >/dev/null 2>&1; then
    echo "✅ Installation successful!"
    echo ""
    echo "📦 Installed binaries:"
    echo "  • cx"
    echo ""
    cx version
    echo ""
    echo "Try:"
    echo "  • cx demo"
    echo "  • cx blocks --live"
else
    echo "⚠️  Installation completed but binaries not found in PATH"
    echo "You may need to restart your terminal or add $INSTALL_DIR to PATH"
    echo ""
    echo "Installed:"
    echo "  • $INSTALL_DIR/cx"
fi
