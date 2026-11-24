import type { CapacitorConfig } from '@capacitor/cli';

/// <reference types="@capacitor/background-runner" />

const config: CapacitorConfig = {
  appId: 'app.kiroku.finance',
  appName: 'Kiroku',
  webDir: 'dist',
  plugins: {
    SecureStoragePlugin: {},
    BackgroundRunner: {
      label: 'app.kiroku.finance.monthly.dirs',
      src: 'background-runner.js',
      event: 'kiroku_monthly_dirs',
      repeat: true,
      interval: 60,
      autoStart: true,
    },
  },
  server: {
    url: 'http://localhost:5173',
    cleartext: true,
  },
};

export default config;
