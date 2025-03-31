#!/bin/bash

set -e

# Handle flags
CLEAN_FLAG=""
DEBUG_FLAG=""

for arg in "$@"; do
  if [[ "$arg" == "--clean" ]]; then
    CLEAN_FLAG="--clean"
  elif [[ "$arg" == "--debug" ]]; then
    DEBUG_FLAG="--debug"
    set -x  # Enable debug output
  fi
done

echo "🔧 Starting cross-platform build..."

# Build Linux AppImage
echo "🐧 Building Linux AppImage..."
bash "$(dirname "$0")/build-appimage.sh" $CLEAN_FLAG $DEBUG_FLAG

# Build Windows Zip
echo "🪟 Building Windows zip..."
bash "$(dirname "$0")/build-windows.sh" $CLEAN_FLAG $DEBUG_FLAG

echo "✅ All builds completed successfully."
