import { computed } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { getCurrencySymbol } from '@/utils/amountDisplay'

export const useDashboardFormatting = () => {
  const settingsStore = useSettingsStore()
  const isHideAmounts = computed(() => settingsStore.isHideAmounts)
  const baseCurrency = computed(() => settingsStore.baseCurrency)
  const numberFormat = computed(() => settingsStore.numberFormat)
  const currencySymbol = computed(() => getCurrencySymbol(baseCurrency.value))

  const formatNumber = (value: number) => new Intl.NumberFormat(numberFormat.value).format(Math.abs(value))
  const maskAmount = (sign: string) => `${sign}${currencySymbol.value}*****`

  const formatIncome = (value: number): string => {
    if (isHideAmounts.value) return maskAmount('+')
    return `+${currencySymbol.value}${formatNumber(value)}`
  }

  const formatExpense = (value: number): string => {
    if (isHideAmounts.value) return maskAmount('−')
    return `−${currencySymbol.value}${formatNumber(value)}`
  }

  const formatNet = (value: number): string => {
    if (isHideAmounts.value) {
      const sign = value >= 0 ? '+' : '−'
      return maskAmount(sign)
    }
    const sign = value >= 0 ? '+' : '−'
    return `${sign}${currencySymbol.value}${formatNumber(value)}`
  }

  const formatBalance = (value: number): string => {
    if (isHideAmounts.value) {
      const sign = value < 0 ? '−' : ''
      return `${sign}${currencySymbol.value}*****`
    }
    const sign = value < 0 ? '−' : '+'
    return `${sign}${currencySymbol.value}${formatNumber(value)}`
  }

  return { formatIncome, formatExpense, formatNet, formatBalance }
}
