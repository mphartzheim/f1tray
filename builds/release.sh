#!/bin/bash

set -e

# --- Parse arguments ---
VERSION=$1
shift

CLEAN_FLAG=""
DEBUG_FLAG=""

# Parse remaining flags
while [[ $# -gt 0 ]]; do
  case "$1" in
    --clean)
      CLEAN_FLAG="--clean"
      ;;
    --debug)
      DEBUG_FLAG="--debug"
      ;;
    *)
      echo "❌ Unknown option: $1"
      echo "Usage: ./release.sh vX.Y.Z [--clean] [--debug]"
      exit 1
      ;;
  esac
  shift
done

if [[ -z "$VERSION" ]]; then
  echo "Usage: ./release.sh vX.Y.Z [--clean] [--debug]"
  exit 1
fi

if [[ "$VERSION" != v* ]]; then
  echo "Error: Version must start with 'v' (e.g., v0.2.1)"
  exit 1
fi

# --- Run build ---
echo "🔧 Starting cross-platform build..."

# Build Linux AppImage
echo "🐧 Building Linux AppImage..."
bash "$(dirname "$0")/build-appimage.sh" $CLEAN_FLAG $DEBUG_FLAG

# Build Windows Zip
echo "🪟 Building Windows zip..."
bash "$(dirname "$0")/build-windows.sh" $CLEAN_FLAG $DEBUG_FLAG

echo "✅ All builds completed successfully."

# --- Validate release notes ---
if ! grep -q "$VERSION" RELEASE_NOTES.md; then
  echo "❌ RELEASE_NOTES.md does not mention version $VERSION"
  exit 1
fi

# --- Git commit ---
echo "📦 Committing release notes..."
git add RELEASE_NOTES.md
git commit -m "Release $VERSION"

# --- Tag and push ---
echo "🏷️ Tagging as $VERSION"
git tag "$VERSION"

echo "🚀 Pushing code and tag to origin..."
git push origin main
git push origin "$VERSION"

echo "✅ Release $VERSION is complete and pushed."
