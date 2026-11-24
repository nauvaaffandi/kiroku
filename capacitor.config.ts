import type { CapacitorConfig } from '@capacitor/cli';

/// <reference types="@capacitor/preferences" />
/// <reference types="@capacitor/local-notifications" />
/// <reference types="@capacitor/toast" />
/// <reference types="@capacitor/clipboard" />
/// <reference types="@capacitor/camera" />
/// <reference types="@capacitor/geolocation" />

const config: CapacitorConfig = {
  appId: 'app.kiroku.finance',
  appName: 'Kiroku',
  webDir: 'dist',
  plugins: {
    SecureStoragePlugin: {},
    Camera: {},
    Clipboard: {},
    Geolocation: {},
    LocalNotifications: {},
    Preferences: {},
    Toast: {},
  },
  server: {
    url: 'http://localhost:5173',
    cleartext: true,
  },
};

export default config;
