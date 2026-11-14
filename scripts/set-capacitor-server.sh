#!/usr/bin/env bash
# set-capacitor-server.sh
# Menambah/mengganti konfigurasi server Capacitor (Live Reload) di capacitor.config.ts
# Penggunaan: ./scripts/set-capacitor-server.sh <host>:<port>
# Contoh:   ./scripts/set-capacitor-server.sh 192.168.0.160:5173

set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Penggunaan: $0 <host>:<port>"
  echo "Contoh:   $0 192.168.0.160:5173"
  exit 1
fi

TARGET="$1"
HOST="${TARGET%%:*}"
PORT="${TARGET##*:}"

if [ "$HOST" = "$TARGET" ] || [ -z "$PORT" ]; then
  echo "Error: format harus <host>:<port>, contoh: 192.168.0.160:5173"
  exit 1
fi

CONFIG="capacitor.config.ts"

if [ ! -f "$CONFIG" ]; then
  echo "Error: $CONFIG tidak ditemukan (jalankan dari root repo)"
  exit 1
fi

# Hapus blok server yang sudah ada (jika ada)
sed -i '/  server: {/,/  },/d' "$CONFIG"

# Tambahkan koma pada baris webDir jika belum ada (supaya valid TS)
if grep -q "webDir: 'dist'$" "$CONFIG"; then
  sed -i "s/webDir: 'dist'$/webDir: 'dist',/" "$CONFIG"
fi

# Sisipkan blok server sebelum baris penutupan objek config (});
sed -i "/^};/i\  server: {\n    url: 'http://${HOST}:${PORT}',\n    cleartext: true,\n  }," "$CONFIG"

echo "OK: server.url di-set ke http://${HOST}:${PORT} di $CONFIG"
echo " Jalankan 'npx cap sync android' setelah ini untuk menyinkronkan ke Android."
