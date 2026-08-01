#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="/usr/local/bin"
BINARY="$INSTALL_DIR/archstack"

if [ ! -f "$BINARY" ]; then
    echo "[info] archstack tidak ditemukan di $BINARY"
    exit 0
fi

echo "Menghapus $BINARY..."

if [ -w "$INSTALL_DIR" ]; then
    rm "$BINARY"
else
    sudo rm "$BINARY"
fi

echo "[done] archstack berhasil dihapus dari $BINARY"
