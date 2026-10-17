import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './assets/main.css'
import { registerSW } from 'virtual:pwa-register'
import { useI18nStore } from '@/locales'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)

// Initialize i18n store (sets document lang from localStorage)
useI18nStore()

// Register service worker
if ('serviceWorker' in navigator) {
  registerSW({
    onNeedRefresh() {
      if (confirm('New content available. Reload?')) {
        window.location.reload()
      }
    },
    onOfflineReady() {
      console.log('App ready to work offline')
    }
  })
}

app.mount('#app')