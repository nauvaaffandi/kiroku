import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'
import { VitePWA } from 'vite-plugin-pwa'

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    VitePWA({
      registerType: 'autoUpdate',
      strategies: 'generateSW',
      includeAssets: ['favicon.ico', 'robots.txt'],
      manifest: {
        id: 'kiroku',
        name: 'Kiroku',
        short_name: 'Kiroku',
        description: 'Self-hosted Personal Finance Tracker',
        theme_color: '#2b67ad',
        background_color: '#f6f5f2',
        display: 'standalone',
        orientation: 'portrait-primary',
        scope: '/',
        start_url: '/',
        lang: 'ja',
        icons: [
          { src: '/icon.svg', sizes: 'any', type: 'image/svg+xml', purpose: 'any maskable' },
          { src: '/icons/icon-72x72.png', sizes: '72x72', type: 'image/png', purpose: 'any maskable' },
          { src: '/icons/icon-96x96.png', sizes: '96x96', type: 'image/png', purpose: 'any maskable' },
          { src: '/icons/icon-128x128.png', sizes: '128x128', type: 'image/png', purpose: 'any maskable' },
          { src: '/icons/icon-144x144.png', sizes: '144x144', type: 'image/png', purpose: 'any maskable' },
          { src: '/icons/icon-152x152.png', sizes: '152x152', type: 'image/png', purpose: 'any maskable' },
          { src: '/icons/icon-192x192.png', sizes: '192x192', type: 'image/png', purpose: 'any maskable' },
          { src: '/icons/icon-384x384.png', sizes: '384x384', type: 'image/png', purpose: 'any maskable' },
          { src: '/icons/icon-512x512.png', sizes: '512x512', type: 'image/png', purpose: 'any maskable' }
        ]
      },
      workbox: {
        mode: 'generateSW',
        runtimeCaching: [
          {
            urlPattern: /^https?.*\/api\/v1\//,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'api-cache',
              networkTimeoutSeconds: 10,
              expiration: { maxEntries: 100, maxAgeSeconds: 60 * 60 * 24 },
              cacheableResponse: { statuses: [0, 200] }
            }
          },
          {
            urlPattern: /^https?.*\/shared\/.*/,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'shared-cache',
              networkTimeoutSeconds: 10,
              expiration: { maxEntries: 50, maxAgeSeconds: 60 * 60 * 24 },
              cacheableResponse: { statuses: [0, 200] }
            }
          }
        ],
        cleanupOutdatedCaches: true,
        skipWaiting: true,
        clientsClaim: true,
        disableDevLogs: false
      },
      devOptions: { enabled: true }
    })
  ],
  resolve: { alias: { '@': path.resolve(__dirname, './src') } },
  build: { minify: 'esbuild' },
  server: { port: 5173, host: true }
})