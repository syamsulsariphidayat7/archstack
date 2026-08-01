#!/usr/bin/env bash
set -euo pipefail

REPO="syamsulsariphidayat7/archstack"
INSTALL_DIR="/usr/local/bin"

# Deteksi arch
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  GOARCH="amd64" ;;
    aarch64|arm64) GOARCH="arm64" ;;
    *)
        echo "[error] Arch tidak didukung: $ARCH"
        exit 1
        ;;
esac

echo "[info] Mendapatkan rilis terbaru dari $REPO..."
LATEST=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" 2>/dev/null | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": "\(.*\)",/\1/')

if [ -z "$LATEST" ]; then
    echo "[error] Gagal mendapatkan rilis terbaru"
    exit 1
fi

echo "[info] Versi: $LATEST"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST/archstack-linux-$GOARCH"
echo "[info] Download dari $DOWNLOAD_URL"

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/archstack"
chmod +x "$TMP_DIR/archstack"

if [ -w "$INSTALL_DIR" ]; then
    cp "$TMP_DIR/archstack" "$INSTALL_DIR/archstack"
else
    echo "[info] $INSTALL_DIR tidak writable, pakai sudo..."
    sudo cp "$TMP_DIR/archstack" "$INSTALL_DIR/archstack"
fi

echo "[done] archstack berhasil diinstal ke $INSTALL_DIR/archstack"
echo "Jalankan: archstack"
