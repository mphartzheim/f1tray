#!/bin/bash

APP_NAME="f1tray"
TARGETS="windows/amd64,darwin/universal"

echo "🚀 Building $APP_NAME for: $TARGETS"
xgo --targets=$TARGETS --out $APP_NAME --pkg cmd/f1tray .

echo "✅ Build complete:"
ls -1 ${APP_NAME}-*
