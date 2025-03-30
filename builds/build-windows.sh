#!/bin/bash

# Check for docker group access
if ! groups | grep -qw docker; then
  echo "❌ Error: You are not in the docker group. Please add your user to the docker group and log out/in."
  exit 1
fi

# Check for required tools
for cmd in xgo; do
  if ! command -v $cmd &>/dev/null; then
    echo "❌ Error: '$cmd' not found. Please install it before continuing."
    exit 1
  fi
done

# Directories and version info
OUTPUT_DIR="builds/output/windows"
RESOURCE_SYSO_SOURCE="builds/windows/resource_windows.syso"
RESOURCE_SYSO_TARGET="cmd/f1tray/resource_windows.syso"
mkdir -p "$OUTPUT_DIR"

VERSION=$(git describe --tags --always)
APP_NAME="f1tray"
APP_NAME_WITH_VERSION="${APP_NAME}-${VERSION}"
TARGETS="windows/amd64,darwin/universal"

EXE_PATH="${OUTPUT_DIR}/${APP_NAME_WITH_VERSION}"

# Clean old build artifacts
for arg in "$@"; do
  if [[ "$arg" == "--clean" ]]; then
    echo "🧹 Cleaning previous build artifacts..."
    rm -f ${OUTPUT_DIR}/${APP_NAME}-*.exe
    rm -f ${OUTPUT_DIR}/${APP_NAME}-*.zip
  fi
done

# Ensure tray icon exists for go:embed
if [ ! -f cmd/f1tray/assets/tray_icon.png ]; then
  echo "❌ Error: cmd/f1tray/assets/tray_icon.png is missing — required for go:embed"
  exit 1
fi

# Handle Windows resource file
if [ -f "$RESOURCE_SYSO_SOURCE" ]; then
  echo "📦 Copying Windows resource file..."
  cp "$RESOURCE_SYSO_SOURCE" "$RESOURCE_SYSO_TARGET"
else
  echo "⚠️  Warning: builds/windows/resource_windows.syso not found. Continuing without Windows metadata."
fi

echo "🚀 Building $APP_NAME_WITH_VERSION for: $TARGETS"
xgo --targets=$TARGETS \
    --out "$EXE_PATH" \
    --pkg cmd/f1tray \
    -ldflags="-H=windowsgui -X main.version=${VERSION}" .

# Remove resource file to keep repo clean
rm -f "$RESOURCE_SYSO_TARGET"

echo "✅ Build complete:"
ls -1 ${OUTPUT_DIR}/${APP_NAME}-*

# Zip the Windows executable
WINDOWS_EXE=$(ls ${OUTPUT_DIR}/${APP_NAME_WITH_VERSION}-windows-*.exe 2>/dev/null | head -n 1)
if [ -f "$WINDOWS_EXE" ]; then
  ZIP_NAME="${OUTPUT_DIR}/${APP_NAME_WITH_VERSION}-windows.zip"
  zip -j "$ZIP_NAME" "$WINDOWS_EXE"
  echo "✅ Zipped $WINDOWS_EXE into $ZIP_NAME"
else
  echo "❌ Windows executable not found, skipping zip."
fi
