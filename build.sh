#!/bin/bash
# Script to compile TinyMonitor into a standalone binary

# Clean up previous builds
rm -rf build dist tinymonitor.spec

# Install build dependencies if needed
# pip install pyinstaller

# Compilation
echo "🔨 Building TinyMonitor..."
pyinstaller --name tinymonitor --onefile --paths src src/tinymonitor/__main__.py

echo "✅ Build complete!"
echo "Binary is located at: dist/tinymonitor"

# Quick test
echo "🧪 Testing binary version..."
./dist/tinymonitor --help
