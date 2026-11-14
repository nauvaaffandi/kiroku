import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'app.kiroku.finance',
  appName: 'Kiroku',
  webDir: 'dist',
  server: {
    url: 'http://192.168.0.160:5173',
    cleartext: true,
  },
};

export default config;
