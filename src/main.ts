import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './assets/main.css'
import { useI18nStore } from '@/locales'
import { cleanup, manifestManager, seed } from '@/repository'
import { checkAndProcessNewMonth } from '@/repository'

async function init() {
  // Initialize filesystem repository layer on app startup.
  // Wrapped in try/catch so a missing native plugin (e.g. browser/dev) does not
  // prevent the app from mounting.
  try {
    await seed()
    await cleanup.runStartup()
    await manifestManager.load()
    // Auto-create transaction day files for a new month (year taken from current date).
    await checkAndProcessNewMonth()
  } catch (e) {
    console.error('[Repository] startup init failed:', e)
  }

  const app = createApp(App)
  const pinia = createPinia()
  app.use(pinia)
  app.use(router)

  // Initialize i18n store (sets document lang from localStorage)
  useI18nStore()

  app.mount('#app')
}

init().catch((e) => {
  console.error('[Repository] init error:', e)
})
