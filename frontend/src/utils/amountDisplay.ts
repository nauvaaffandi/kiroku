import { useSettingsStore } from '@/stores/settings'

const CURRENCY_SYMBOLS: Record<string, string> = {
  JPY: '¥',
  USD: '$',
  EUR: '€',
  GBP: '£',
  CNY: '¥',
  KRW: '₩',
  INR: '₹',
  SGD: 'S$',
  HKD: 'HK$',
  AUD: 'A$',
  CAD: 'C$',
  CHF: 'CHF',
  THB: '฿',
  VND: '₫',
  IDR: 'Rp',
  MYR: 'RM',
  PHP: '₱',
  TWD: 'NT$',
  BRL: 'R$',
  MXN: 'MX$',
  RUB: '₽',
  TRY: '₺',
  PLN: 'zł',
  CZK: 'Kč',
  HUF: 'Ft',
  ILS: '₪',
  ZAR: 'R',
  SEK: 'kr',
  NOK: 'kr',
  DKK: 'kr',
  NZD: 'NZ$'
}

export function getCurrencySymbol(currency: string): string {
  return CURRENCY_SYMBOLS[currency] || currency
}

export function formatNumber(value: number, locale = 'ja-JP'): string {
  return new Intl.NumberFormat(locale).format(Math.abs(value))
}

export function formatAmountForDisplay(
  amount: number,
  type: 'income' | 'expense' | 'transfer',
  signOverride?: '+' | '-'
): string {
  const settingsStore = useSettingsStore()

  if (settingsStore.isHideAmounts) {
    const masked = '*****'
    const sign = signOverride || (type === 'income' ? '+' : type === 'expense' ? '-' : '')
    const symbol = getCurrencySymbol(settingsStore.baseCurrency)
    return `${sign}${symbol}${masked}`
  }

  const sign = signOverride || (type === 'income' ? '+' : type === 'expense' ? '-' : '')
  const symbol = getCurrencySymbol(settingsStore.baseCurrency)
  const formatted = formatNumber(amount)
  return `${sign}${symbol}${formatted}`
}