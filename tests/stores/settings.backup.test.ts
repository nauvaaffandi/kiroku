import { describe, test, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const { mockList, mockFile, mockGenerateAsync } = vi.hoisted(() => ({
  mockList: vi.fn(),
  mockFile: vi.fn(),
  mockGenerateAsync: vi.fn()
}))

vi.mock('@/storage/repository', () => ({
  localDb: { list: mockList }
}))
vi.mock('jszip', () => ({
  default: vi.fn().mockImplementation(() => ({
    file: mockFile,
    generateAsync: mockGenerateAsync
  }))
}))

import { useSettingsStore } from '@/stores/settings'

const ENTITY_KINDS = [
  'account',
  'category',
  'label',
  'notification',
  'transaction',
  'template',
  'scheduled',
  'budget',
  'goal',
  'setting'
] as const

const fakeBlob = new Blob(['fake-zip-bytes'])

beforeEach(() => {
  vi.clearAllMocks()
  mockList.mockResolvedValue([])
  mockFile.mockReturnValue(undefined)
  mockGenerateAsync.mockResolvedValue(fakeBlob)
  setActivePinia(createPinia())
})

describe('settings store createBackup', () => {
  test('returns a Blob', async () => {
    const store = useSettingsStore()
    const result = await store.createBackup()
    expect(result).toBeInstanceOf(Blob)
  })

  test('blob is not empty', async () => {
    const store = useSettingsStore()
    const result = await store.createBackup()
    expect(result.size).toBeGreaterThan(0)
  })

  test('calls localDb.list for all 10 entity kinds', async () => {
    const store = useSettingsStore()
    await store.createBackup()
    expect(mockList).toHaveBeenCalledTimes(ENTITY_KINDS.length)
    ENTITY_KINDS.forEach((kind) => {
      expect(mockList).toHaveBeenCalledWith(kind)
    })
  })

  test('uses JSZip.generateAsync to build the zip blob', async () => {
    const store = useSettingsStore()
    await store.createBackup()
    expect(mockGenerateAsync).toHaveBeenCalledTimes(1)
    expect(mockGenerateAsync).toHaveBeenCalledWith(expect.objectContaining({ type: 'blob' }))
  })

  test('propagates localDb errors so the caller can surface them', async () => {
    mockList.mockRejectedValueOnce(new Error('database read failed'))
    const store = useSettingsStore()
    await expect(store.createBackup()).rejects.toThrow('database read failed')
  })
})
