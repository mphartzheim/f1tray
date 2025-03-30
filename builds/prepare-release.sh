#!/bin/bash

set -e

# --- Confirm you're on dev branch ---
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != "dev" ]]; then
  echo "❌ You must run this script from the 'dev' branch (currently on '$CURRENT_BRANCH')"
  exit 1
fi

# --- Get release version ---
read -rp "🔖 Enter release version (e.g., v0.2.1): " VERSION

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([-a-zA-Z0-9]*)?$ ]]; then
  echo "❌ Invalid version format. Use 'vX.Y.Z' or 'vX.Y.Z-suffix'"
  exit 1
fi

# --- Optional flags ---
CLEAN_FLAG=""
DEBUG_FLAG=""

read -rp "🧼 Run clean build? (y/N): " CLEAN_INPUT
[[ "$CLEAN_INPUT" =~ ^[Yy]$ ]] && CLEAN_FLAG="--clean"

read -rp "🐞 Include debug flag? (y/N): " DEBUG_INPUT
[[ "$DEBUG_INPUT" =~ ^[Yy]$ ]] && DEBUG_FLAG="--debug"

# --- Check for GitHub CLI ---
if ! command -v gh >/dev/null 2>&1; then
  echo "❌ GitHub CLI (gh) not found. Please install it: https://cli.github.com/"
  exit 1
fi

# --- Validate release notes ---
if ! grep -q "$VERSION" RELEASE_NOTES.md; then
  echo "❌ RELEASE_NOTES.md does not mention version $VERSION"
  exit 1
fi

echo ""
echo "📋 Summary:"
echo "  Version:     $VERSION"
echo "  Clean build: ${CLEAN_FLAG:-no}"
echo "  Debug mode:  ${DEBUG_FLAG:-no}"
echo ""

read -rp "📤 Push dev and open PR into main? (y/N): " CONFIRM
if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
  echo "❌ Aborting."
  exit 1
fi

# --- Push dev branch ---
echo "🔼 Pushing dev to origin..."
git push origin dev

# --- Check for existing PR ---
echo "🔍 Checking for existing PR from dev to main..."
EXISTING_PR_URL=$(gh pr list --base main --head dev --json url --jq '.[0].url')

if [[ -n "$EXISTING_PR_URL" ]]; then
  echo "ℹ️ An open PR from dev to main already exists:"
  echo "   $EXISTING_PR_URL"
else
  echo "🔃 Creating pull request: dev → main"
  gh pr create --base main --head dev \
    --title "Release $VERSION" \
    --body-file RELEASE_NOTES.md
fi

# --- Final instructions ---
echo ""
echo "✅ Pull request is ready. Please review and merge it on GitHub."
echo "🧠 Once merged, finish the release with:"
echo ""
echo "    git checkout main && git pull origin main"
echo "    ./release.sh $VERSION $CLEAN_FLAG $DEBUG_FLAG"
echo ""
