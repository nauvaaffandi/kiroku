import { describe, test, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useSettingsStore } from '@/stores/settings'
import { formatNumber, formatAmountForDisplay } from '@/utils/amountDisplay'

const BULLET_MASK = '••••'

describe('amountDisplay mask (L1): •••• when amounts hidden', () => {
  let settingsStore: ReturnType<typeof useSettingsStore>

  beforeEach(() => {
    setActivePinia(createPinia())
    settingsStore = useSettingsStore()
    settingsStore.settings.number_format = 'en-US'
    settingsStore.settings.base_currency = 'USD'
    settingsStore.settings.hide_amounts = false
  })

  describe('formatNumber', () => {
    test('returns •••• (4 bullets) when hide_amounts is true', () => {
      settingsStore.settings.hide_amounts = true
      expect(formatNumber(1000)).toBe(BULLET_MASK)
      expect(formatNumber(0)).toBe(BULLET_MASK)
      expect(formatNumber(-250)).toBe(BULLET_MASK)
    })

    test('does NOT return ***** when hide_amounts is true', () => {
      settingsStore.settings.hide_amounts = true
      expect(formatNumber(1000)).not.toBe('*****')
      expect(formatNumber(0)).not.toBe('*****')
    })

    test('returns formatted number when hide_amounts is false', () => {
      expect(formatNumber(1234567.89)).toBe('1,234,567.89')
    })
  })

  describe('formatAmountForDisplay', () => {
    test('returns sign + symbol + •••• when hide_amounts is true (income)', () => {
      settingsStore.settings.hide_amounts = true
      settingsStore.settings.base_currency = 'USD'
      expect(formatAmountForDisplay(5000, 'income')).toBe('+$••••')
    })

    test('returns sign + symbol + •••• when hide_amounts is true (expense)', () => {
      settingsStore.settings.hide_amounts = true
      settingsStore.settings.base_currency = 'USD'
      expect(formatAmountForDisplay(300, 'expense')).toBe('−$••••')
    })

    test('returns symbol + •••• when hide_amounts is true (transfer)', () => {
      settingsStore.settings.hide_amounts = true
      settingsStore.settings.base_currency = 'USD'
      expect(formatAmountForDisplay(100, 'transfer')).toBe('$••••')
    })

    test('does NOT return ***** when hide_amounts is true', () => {
      settingsStore.settings.hide_amounts = true
      expect(formatAmountForDisplay(5000, 'income')).not.toBe('+$*****')
      expect(formatAmountForDisplay(300, 'expense')).not.toBe('−$*****')
      expect(formatAmountForDisplay(100, 'transfer')).not.toBe('$*****')
    })

    test('returns full formatted amount with symbol when hide_amounts is false', () => {
      expect(formatAmountForDisplay(1500, 'income')).toBe('+$1,500')
    })
  })
})
