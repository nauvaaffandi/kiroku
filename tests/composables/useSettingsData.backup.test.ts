import { describe, test, expect, vi, beforeEach, afterEach } from 'vitest'

const mockCreateBackup = vi.fn()
const mockSuccess = vi.fn()
const mockError = vi.fn()

vi.mock('@/stores/settings', () => ({
  useSettingsStore: vi.fn(() => ({ createBackup: mockCreateBackup }))
}))
vi.mock('@/stores/ui', () => ({
  useUIStore: vi.fn(() => ({ success: mockSuccess, error: mockError }))
}))
vi.mock('@/locales', () => ({
  useI18n: vi.fn(() => ({ t: (k: string) => k }))
}))

import { useSettingsData } from '@/composables/useSettingsData'

const fakeBlob = new Blob(['zip-bytes'])

let createElementSpy: ReturnType<typeof vi.spyOn>
let clickSpy: ReturnType<typeof vi.spyOn>
let mockCreateObjectURL: ReturnType<typeof vi.fn>
let mockRevokeObjectURL: ReturnType<typeof vi.fn>

beforeEach(() => {
  vi.clearAllMocks()
  mockCreateBackup.mockResolvedValue(fakeBlob)
  // Let createElement build real DOM nodes so appendChild/removeChild work in jsdom
  createElementSpy = vi.spyOn(document, 'createElement')
  clickSpy = vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => {})
  // jsdom does not implement these; assign directly
  mockCreateObjectURL = vi.fn().mockReturnValue('blob:fake-url')
  mockRevokeObjectURL = vi.fn()
  globalThis.URL.createObjectURL = mockCreateObjectURL
  globalThis.URL.revokeObjectURL = mockRevokeObjectURL
})

afterEach(() => {
  createElementSpy.mockRestore()
  clickSpy.mockRestore()
})

describe('useSettingsData createBackup', () => {
  test('triggers a download anchor click', async () => {
    const { createBackup } = useSettingsData()
    await createBackup()
    expect(createElementSpy).toHaveBeenCalledWith('a')
    expect(clickSpy).toHaveBeenCalledTimes(1)
  })

  test('calls uiStore.success on success', async () => {
    const { createBackup } = useSettingsData()
    await createBackup()
    expect(mockSuccess).toHaveBeenCalledTimes(1)
    expect(mockError).not.toHaveBeenCalled()
  })

  test('creates and revokes object URL', async () => {
    const { createBackup } = useSettingsData()
    await createBackup()
    expect(mockCreateObjectURL).toHaveBeenCalledWith(fakeBlob)
    expect(mockRevokeObjectURL).toHaveBeenCalledTimes(1)
  })

  test('calls uiStore.error on failure and does not leak object URL', async () => {
    mockCreateBackup.mockRejectedValue(new Error('backup failed'))
    const { createBackup } = useSettingsData()
    await createBackup()
    expect(mockError).toHaveBeenCalledTimes(1)
    expect(mockSuccess).not.toHaveBeenCalled()
    expect(mockRevokeObjectURL).not.toHaveBeenCalled()
  })
})
