import { computed } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useUIStore } from '@/stores/ui'

/**
 * Thin theme composable. All theme reads/writes go through
 * useSettingsStore / useUIStore (single source of truth: data-theme only).
 * Components never manipulate data-theme or the .dark class directly.
 */
export function useTheme() {
  const settingsStore = useSettingsStore()
  const uiStore = useUIStore()

  const theme = computed(() => settingsStore.theme)

  const isDark = computed(() => {
    const t = theme.value
    if (t === 'system') {
      return window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
    }
    return t === 'dark'
  })

  const setTheme = (value: 'light' | 'dark' | 'system') => {
    console.log('[useTheme] setTheme', { value })
    uiStore.setTheme(value)
  }

  return { theme, isDark, setTheme }
}
