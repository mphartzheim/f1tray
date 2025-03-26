#!/bin/bash

# Retrieve the version from git (e.g., v1.2.3 or commit hash)
VERSION=$(git describe --tags --always)
APP_NAME="f1tray-${VERSION}"
TARGETS="windows/amd64,darwin/universal"

# Delete previous build artifacts if they exist
rm -f f1tray-*.exe
rm -f f1tray-*.zip

# Ensure tray icon exists for go:embed
if [ ! -f cmd/f1tray/assets/tray_icon.png ]; then
  echo "❌ Error: cmd/f1tray/assets/tray_icon.png is missing — required for go:embed"
  exit 1
fi

echo "🚀 Building $APP_NAME for: $TARGETS"
# The -ldflags flag embeds the version information into your binary (assuming main.version is defined)
xgo --targets=$TARGETS --out $APP_NAME --pkg cmd/f1tray -ldflags "-X main.version=${VERSION}" .

echo "✅ Build complete:"
ls -1 ${APP_NAME}-*

# Automatically zip the Windows executable
WINDOWS_EXE=$(ls ${APP_NAME}-windows-*.exe 2>/dev/null | head -n 1)
if [ -f "$WINDOWS_EXE" ]; then
  ZIP_NAME="${APP_NAME}-windows.zip"
  zip -j "$ZIP_NAME" "$WINDOWS_EXE"
  echo "✅ Zipped $WINDOWS_EXE into $ZIP_NAME"
else
  echo "❌ Windows executable not found, skipping zip."
fi
