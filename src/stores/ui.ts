import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Toast, ToastType } from '@/types'

export const useUIStore = defineStore('ui', () => {
  const sidebarOpen = ref(false)
  const toasts = ref<Toast[]>([])
  const theme = ref<'light' | 'dark' | 'system'>(
    (localStorage.getItem('theme') as 'light' | 'dark' | 'system') || 'system'
  )
  let toastId = 0

  const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
  }

  const setTheme = (newTheme: 'light' | 'dark' | 'system') => {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    applyTheme(newTheme)
  }

  const applyTheme = (currentTheme: 'light' | 'dark' | 'system') => {
    const root = document.documentElement
    if (currentTheme === 'system') {
      root.removeAttribute('data-theme')
    } else {
      root.setAttribute('data-theme', currentTheme)
    }
  }

  const showToast = (type: ToastType, message: string, duration = 4000) => {
    const id = ++toastId
    toasts.value.push({ id, type, message })
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, duration)
  }

  const success = (message: string) => showToast('success', message)
  const error = (message: string) => showToast('error', message)
  const warning = (message: string) => showToast('warning', message)
  const info = (message: string) => showToast('info', message)

  // Initialize theme on store creation
  applyTheme(theme.value)

  return {
    sidebarOpen,
    toasts,
    theme,
    toggleSidebar,
    setTheme,
    showToast,
    success,
    error,
    warning,
    info
  }
})