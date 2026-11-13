import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './assets/main.css'
import { useI18nStore } from '@/locales'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)

// Initialize i18n store (sets document lang from localStorage)
useI18nStore()

app.mount('#app')