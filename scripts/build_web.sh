#!/usr/bin/env bash

# This bash script does the following:
# Complie our Go game into WebAssemply, inside game.wasm file.
# Copy the correct wasm_exec.js from the Go installation.
# Copy assets/ into web/.
# Produce a ready-to-publish web/ folder.
# It fails and exits if one of the following directories/files doesn't exist:
# assets/
# wasm_exec.js
# web/index.html

set -euo pipefail


REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
WEB_DIR="$REPO_ROOT/web"
ASSETS_SRC="$REPO_ROOT/assets"
WASM_OUT="$WEB_DIR/game.wasm"

# Ensure web dir exists
mkdir -p "$WEB_DIR"

# Clean previous build artifacts if they exist (keeps web/ deterministic)
echo "[build_web] Deleting previous build output if it exists ..."
rm -f "$WASM_OUT"
rm -rf "$WEB_DIR/assets"

# Build WASM
echo "[build_web] Compiling Go program to WebAssembly ..."
# Build just the package that contains main (assume module root or provide path)
GOOS=js GOARCH=wasm go build -o "$WASM_OUT" ./...

# Copy the correct wasm_exec.js depending on Go layout (supports old/new layouts)
# Newer Go: $(go env GOROOT)/lib/wasm/wasm_exec.js
# Older Go: $(go env GOROOT)/misc/wasm/wasm_exec.js
GOROOT="$(go env GOROOT)"
if [ -f "$GOROOT/lib/wasm/wasm_exec.js" ]; then
  cp "$GOROOT/lib/wasm/wasm_exec.js" "$WEB_DIR/wasm_exec.js"
elif [ -f "$GOROOT/misc/wasm/wasm_exec.js" ]; then
  cp "$GOROOT/misc/wasm/wasm_exec.js" "$WEB_DIR/wasm_exec.js"
else
  echo "ERROR: wasm_exec.js not found in GOROOT ($GOROOT)."
  exit 2
fi

# Copy assets into web/ (fresh copy)
if [ -d "$ASSETS_SRC" ]; then
  echo "[build_web] Copying assets/ inside web/ ..."
  cp -r "$ASSETS_SRC" "$WEB_DIR/"
else
  echo "[build_web] Error: Source assets dir not found at $ASSETS_SRC "
  echo "Build failed."
  exit 3
fi

# Ensure index.html exists
if [ ! -f "$WEB_DIR/index.html" ]; then
  echo "[build_web] Error: web/index.html is missing. Please create it."
  echo "Build failed."
  exit 4
fi

echo "[build_web] Done, web/ is ready for publish."