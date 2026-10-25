import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api/client'
import type { Settings, Theme, Language, DateFormat, NumberFormat, WeekStart } from '@/types'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<Settings>({
    theme: 'system',
    language: 'ja',
    base_currency: 'USD',
    date_format: 'YYYY-MM-DD',
    number_format: 'ja-JP',
    week_start: 1,
    totp_enabled: false,
    hide_amounts: false
  })
  const loading = ref(false)

  const theme = computed<Theme>(() => settings.value.theme || 'system')
  const language = computed<Language>(() => settings.value.language || 'ja')
  const baseCurrency = computed(() => settings.value.base_currency || 'USD')
  const dateFormat = computed<DateFormat>(() => settings.value.date_format || 'YYYY-MM-DD')
  const numberFormat = computed<NumberFormat>(() => settings.value.number_format || 'ja-JP')
  const weekStart = computed<WeekStart>(() => {
    const val = parseInt(String(settings.value.week_start || '1'), 10)
    return (val as WeekStart)
  })
  const isTotpEnabled = computed(() => settings.value.totp_enabled === true)
  const isHideAmounts = computed(() => settings.value.hide_amounts === true)

  const fetchSettings = async () => {
    loading.value = true
    try {
      const response = await api.get('/api/v1/settings')
      settings.value = { ...settings.value, ...response.data } as Settings
      applyTheme(settings.value.theme || 'system')
    } catch (e) {
      console.error('Failed to fetch settings:', e)
    } finally {
      loading.value = false
    }
  }

  const updateSettings = async (newSettings: Partial<Settings>) => {
    try {
      await api.put('/api/v1/settings', newSettings)
      settings.value = { ...settings.value, ...newSettings } as Settings
      if (newSettings.theme) applyTheme(newSettings.theme)
      return true
    } catch (e) {
      console.error('Failed to update settings:', e)
      return false
    }
  }

  const getSetting = <K extends keyof Settings>(key: K): Settings[K] => settings.value[key]
  const setSetting = async <K extends keyof Settings>(key: K, value: Settings[K]) => {
    await updateSettings({ [key]: value } as Partial<Settings>)
  }

  const createBackup = async () => {
    const response = await api.post('/api/v1/backup/create')
    return response.data
  }

  const restoreBackup = async (file: File, password: string) => {
    const formData = new FormData()
    formData.append('backup', file)
    formData.append('password', password)
    await api.post('/api/v1/backup/restore', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }

  const clearAllData = async () => {
    await api.post('/api/v1/settings/clear')
  }

  const applyTheme = (themeValue: string) => {
    const root = document.documentElement
    if (themeValue === 'system') {
      root.removeAttribute('data-theme')
    } else {
      root.setAttribute('data-theme', themeValue)
    }
  }

  return {
    settings,
    loading,
    theme,
    language,
    baseCurrency,
    dateFormat,
    numberFormat,
    weekStart,
    isTotpEnabled,
    isHideAmounts,
    fetchSettings,
    updateSettings,
    getSetting,
    setSetting,
    createBackup,
    restoreBackup,
    clearAllData
  }
})