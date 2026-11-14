#!/usr/bin/env bash
# reset-capacitor-server.sh
# Menghapus konfigurasi server Capacitor (Live Reload) di capacitor.config.ts
# Mengembalikannya ke mode produksi/lokal (tanpa URL server eksternal)
# Penggunaan: ./scripts/reset-capacitor-server.sh

set -euo pipefail

CONFIG="capacitor.config.ts"

if [ ! -f "$CONFIG" ]; then
  echo "Error: $CONFIG tidak ditemukan (jalankan dari root repo)"
  exit 1
fi

# Hapus blok server
sed -i '/  server: {/,/  },/d' "$CONFIG"

# Hapus koma berlebih di webDir jika ada
sed -i "s/webDir: 'dist',$/webDir: 'dist'/" "$CONFIG"

echo "OK: server config dihapus dari $CONFIG (siap untuk build CI/CD)"
