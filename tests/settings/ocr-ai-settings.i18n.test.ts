import { describe, test, expect } from 'vitest'

import en from '@/locales/en'
import ja from '@/locales/ja'
import id from '@/locales/id'

const NEW_KEYS = [
  'ocrTitle',
  'ocrToggle',
  'ocrApiKey',
  'aiTitle',
  'aiToggle',
  'aiProviderUrl',
  'aiModel',
  'aiApiKey',
  'aiLocalOnlyWarning'
] as const

function settingsOf(locale: Record<string, unknown>) {
  return (locale as Record<string, Record<string, unknown>>).settings
}

describe('i18n OCR & AI provider settings keys (TDD red phase)', () => {
  test.each(NEW_KEYS)('settings.%s exists and is a non-empty string in en', (key) => {
    const val = settingsOf(en)[key]
    expect(typeof val).toBe('string')
    expect(val.length).toBeGreaterThan(0)
  })

  test.each(NEW_KEYS)('settings.%s exists and is a non-empty string in ja', (key) => {
    const val = settingsOf(ja)[key]
    expect(typeof val).toBe('string')
    expect(val.length).toBeGreaterThan(0)
  })

  test.each(NEW_KEYS)('settings.%s exists and is a non-empty string in id', (key) => {
    const val = settingsOf(id)[key]
    expect(typeof val).toBe('string')
    expect(val.length).toBeGreaterThan(0)
  })

  test.each(NEW_KEYS)('no English leakage: ja.settings.%s !== en.settings.%s', (key) => {
    expect(settingsOf(ja)[key]).not.toBe(settingsOf(en)[key])
  })

  test.each(NEW_KEYS)('no English leakage: id.settings.%s !== en.settings.%s', (key) => {
    expect(settingsOf(id)[key]).not.toBe(settingsOf(en)[key])
  })
})
