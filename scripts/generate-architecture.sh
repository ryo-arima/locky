#!/bin/bash

# Generate architecture SVG from mermaid file
# Requires: npm install -g @mermaid-js/mermaid-cli

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DOCS_DIR="$PROJECT_ROOT/docs/architecture"
MMD_FILE="$DOCS_DIR/high-level-architecture.mmd"
SVG_FILE="$DOCS_DIR/high-level-architecture.svg"

# Check if mmdc is installed
if ! command -v mmdc &> /dev/null; then
    echo "Error: mmdc (mermaid-cli) is not installed."
    echo "Please install it with:"
    echo "  npm install -g @mermaid-js/mermaid-cli"
    exit 1
fi

# Generate SVG
echo "Generating architecture diagram..."
mmdc -i "$MMD_FILE" -o "$SVG_FILE" -b transparent

echo "Successfully generated: $SVG_FILE"
