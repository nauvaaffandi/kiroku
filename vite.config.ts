import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

export default defineConfig(({ mode }) => {
  const alias: Record<string, string> = {
    '@': path.resolve(__dirname, './src')
  }

  // capacitor-secure-storage-plugin resolves to the in-memory fallback shim only
  // in non-production modes (dev / test). Production must refuse to ship the
  // unencrypted stub; a real native plugin + `npx cap sync android` + native
  // rebuild is required before a production release.
  if (mode !== 'production') {
    alias['capacitor-secure-storage-plugin'] = path.resolve(
      __dirname,
      './src/stubs/secure-storage-plugin.ts'
    )
  }

  return {
    plugins: [vue(), tailwindcss()],
    resolve: { alias },
    build: { minify: 'esbuild' },
    server: { port: 5173, host: true }
  }
})
