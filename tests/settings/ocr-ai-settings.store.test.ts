import { describe, test, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import type { Settings } from '@/types'

const { mockPut, mockGet, persisted } = vi.hoisted(() => ({
  mockPut: vi.fn(async (_kind: string, _id: string, payload: unknown) => {
    persisted.set('settings', payload as Settings)
  }),
  mockGet: vi.fn(async (_kind: string, _id: string) => (persisted.get('settings') as Settings | undefined) ?? null),
  persisted: new Map<string, Settings>()
}))

vi.mock('@/storage/repository', () => ({
  localDb: { put: mockPut, get: mockGet }
}))

const { SECURE_KEYS, setSecure, getSecure, removeSecure, hasSecure, migratePlaintextKeys, secureValues } = vi.hoisted(() => ({
  SECURE_KEYS: ['ocr_api_key', 'ai_api_key'] as const,
  secureValues: new Map<string, string>(),
  setSecure: vi.fn(async (key: string, value: string) => { secureValues.set(key, value); return true }),
  getSecure: vi.fn(async (key: string) => secureValues.get(key) ?? null),
  removeSecure: vi.fn(async (key: string) => { secureValues.delete(key) }),
  hasSecure: vi.fn(async (key: string) => secureValues.has(key)),
  migratePlaintextKeys: vi.fn(async (settingsRow: Settings) => {
    const redactedKeys: string[] = []
    let migrated = false
    for (const k of ['ocr_api_key', 'ai_api_key'] as const) {
      const v = settingsRow[k]
      if (v) {
        secureValues.set(k, v)
        redactedKeys.push(k)
        migrated = true
      }
    }
    return { migrated, redactedKeys, failedKeys: [] }
  })
}))

vi.mock('@/composables/useSecureStorage', () => ({
  SECURE_KEYS,
  useSecureStorage: vi.fn(() => ({
    setSecure,
    getSecure,
    removeSecure,
    hasSecure,
    migratePlaintextKeys
  }))
}))

import { useSettingsStore } from '@/stores/settings'

describe('settings store — OCR/AI API keys at rest in secure storage (TDD red phase)', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    persisted.clear()
    secureValues.clear()
  })

  describe('defaults', () => {
    test('ocr_api_key defaults to empty string', () => {
      expect(useSettingsStore().settings.ocr_api_key).toBe('')
    })

    test('ai_api_key defaults to empty string', () => {
      expect(useSettingsStore().settings.ai_api_key).toBe('')
    })

    test('ocr_enabled defaults to false', () => {
      expect(useSettingsStore().settings.ocr_enabled).toBe(false)
    })

    test('ai_provider_enabled defaults to false', () => {
      expect(useSettingsStore().settings.ai_provider_enabled).toBe(false)
    })
  })

  describe('setSetting routes API keys to secure storage, not SQLite payload', () => {
    test('setSetting ocr_api_key writes to secure storage and keeps value in store state', async () => {
      const store = useSettingsStore()
      await store.setSetting('ocr_api_key', 'test-key-123')
      expect(setSecure).toHaveBeenCalledWith('ocr_api_key', 'test-key-123')
      expect(store.getSetting('ocr_api_key')).toBe('test-key-123')
    })

    test('setSetting ai_api_key writes to secure storage and keeps value in store state', async () => {
      const store = useSettingsStore()
      await store.setSetting('ai_api_key', 'test-key-456')
      expect(setSecure).toHaveBeenCalledWith('ai_api_key', 'test-key-456')
      expect(store.getSetting('ai_api_key')).toBe('test-key-456')
    })

    test('setSetting ocr_api_key does not write the key into the SQLite payload', async () => {
      const store = useSettingsStore()
      await store.setSetting('ocr_api_key', 'test-key-123')
      mockPut.mock.calls.forEach((call) => {
        const payload = call[2] as Settings
        expect(payload.ocr_api_key).not.toBe('test-key-123')
      })
    })

    test('setSetting non-key field still writes to SQLite and does not touch secure storage', async () => {
      const store = useSettingsStore()
      await store.setSetting('ocr_enabled', true)
      expect(mockPut).toHaveBeenCalledTimes(1)
      const [kind, id, payload] = mockPut.mock.calls[0]
      expect(kind).toBe('setting')
      expect(id).toBe('settings')
      expect((payload as Record<string, unknown>).ocr_enabled).toBe(true)
      expect(setSecure).not.toHaveBeenCalled()
    })
  })

  describe('updateSettings writes non-key fields to SQLite and API keys to secure storage', () => {
    test('API keys are redacted in the SQLite payload and stored in secure storage', async () => {
      const store = useSettingsStore()
      await store.updateSettings({ ocr_api_key: 'test-key-123', ai_api_key: 'test-key-456' })
      expect(setSecure).toHaveBeenCalledWith('ocr_api_key', 'test-key-123')
      expect(setSecure).toHaveBeenCalledWith('ai_api_key', 'test-key-456')
      expect(mockPut).toHaveBeenCalledTimes(1)
      const payload = mockPut.mock.calls[0][2] as Settings
      expect(payload.ocr_api_key).toBe('')
      expect(payload.ai_api_key).toBe('')
      expect(store.settings.ocr_api_key).toBe('test-key-123')
      expect(store.settings.ai_api_key).toBe('test-key-456')
    })

    test('non-key fields are written to the SQLite payload', async () => {
      const store = useSettingsStore()
      await store.updateSettings({ ocr_enabled: true, base_currency: 'EUR', ai_provider_url: 'http://localhost:11434' })
      expect(mockPut).toHaveBeenCalledTimes(1)
      const payload = mockPut.mock.calls[0][2] as Settings
      expect(payload.ocr_enabled).toBe(true)
      expect(payload.base_currency).toBe('EUR')
      expect(payload.ai_provider_url).toBe('http://localhost:11434')
      expect(setSecure).not.toHaveBeenCalled()
    })
  })

  describe('store state returns decrypted API keys from secure storage', () => {
    test('settings.ocr_api_key returns the value from secure storage when SQLite is empty', async () => {
      secureValues.set('ocr_api_key', 'test-key-123')
      mockGet.mockResolvedValue(null)
      const store = useSettingsStore()
      await store.fetchSettings()
      expect(store.settings.ocr_api_key).toBe('test-key-123')
    })

    test('secure storage value takes precedence over an empty SQLite value', async () => {
      secureValues.set('ai_api_key', 'test-key-456')
      mockGet.mockResolvedValue({ ocr_api_key: '', ai_api_key: '' } as Settings)
      const store = useSettingsStore()
      await store.fetchSettings()
      expect(store.settings.ai_api_key).toBe('test-key-456')
    })
  })

  describe('fetchSettings round-trip', () => {
    test('fresh store instance loads API key previously written to secure storage', async () => {
      const store1 = useSettingsStore()
      await store1.setSetting('ocr_api_key', 'test-key-123')
      const store2 = useSettingsStore()
      await store2.fetchSettings()
      expect(store2.settings.ocr_api_key).toBe('test-key-123')
    })
  })

  describe('migration is idempotent', () => {
    test('fetchSettings runs migratePlaintextKeys on the SQLite row', async () => {
      mockGet.mockResolvedValue({
        ocr_api_key: 'test-key-123',
        ai_api_key: '',
        ocr_enabled: true,
        ai_provider_enabled: false
      } as Settings)
      const store = useSettingsStore()
      await store.fetchSettings()
      expect(migratePlaintextKeys).toHaveBeenCalledTimes(1)
      expect(secureValues.get('ocr_api_key')).toBe('test-key-123')
    })

    test('running fetchSettings on an already-redacted row does not re-migrate', async () => {
      mockGet.mockResolvedValue({ ocr_api_key: '', ai_api_key: '' } as Settings)
      const store = useSettingsStore()
      await store.fetchSettings()
      await store.fetchSettings()
      expect(migratePlaintextKeys).toHaveBeenCalledTimes(2)
      const lastRow = migratePlaintextKeys.mock.calls[migratePlaintextKeys.mock.calls.length - 1][0] as Settings
      expect(lastRow.ocr_api_key).toBe('')
      expect(lastRow.ai_api_key).toBe('')
    })
  })
})
