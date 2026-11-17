import { describe, test, expect, vi, beforeEach } from 'vitest'
import type { Settings } from '@/types'

const mockSet = vi.fn(async () => true)
const mockGet = vi.fn(async () => null)
const mockRemove = vi.fn(async () => true)
const mockKeys = vi.fn(async () => ({ keys: [] }))

vi.mock('capacitor-secure-storage-plugin', () => ({
  SecureStorage: {
    set: mockSet,
    get: mockGet,
    remove: mockRemove,
    keys: mockKeys
  }
}))

import { useSecureStorage } from '@/composables/useSecureStorage'

describe('useSecureStorage — at-rest encryption wrapper for OCR/AI API keys (TDD red phase)', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockGet.mockResolvedValue(null)
    mockKeys.mockResolvedValue({ keys: [] })
    mockSet.mockResolvedValue(true)
    mockRemove.mockResolvedValue(true)
  })

  describe('setSecure', () => {
    test('writes key/value to the secure storage plugin', async () => {
      const { setSecure } = useSecureStorage()
      await setSecure('ocr_api_key', 'test-key-123')
      expect(mockSet).toHaveBeenCalledTimes(1)
      expect(mockSet).toHaveBeenCalledWith({ key: 'ocr_api_key', value: 'test-key-123' })
    })

    test('resolves void on success', async () => {
      const { setSecure } = useSecureStorage()
      await expect(setSecure('ai_api_key', 'test-key-456')).resolves.toBe(true)
    })
  })

  describe('getSecure', () => {
    test('returns the decrypted value when the plugin has the key', async () => {
      mockGet.mockResolvedValue({ value: 'test-key-123' })
      const { getSecure } = useSecureStorage()
      const result = await getSecure('ocr_api_key')
      expect(mockGet).toHaveBeenCalledWith({ key: 'ocr_api_key' })
      expect(result).toBe('test-key-123')
    })

    test('returns null when the plugin returns null', async () => {
      mockGet.mockResolvedValue(null)
      const { getSecure } = useSecureStorage()
      expect(await getSecure('ai_api_key')).toBeNull()
    })

    test('returns null when the plugin returns an object without a value', async () => {
      mockGet.mockResolvedValue({})
      const { getSecure } = useSecureStorage()
      expect(await getSecure('ocr_api_key')).toBeNull()
    })
  })

  describe('removeSecure', () => {
    test('removes the key from the secure storage plugin', async () => {
      const { removeSecure } = useSecureStorage()
      await removeSecure('ocr_api_key')
      expect(mockRemove).toHaveBeenCalledTimes(1)
      expect(mockRemove).toHaveBeenCalledWith({ key: 'ocr_api_key' })
    })

    test('resolves void on success', async () => {
      const { removeSecure } = useSecureStorage()
      await expect(removeSecure('ai_api_key')).resolves.toBe(true)
    })
  })

  describe('hasSecure', () => {
    test('returns true when the key is in the plugin keys list', async () => {
      mockKeys.mockResolvedValue({ keys: ['ocr_api_key', 'ai_api_key'] })
      const { hasSecure } = useSecureStorage()
      expect(await hasSecure('ocr_api_key')).toBe(true)
      expect(mockKeys).toHaveBeenCalledTimes(1)
    })

    test('returns false when the key is not in the plugin keys list', async () => {
      mockKeys.mockResolvedValue({ keys: ['ai_api_key'] })
      const { hasSecure } = useSecureStorage()
      expect(await hasSecure('ocr_api_key')).toBe(false)
    })
  })

  describe('migratePlaintextKeys', () => {
    test('migrates non-empty ocr_api_key and ai_api_key into secure storage', async () => {
      const { migratePlaintextKeys } = useSecureStorage()
      const row = { ocr_api_key: 'test-key-123', ai_api_key: 'test-key-456', ocr_enabled: true } as Settings
      const result = await migratePlaintextKeys(row)
      expect(result.migrated).toBe(true)
      expect(result.redactedKeys).toEqual(expect.arrayContaining(['ocr_api_key', 'ai_api_key']))
      expect(mockSet).toHaveBeenCalledWith({ key: 'ocr_api_key', value: 'test-key-123' })
      expect(mockSet).toHaveBeenCalledWith({ key: 'ai_api_key', value: 'test-key-456' })
    })

    test('migrates nothing and returns migrated false when both keys are empty', async () => {
      const { migratePlaintextKeys } = useSecureStorage()
      const result = await migratePlaintextKeys({ ocr_api_key: '', ai_api_key: '' } as Settings)
      expect(result).toEqual({ migrated: false, redactedKeys: [], failedKeys: [] })
      expect(mockSet).not.toHaveBeenCalled()
    })

    test('migrates only the non-empty key', async () => {
      const { migratePlaintextKeys } = useSecureStorage()
      const result = await migratePlaintextKeys({ ocr_api_key: 'test-key-123', ai_api_key: '' } as Settings)
      expect(result.migrated).toBe(true)
      expect(result.redactedKeys).toEqual(['ocr_api_key'])
      expect(mockSet).toHaveBeenCalledTimes(1)
      expect(mockSet).toHaveBeenCalledWith({ key: 'ocr_api_key', value: 'test-key-123' })
    })
  })

  describe('graceful fallback when the plugin is unavailable', () => {
    test('getSecure returns null instead of throwing when the plugin rejects', async () => {
      mockGet.mockRejectedValue(new Error('secure storage unavailable'))
      const { getSecure } = useSecureStorage()
      await expect(getSecure('ocr_api_key')).resolves.toBeNull()
    })

    test('setSecure resolves void instead of throwing when the plugin rejects', async () => {
      mockSet.mockRejectedValue(new Error('secure storage unavailable'))
      const { setSecure } = useSecureStorage()
      await expect(setSecure('ocr_api_key', 'test-key-123')).resolves.toBe(false)
    })

    test('hasSecure returns false when keys() rejects', async () => {
      mockKeys.mockRejectedValue(new Error('secure storage unavailable'))
      const { hasSecure } = useSecureStorage()
      expect(await hasSecure('ocr_api_key')).toBe(false)
    })
  })
})
