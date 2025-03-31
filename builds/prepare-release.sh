#!/bin/bash
set -e

# === Config ===
RELEASE_BRANCH="main"
DEVELOP_BRANCH="dev"

# === 1. Ensure clean working directory ===
if [[ -n $(git status --porcelain) ]]; then
  echo "❌ Working directory not clean. Please commit or stash changes first."
  exit 1
fi

# === 2. Ensure you're on the dev branch ===
current_branch=$(git rev-parse --abbrev-ref HEAD)
if [[ "$current_branch" != "$DEVELOP_BRANCH" ]]; then
  echo "❌ You must be on the '$DEVELOP_BRANCH' branch to prepare a release."
  exit 1
fi

# === 3. Ensure branch is up to date with remote ===
git fetch origin
if [[ $(git rev-parse "$DEVELOP_BRANCH") != $(git rev-parse "origin/$DEVELOP_BRANCH") ]]; then
  echo "❌ Local '$DEVELOP_BRANCH' is not up to date with origin. Please pull first."
  exit 1
fi

# === 4. Get latest tag ===
if ! latest_tag=$(git describe --tags --abbrev=0 2>/dev/null); then
  echo "❌ No tags found. Please create a tag for this release (e.g., v1.0.0)."
  exit 1
fi

# === 5. Ensure HEAD is tagged ===
tag_commit=$(git rev-list -n 1 "$latest_tag")
head_commit=$(git rev-parse HEAD)
if [[ "$tag_commit" != "$head_commit" ]]; then
  echo "❌ HEAD is not tagged. Please create a tag on the latest commit."
  echo "Suggestion: git tag -a vX.Y.Z -m 'Release vX.Y.Z' && git push --tags"
  exit 1
fi

# === 6. Optional: Confirm changelog contains latest tag ===
if ! grep -q "$latest_tag" CHANGELOG.md 2>/dev/null; then
  echo "⚠️  WARNING: Tag $latest_tag not found in CHANGELOG.md. Did you forget to update it?"
fi

# === 7. Create PR from dev → main ===
echo "🚀 Creating a pull request from $DEVELOP_BRANCH into $RELEASE_BRANCH for tag $latest_tag..."
gh pr create --base "$RELEASE_BRANCH" --head "$DEVELOP_BRANCH" --title "Release $latest_tag" --body "Automated PR to release version $latest_tag"

echo "✅ Release PR created! Merge it once approved."
