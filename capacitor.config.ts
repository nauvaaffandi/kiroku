import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'app.kiroku.finance',
  appName: 'Kiroku',
  webDir: 'dist',
  plugins: {
    SecureStoragePlugin: {}
  },
  server: {
    url: 'http://localhost:5173',
    cleartext: true,
  },
};

export default config;
