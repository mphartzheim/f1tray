#!/bin/bash

set -e

REQUIRED_BRANCH="main"
CURRENT_BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null || echo "detached")

if [[ "$CURRENT_BRANCH" != "$REQUIRED_BRANCH" ]]; then
  echo "❌ You must be on the '$REQUIRED_BRANCH' branch to run this script (currently on '$CURRENT_BRANCH')"
  echo "💡 Finish merging your release PR into '$REQUIRED_BRANCH' and run:"
  echo "   git checkout $REQUIRED_BRANCH && git pull origin $REQUIRED_BRANCH"
  exit 1
fi

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

# --- Commit if not already committed ---
LAST_COMMIT_MSG=$(git log -1 --pretty=%s)
if [[ "$LAST_COMMIT_MSG" == "Release $VERSION" ]]; then
  echo "ℹ️ Release commit for $VERSION already exists. Skipping commit."
else
  echo "📦 Committing release notes..."
  git add RELEASE_NOTES.md
  git commit -m "Release $VERSION"
fi

# --- Tag if not already tagged ---
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "ℹ️ Tag '$VERSION' already exists. Skipping tag."
else
  echo "🏷️ Tagging as $VERSION"
  git tag "$VERSION"
fi

# --- Always build, even if commit/tag already exist ---
echo "🔧 Starting cross-platform build..."

echo "🐧 Building Linux AppImage..."
bash "$(dirname "$0")/build-appimage.sh" $CLEAN_FLAG $DEBUG_FLAG

echo "🪟 Building Windows zip..."
bash "$(dirname "$0")/build-windows.sh" $CLEAN_FLAG $DEBUG_FLAG

echo "✅ All builds completed successfully."

# --- Push code and tag ---
echo "🚀 Pushing code and tag to origin..."
git push origin main || echo "⚠️ Could not push main (may already be up to date)"
git push origin "$VERSION" || echo "⚠️ Could not push tag (may already be pushed)"

echo ""
echo "🎉 Release $VERSION is complete!"
