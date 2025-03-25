#!/bin/bash

APP_NAME="F1Tray"
BINARY_NAME="f1tray"
BINARY_PATH="./cmd/$BINARY_NAME/$BINARY_NAME"
APPDIR="${APP_NAME}.AppDir"

# Grab version from Git (fallback to 'dev' if no tag)
GIT_VERSION=$(git describe --tags --dirty 2>/dev/null || echo "dev")
OUTPUT="${APP_NAME}-${GIT_VERSION}-x86_64.AppImage"

echo "🚀 Building Linux AppImage for $APP_NAME v$GIT_VERSION"

# Check for required tools
for cmd in go magick appimagetool git; do
  if ! command -v $cmd &>/dev/null; then
    echo "❌ Error: '$cmd' not found. Please install it before continuing."
    exit 1
  fi
done

# Handle optional --clean flag
if [[ "$1" == "--clean" ]]; then
  echo "🧹 Cleaning previous build artifacts..."
  rm -rf "$APPDIR" "$OUTPUT"
fi

# Compile the Go binary
echo "🛠 Compiling Linux binary..."
go build -o "$BINARY_PATH" ./cmd/$BINARY_NAME
if [ $? -ne 0 ]; then
  echo "❌ Go build failed. Aborting."
  exit 1
fi
chmod +x "$BINARY_PATH"

# Create AppDir structure
mkdir -p "$APPDIR/usr/bin"

# Copy binary into AppDir
cp "$BINARY_PATH" "$APPDIR/usr/bin/$APP_NAME"

# Create AppRun script
cat > "$APPDIR/AppRun" <<EOF
#!/bin/sh
exec usr/bin/$APP_NAME "\$@"
EOF
chmod +x "$APPDIR/AppRun"

# Create .desktop file
cat > "$APPDIR/$APP_NAME.desktop" <<EOF
[Desktop Entry]
Name=$APP_NAME
Exec=$APP_NAME
Icon=$BINARY_NAME
Type=Application
Categories=Utility;
EOF

# Generate icon with ImageMagick (using 'magick')
echo "🎨 Converting icon..."
magick assets/tray_icon.png -resize 256x256 "$APPDIR/$BINARY_NAME.png"
chmod 644 "$APPDIR/$BINARY_NAME.png"

# Package AppImage
echo "📦 Creating AppImage: $OUTPUT"
ARCH=x86_64 appimagetool "$APPDIR" "$OUTPUT"

if [ $? -eq 0 ]; then
  echo "✅ Done! AppImage created: $OUTPUT"
else
  echo "❌ AppImage creation failed."
  exit 1
fi
