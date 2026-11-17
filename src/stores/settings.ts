import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import JSZip from 'jszip'
import { localDb, type EntityKind } from '@/storage/repository'
import { useSecureStorage, SECURE_KEYS } from '@/composables/useSecureStorage'
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
  hide_amounts: false,
  ocr_enabled: false,
  ocr_api_key: '',
  ai_provider_enabled: false,
  ai_provider_url: '',
  ai_model: '',
  ai_api_key: ''
}

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<Settings>({ ...defaultSettings })
  const loading = ref(false)
  const secure = useSecureStorage()

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
  const isOcrEnabled = computed(() => settings.value.ocr_enabled === true)
  const isAiProviderEnabled = computed(() => settings.value.ai_provider_enabled === true)

  const fetchSettings = async () => {
    console.log('[Settings] fetchSettings')
    loading.value = true
    try {
      const stored = await localDb.get<Settings>('setting', SETTINGS_ID)
      console.log('[Settings] settings loaded', stored ? 'found' : 'empty')
      if (stored) {
        const { migrated, redactedKeys } = await secure.migratePlaintextKeys(stored)
        if (migrated) {
          const payload = { ...stored }
          for (const key of redactedKeys) {
            payload[key as keyof Settings] = ''
          }
          // Keys whose secure write failed (returned by migratePlaintextKeys as failedKeys)
          // are intentionally left in SQLite plaintext so they are not lost.
          await localDb.put('setting', SETTINGS_ID, payload)
          console.log('[Settings] redacted plaintext API keys from SQLite')
        }
        settings.value = { ...defaultSettings, ...stored } as Settings
      } else {
        settings.value = { ...defaultSettings }
      }
      // Secure storage wins over (possibly redacted/empty) SQLite values.
      const secureOcr = await secure.getSecure('ocr_api_key')
      const secureAi = await secure.getSecure('ai_api_key')
      if (secureOcr) settings.value.ocr_api_key = secureOcr
      if (secureAi) settings.value.ai_api_key = secureAi
      applyTheme(settings.value.theme || 'system')
    } catch (e) {
      console.error('Failed to load settings:', e)
    } finally {
      loading.value = false
    }
  }

  const updateSettings = async (newSettings: Partial<Settings>) => {
    console.log('[Settings] updateSettings')
    try {
      const next = { ...settings.value, ...newSettings } as Settings
      // API keys live in secure storage only; redact them from the SQLite payload.
      const payload = { ...next }
      for (const key of SECURE_KEYS) {
        if (key in newSettings) {
          const value = next[key] as string
          if (!value) continue
          const ok = await secure.setSecure(key as string, value)
          if (ok) {
            payload[key] = ''
          } else {
            console.log('[Settings] setSecure failed')
            payload[key] = value
          }
        }
      }
      settings.value = next
      await localDb.put('setting', SETTINGS_ID, payload)
      console.log('[Settings] settings updated')
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
    isOcrEnabled,
    isAiProviderEnabled,
    fetchSettings,
    updateSettings,
    getSetting,
    setSetting,
    createBackup,
    clearAllData
  }
})