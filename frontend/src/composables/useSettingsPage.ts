import { ref, computed } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useUIStore } from '@/stores/ui'
import type { Theme, WeekStart } from '@/types'
import { useCurrencies } from '@/composables/useCurrencies'
import { useBiometric } from '@/composables/useBiometric'

export function useSettingsPage() {
  const settingsStore = useSettingsStore()
  const uiStore = useUIStore()
  const { status, enable, disable, loading: biometricLoading, refreshStatus } = useBiometric()
  const settings = computed(() => settingsStore.settings)

  const theme = computed<Theme>({
    get: () => settingsStore.theme,
    set: (val: Theme) => settingsStore.setSetting('theme', val)
  })

  const weekStartStr = computed<string>({
    get: () => String(settings.value.week_start ?? 1),
    set: (val: string) => settingsStore.setSetting('week_start', parseInt(val, 10) as WeekStart)
  })

  const newCurrencyCode = ref('')

  const totpButtonVariant = computed<'primary' | 'danger'>(() =>
    settingsStore.isTotpEnabled ? 'danger' : 'primary'
  )

  const { currencies, addCurrency } = useCurrencies()
  const currencyOptions = computed(() => currencies.value.map(c => ({ value: c, label: c })))

  const themeOptions = [
    { value: 'system', label: 'System' },
    { value: 'light', label: 'Light' },
    { value: 'dark', label: 'Dark' }
  ]

  const dateFormatOptions = [
    { value: 'YYYY-MM-DD', label: 'YYYY-MM-DD' },
    { value: 'DD/MM/YYYY', label: 'DD/MM/YYYY' },
    { value: 'MM/DD/YYYY', label: 'MM/DD/YYYY' },
    { value: 'YYYY年M月D日', label: 'YYYY年M月D日 (Japanese)' }
  ]

  const numberFormatOptions = [
    { value: 'ja-JP', label: 'Japanese (1,234,567)' },
    { value: 'en-US', label: 'US English (1,234,567)' },
    { value: 'en-GB', label: 'UK English (1,234,567)' },
    { value: 'de-DE', label: 'German (1.234.567)' }
  ]

  const weekStartOptions = [
    { value: '0', label: 'Sunday' },
    { value: '1', label: 'Monday' },
    { value: '6', label: 'Saturday' }
  ]

  const handleAddCurrency = async () => {
    const code = newCurrencyCode.value.trim().toUpperCase()
    if (!code) return
    const success = await addCurrency(code)
    if (success) {
      uiStore.success(`Currency ${code} added`)
      newCurrencyCode.value = ''
    } else {
      uiStore.error('Failed to add currency')
    }
  }

  const onThemeChange = async (value: string) => settingsStore.setSetting('theme', value as Theme)

  const onWeekStartChange = async (value: string) => settingsStore.setSetting('week_start', parseInt(value, 10) as WeekStart)

  const saveSettings = async () => {
    await settingsStore.updateSettings(settings.value)
    uiStore.success('Settings saved')
  }

  const toggleBiometric = async () => {
    if (status.value.enabled) { const ok = await disable(); if (ok) uiStore.success('Biometric disabled') }
    else { const ok = await enable(); if (ok) uiStore.success('Biometric enabled') }
  }

  return {
    settings, theme, weekStartStr, newCurrencyCode, totpButtonVariant, currencies,
    currencyOptions, themeOptions, dateFormatOptions, numberFormatOptions, weekStartOptions,
    handleAddCurrency, onThemeChange, onWeekStartChange, saveSettings, toggleBiometric,
    biometricLoading, status, refreshStatus
  }
}