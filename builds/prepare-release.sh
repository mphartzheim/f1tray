#!/bin/bash
set -e

# Check for GitHub CLI (gh)
if ! command -v gh &> /dev/null; then
  echo "GitHub CLI (gh) is required. Please install it from https://cli.github.com/"
  exit 1
fi

# Define your version tag
VERSION_TAG=$(git describe --tags --abbrev=0)

# Create a pull request from dev into main using GitHub CLI.
echo "Creating a pull request from dev into main for release ${VERSION_TAG}..."
gh pr create --base main --head dev --title "Release ${VERSION_TAG}" --body "Automated PR for release ${VERSION_TAG}. Please review and merge."
echo "PR created successfully. Please review and merge the PR on GitHub to continue with the release."
