import { describe, test, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useSettingsStore } from '@/stores/settings'
import { formatNumber, formatAmountForDisplay, getCurrencySymbol } from '@/utils/amountDisplay'

describe('amountDisplay utility', () => {
  let settingsStore: ReturnType<typeof useSettingsStore>

  beforeEach(() => {
    setActivePinia(createPinia())
    settingsStore = useSettingsStore()
    settingsStore.settings.number_format = 'en-US'
    settingsStore.settings.base_currency = 'USD'
    settingsStore.settings.hide_amounts = false
  })

  describe('getCurrencySymbol', () => {
    test('returns correct symbol for known currency', () => {
      expect(getCurrencySymbol('USD')).toBe('$')
      expect(getCurrencySymbol('JPY')).toBe('¥')
      expect(getCurrencySymbol('EUR')).toBe('€')
      expect(getCurrencySymbol('IDR')).toBe('Rp')
      expect(getCurrencySymbol('GBP')).toBe('£')
    })

    test('returns currency code as fallback for unknown currency', () => {
      expect(getCurrencySymbol('XYZ')).toBe('XYZ')
    })
  })

  describe('formatNumber', () => {
    test('formats positive numbers using settings number format', () => {
      settingsStore.settings.number_format = 'en-US'
      expect(formatNumber(1234567.89)).toBe('1,234,567.89')
    })

    test('formats numbers with custom locale override', () => {
      expect(formatNumber(1234.5, 'de-DE')).toBe('1.234,5')
    })

    test('formats absolute value for negative numbers', () => {
      expect(formatNumber(-500)).toBe('500')
    })

    test('masks value with •••• when hide_amounts is enabled', () => {
      settingsStore.settings.hide_amounts = true
      expect(formatNumber(1000)).toBe('••••')
      expect(formatNumber(0)).toBe('••••')
      expect(formatNumber(-250)).toBe('••••')
    })
  })

  describe('formatAmountForDisplay', () => {
    test('formats income with + sign and currency symbol', () => {
      settingsStore.settings.base_currency = 'USD'
      expect(formatAmountForDisplay(1500, 'income')).toBe('+$1,500')
    })

    test('formats expense with - sign and currency symbol', () => {
      settingsStore.settings.base_currency = 'USD'
      expect(formatAmountForDisplay(75, 'expense')).toBe('−$75')
    })

    test('formats transfer with no sign prefix', () => {
      settingsStore.settings.base_currency = 'USD'
      expect(formatAmountForDisplay(200, 'transfer')).toBe('$200')
    })

    test('respects explicit signOverride', () => {
      expect(formatAmountForDisplay(100, 'expense', '+')).toBe('+$100')
      expect(formatAmountForDisplay(100, 'income', '-')).toBe('-$100')
    })

    test('masks amount with •••• when hide_amounts is active', () => {
      settingsStore.settings.hide_amounts = true
      expect(formatAmountForDisplay(5000, 'income')).toBe('+$••••')
      expect(formatAmountForDisplay(300, 'expense')).toBe('−$••••')
      expect(formatAmountForDisplay(100, 'transfer')).toBe('$••••')
    })

    test('correctly uses JPY currency symbol and integer format', () => {
      settingsStore.settings.base_currency = 'JPY'
      settingsStore.settings.number_format = 'ja-JP'
      expect(formatAmountForDisplay(10000, 'income')).toBe('+¥10,000')
    })

    test('correctly uses IDR currency symbol', () => {
      settingsStore.settings.base_currency = 'IDR'
      settingsStore.settings.number_format = 'id-ID'
      expect(formatAmountForDisplay(50000, 'expense')).toBe('−Rp50.000')
    })
  })
})
