import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import JSZip from 'jszip'
import { localDb, type EntityKind } from '@/storage/repository'
import type { Settings, Theme, Language, DateFormat, NumberFormat, WeekStart } from '@/types'

const BACKUP_KINDS: readonly string[] = [
  'account', 'category', 'label', 'notification', 'transaction', 'template', 'scheduled', 'budget', 'goal', 'setting'
]

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

  const createBackup = async (): Promise<Blob> => {
    console.log('[Backup] Starting backup creation')
    const zip = new JSZip()
    const payload: Record<string, unknown[]> = {}

    for (const kind of BACKUP_KINDS) {
      console.log(`[Backup] Fetching data for kind: ${kind}`)
      const items = await localDb.list(kind as EntityKind)
      console.log(`[Backup] Retrieved ${items.length} records for kind: ${kind}`)
      payload[kind] = items
      zip.file(`${kind}.json`, JSON.stringify(items, null, 2))
    }

    const manifest = {
      version: '1.0',
      exported_at: new Date().toISOString(),
      counts: Object.fromEntries(Object.entries(payload).map(([k, v]) => [k, v.length]))
    }
    console.log('[Backup] Manifest generated:', manifest)
    zip.file('manifest.json', JSON.stringify(manifest, null, 2))

    const blob = await zip.generateAsync({ type: 'blob' })
    console.log(`[Backup] ZIP created successfully, size: ${blob.size} bytes`)
    return blob
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
    clearAllData
  }
})