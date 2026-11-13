import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { localDb } from '@/storage/repository'
import type { Settings, Theme, Language, DateFormat, NumberFormat, WeekStart } from '@/types'

const SETTINGS_ID = 'settings'

const defaultSettings: Settings = {
  theme: 'system',
  language: 'ja',
  base_currency: 'USD',
  date_format: 'YYYY-MM-DD',
  number_format: 'ja-JP',
  week_start: 1,
  totp_enabled: false,
  hide_amounts: false
}

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<Settings>({ ...defaultSettings })
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
      const stored = await localDb.get<Settings>('setting', SETTINGS_ID)
      if (stored) settings.value = { ...defaultSettings, ...stored }
      applyTheme(settings.value.theme || 'system')
    } catch (e) {
      console.error('Failed to load settings:', e)
    } finally {
      loading.value = false
    }
  }

  const updateSettings = async (newSettings: Partial<Settings>) => {
    try {
      settings.value = { ...settings.value, ...newSettings } as Settings
      await localDb.put('setting', SETTINGS_ID, settings.value)
      if (newSettings.theme) applyTheme(newSettings.theme)
      return true
    } catch (e) {
      console.error('Failed to save settings:', e)
      return false
    }
  }

  const getSetting = <K extends keyof Settings>(key: K): Settings[K] => settings.value[key]
  const setSetting = async <K extends keyof Settings>(key: K, value: Settings[K]) => {
    await updateSettings({ [key]: value } as Partial<Settings>)
  }

  const createBackup = async () => {
    throw new Error('Backup is not available offline yet')
  }

  const restoreBackup = async (_file: File, _password: string) => {
    throw new Error('Restore is not available offline yet')
  }

  const clearAllData = async () => {
    throw new Error('Clear all data is not available offline yet')
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