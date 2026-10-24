#!/usr/bin/env bash
set -euo pipefail

BINARY_NAME="${BINARY_NAME:-attention-span}"
BUILD_DIR=bin
VERSION="${VERSION:-1.0.0}"
INSTALL_DIR="${INSTALL_DIR:-${HOME}/.local/bin}"

build_linux_amd64() {
	mkdir -p "$BUILD_DIR"
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$VERSION" -o "$BUILD_DIR/$BINARY_NAME" ./cmd/attention-span
	echo "✓ Linux binary built: $BUILD_DIR/$BINARY_NAME"
}

build_darwin_arm64() {
  mkdir -p "$BUILD_DIR"
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=$VERSION" -o "$BUILD_DIR/$BINARY_NAME" ./cmd/attention-span
	echo "✓ Darwin binary built: $BUILD_DIR/$BINARY_NAME"
}

main() {
  # build at BUILD_DIR 
  if [ "$(uname -s)" = "Darwin" ]; then
    build_darwin_arm64
  elif [ "$(uname -s)" = "Linux" ]; then
    build_linux_amd64
  fi

  # install to INSTALL_DIR
  if command -v install >/dev/null 2>&1; then
    install -m 755 "${BUILD_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
  else
    cp "${BUILD_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod 755 "${INSTALL_DIR}/${BINARY_NAME}"
  fi
  # Darwin (macOS) specific quarantine handling
  if [ "$(uname -s)" = "Darwin" ] && command -v xattr >/dev/null 2>&1; then
    xattr -d com.apple.quarantine "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null || true
  fi
  rm -rf "$BUILD_DIR"
  echo "Installed: $INSTALL_DIR/$BINARY_NAME"
  
  case ":${PATH}:" in
  *:"${INSTALL_DIR}":*)
    # Destination is already in PATH
    ;;
  *)
    echo ""
    echo "$INSTALL_DIR is not on your PATH. Add to ~/.zshrc or ~/.bashrc:"
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    ;;
  esac
}

main