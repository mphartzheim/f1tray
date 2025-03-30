#!/bin/bash

set -e

# --- Smarter branch detection ---
REQUIRED_BRANCH="main"
CURRENT_BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null || echo "detached")

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
  echo "❌ Version must start with 'v' (e.g., v0.2.1)"
  exit 1
fi

# --- Validate release notes ---
if ! grep -q "$VERSION" RELEASE_NOTES.md; then
  echo "❌ RELEASE_NOTES.md does not mention version $VERSION"
  exit 1
fi

# --- Commit release notes if needed ---
LAST_COMMIT_MSG=$(git log -1 --pretty=%s)
if [[ "$LAST_COMMIT_MSG" == "Release $VERSION" ]]; then
  echo "ℹ️ Release commit already exists. Skipping commit step."
else
  echo "📦 Committing release notes..."
  git add RELEASE_NOTES.md
  git commit -m "Release $VERSION"
fi

# --- Tag if it doesn't already exist ---
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "ℹ️ Tag '$VERSION' already exists. Skipping tag creation."
else
  echo "🏷️ Tagging as $VERSION"
  git tag "$VERSION"
fi

# --- Run builds ---
echo "🔧 Starting cross-platform build..."

echo "🐧 Building Linux AppImage..."
bash "$(dirname "$0")/build-appimage.sh" $CLEAN_FLAG $DEBUG_FLAG

echo "🪟 Building Windows zip..."
bash "$(dirname "$0")/build-windows.sh" $CLEAN_FLAG $DEBUG_FLAG

echo "✅ All builds completed successfully."

# --- Push if needed ---
echo "🚀 Pushing code and tag to origin..."

if [[ $(git status --porcelain) ]]; then
  echo "⚠️ Working directory is dirty. Skipping push to avoid conflict."
else
  git push origin main || echo "⚠️ Failed to push main. Handle manually if needed."
  git push origin "$VERSION" || echo "⚠️ Failed to push tag '$VERSION'. It may already be pushed."
fi

echo ""
echo "🎉 Release $VERSION complete (or already finished)."
