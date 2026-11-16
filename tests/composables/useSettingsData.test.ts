import { describe, test, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const { mockList } = vi.hoisted(() => ({
  mockList: vi.fn().mockResolvedValue([])
}))

vi.mock('@/storage/repository', () => ({
  localDb: { list: mockList, put: vi.fn(), get: vi.fn() }
}))
vi.mock('@/stores/ui', () => ({
  useUIStore: vi.fn(() => ({ success: vi.fn(), error: vi.fn() }))
}))
vi.mock('@/locales', () => ({
  useI18n: vi.fn(() => ({ t: (k: string) => k }))
}))

import { useSettingsData } from '@/composables/useSettingsData'

let logSpy: ReturnType<typeof vi.spyOn>
let originalCreateObjectURL: typeof URL.createObjectURL
let originalRevokeObjectURL: typeof URL.revokeObjectURL

beforeEach(() => {
  setActivePinia(createPinia())
  vi.clearAllMocks()
  mockList.mockResolvedValue([])
  logSpy = vi.spyOn(console, 'log').mockImplementation(() => {})
  originalCreateObjectURL = globalThis.URL.createObjectURL
  originalRevokeObjectURL = globalThis.URL.revokeObjectURL
  globalThis.URL.createObjectURL = vi.fn().mockReturnValue('blob:fake-url')
  globalThis.URL.revokeObjectURL = vi.fn()
})

afterEach(() => {
  logSpy.mockRestore()
  globalThis.URL.createObjectURL = originalCreateObjectURL
  globalThis.URL.revokeObjectURL = originalRevokeObjectURL
})

describe('useSettingsData console.log tracing (get/save/update/delete)', () => {
  test('createBackup logs initiation of the get/save flow', async () => {
    const { createBackup } = useSettingsData()
    await createBackup()
    expect(logSpy).toHaveBeenCalledWith('[SettingsData] Initiating backup download')
  })

  test('createBackup traces get operations via settings store createBackup', async () => {
    const { createBackup } = useSettingsData()
    await createBackup()
    expect(logSpy).toHaveBeenCalledWith('[Backup] Starting backup creation')
    expect(logSpy).toHaveBeenCalledWith(expect.stringContaining('[Backup] Fetching data for kind:'))
    expect(logSpy).toHaveBeenCalledWith(expect.stringContaining('[Backup] Retrieved'))
  })

  test('createBackup traces save operation (ZIP blob creation)', async () => {
    const { createBackup } = useSettingsData()
    await createBackup()
    expect(logSpy).toHaveBeenCalledWith(expect.stringContaining('[Backup] ZIP created successfully'))
  })

  test('createBackup logs completion message', async () => {
    const { createBackup } = useSettingsData()
    await createBackup()
    expect(logSpy).toHaveBeenCalledWith(expect.stringContaining('[SettingsData] Backup downloaded as'))
  })

  test('createBackup traces a get for all 10 entity kinds', async () => {
    const { createBackup } = useSettingsData()
    await createBackup()
    const kindLogs = logSpy.mock.calls.filter(
      ([arg]) => typeof arg === 'string' && (arg as string).includes('Fetching data for kind:')
    )
    expect(kindLogs.length).toBeGreaterThanOrEqual(10)
  })

  test('error path logs via console.error (not console.log) on failure', async () => {
    mockList.mockRejectedValueOnce(new Error('db read failed'))
    const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    const { createBackup } = useSettingsData()
    await createBackup()
    expect(errorSpy).toHaveBeenCalledWith(
      expect.stringContaining('[SettingsData] Failed to create backup'),
      expect.any(Error)
    )
    errorSpy.mockRestore()
  })
})
