#!/bin/bash

set -e

VERSION=$1

if [[ -z "$VERSION" ]]; then
  echo "Usage: ./undo-release.sh vX.Y.Z"
  exit 1
fi

# Confirm before proceeding
read -rp "⚠️  This will delete tag '$VERSION' locally and on GitHub. Continue? (y/N): " CONFIRM
if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
  echo "❌ Aborted."
  exit 1
fi

# Delete tag locally
echo "🗑️  Deleting local tag: $VERSION"
git tag -d "$VERSION" || echo "⚠️  Local tag not found"

# Delete tag on remote
echo "🌐 Deleting remote tag: $VERSION"
git push origin :refs/tags/"$VERSION" || echo "⚠️  Remote tag not found"

# Optionally undo the release commit if it's HEAD
HEAD_MSG=$(git log -1 --pretty=%s)
EXPECTED_MSG="Release $VERSION"

if [[ "$HEAD_MSG" == "$EXPECTED_MSG" ]]; then
  echo "↩️  Undoing release commit at HEAD..."
  git reset --hard HEAD~1
  git push origin main --force
else
  echo "ℹ️  Release commit not at HEAD. Skipping commit rollback."
  echo "🧠 If needed, manually revert or reset to remove the commit."
fi

echo "✅ Release $VERSION undone."
