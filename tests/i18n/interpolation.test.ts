import { describe, test, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useI18nStore } from '@/locales'

describe('i18n double-brace interpolation ({{param}})', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  test('interpolates {{count}} double-brace parameters in transactions.summaryCount', () => {
    const i18n = useI18nStore()
    i18n.setLocale('en')
    const result = i18n.t('transactions.summaryCount', { count: 5 })
    // Expected: "5 transactions", but CURRENTLY BUGGY: "{5} transactions"
    expect(result).toBe('5 transactions')
  })

  test('interpolates {{count}} in transactions.selectedCount', () => {
    const i18n = useI18nStore()
    i18n.setLocale('en')
    const result = i18n.t('transactions.selectedCount', { count: 12 })
    expect(result).toBe('12 selected')
  })

  test('interpolates {{name}} double-brace parameter in categories.subcategoriesTitle', () => {
    const i18n = useI18nStore()
    i18n.setLocale('en')
    const result = i18n.t('categories.subcategoriesTitle', { name: 'Food' })
    expect(result).toBe('Subcategories of Food')
  })

  test('interpolates multiple double-brace parameters in common.page', () => {
    const i18n = useI18nStore()
    i18n.setLocale('en')
    const result = i18n.t('common.page', { current: 1, total: 10 })
    expect(result).toBe('Page 1 of 10')
  })

  test('interpolates across non-English locales (id)', () => {
    const i18n = useI18nStore()
    i18n.setLocale('id')
    const result = i18n.t('transactions.summaryCount', { count: 3 })
    // In id.ts: summaryCount is '{{count}} transaksi'
    expect(result).toBe('3 transaksi')
  })

  test('interpolates across non-English locales (ja)', () => {
    const i18n = useI18nStore()
    i18n.setLocale('ja')
    const result = i18n.t('transactions.summaryCount', { count: 42 })
    // In ja.ts: summaryCount is '{{count}}件の取引'
    expect(result).toBe('42件の取引')
  })

  test('supports single-brace {param} backward-compatibility', () => {
    const i18n = useI18nStore()
    i18n.setLocale('en')
    const result = i18n.t('calendar.monthTitle', { month: 'November', year: 2026 })
    expect(result).toBe('November 2026')
  })
})
