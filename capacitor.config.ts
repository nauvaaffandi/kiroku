import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'app.kiroku.finance',
  appName: 'Kiroku',
  webDir: 'dist',
  // Live Reload untuk development di Termux.
  // Ganti IP di bawah dengan alamat LAN Termux (jalankan `ip addr` di Termux untuk mencarinya),
  // contoh: http://192.168.1.100:5173
  // Jangan pakai 127.0.0.1 — WebView Android tidak dapat mencapai localhost Termux.
  // Untuk build produksi via CI/CD, komentar bagian server ini terlebih dahulu agar APK
  // memuat aset web lokal (dist/) alih-alih mencoba menghubungi dev server.
  server: {
    url: 'http://192.168.0.160:5173',
    cleartext: true,
  },
};

export default config;
