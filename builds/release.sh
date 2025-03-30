#!/bin/bash

set -e

# --- Safety check: Must be on main branch ---
REQUIRED_BRANCH="main"
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

if [[ "$CURRENT_BRANCH" != "$REQUIRED_BRANCH" ]]; then
  echo "❌ You must be on the '$REQUIRED_BRANCH' branch to run this script (currently on '$CURRENT_BRANCH')"
  echo "💡 Finish merging your release PR into '$REQUIRED_BRANCH' and run:"
  echo "   git checkout $REQUIRED_BRANCH && git pull origin $REQUIRED_BRANCH"
  exit 1
fi

# --- Parse arguments ---
VERSION=$1
shift

CLEAN_FLAG=""
DEBUG_FLAG=""

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
      echo "Usage: ./builds/release.sh vX.Y.Z [--clean] [--debug]"
      exit 1
      ;;
  esac
  shift
done

if [[ -z "$VERSION" ]]; then
  echo "Usage: ./builds/release.sh vX.Y.Z [--clean] [--debug]"
  exit 1
fi

if [[ "$VERSION" != v* ]]; then
  echo "❌ Version must start with 'v' (e.g., v0.2.1)"
  exit 1
fi

# --- Validate release notes ---
if ! grep -q "$VERSION" RELEASE_NOTES.md; then
  echo "❌ RELEASE_NOTES.md does not mention version $VERSION"
  exit 1
fi

# --- Commit release notes ---
echo "📦 Committing release notes..."
git add RELEASE_NOTES.md
git commit -m "Release $VERSION"

# --- Tag BEFORE building ---
echo "🏷️ Tagging as $VERSION"
git tag "$VERSION"

# --- Run builds ---
echo "🔧 Starting cross-platform build..."

echo "🐧 Building Linux AppImage..."
bash "$(dirname "$0")/build-appimage.sh" $CLEAN_FLAG $DEBUG_FLAG

echo "🪟 Building Windows zip..."
bash "$(dirname "$0")/build-windows.sh" $CLEAN_FLAG $DEBUG_FLAG

echo "✅ All builds completed successfully."

# --- Push everything ---
echo "🚀 Pushing code and tag to origin..."
git push origin main
git push origin "$VERSION"

echo ""
echo "🎉 Release $VERSION is complete and live!"
