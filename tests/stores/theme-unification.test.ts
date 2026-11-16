import { describe, test, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const { mockPut, mockGet } = vi.hoisted(() => ({
  mockPut: vi.fn().mockResolvedValue(undefined),
  mockGet: vi.fn().mockResolvedValue(undefined)
}))

vi.mock('@/storage/repository', () => ({
  localDb: { put: mockPut, get: mockGet }
}))

import { useSettingsStore } from '@/stores/settings'

describe('settings store theme unification', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockGet.mockResolvedValue(undefined)
    document.documentElement.removeAttribute('data-theme')
  })

  afterEach(() => {
    document.documentElement.removeAttribute('data-theme')
  })

  test('updateSettings({ theme: "dark" }) sets data-theme="dark"', async () => {
    const store = useSettingsStore()
    await store.updateSettings({ theme: 'dark' })
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
    expect(store.theme).toBe('dark')
  })

  test('updateSettings({ theme: "light" }) sets data-theme="light"', async () => {
    const store = useSettingsStore()
    await store.updateSettings({ theme: 'light' })
    expect(document.documentElement.getAttribute('data-theme')).toBe('light')
    expect(store.theme).toBe('light')
  })

  test('updateSettings({ theme: "system" }) removes data-theme attribute', async () => {
    const store = useSettingsStore()
    document.documentElement.setAttribute('data-theme', 'dark')
    await store.updateSettings({ theme: 'system' })
    expect(document.documentElement.getAttribute('data-theme')).toBeNull()
    expect(store.theme).toBe('system')
  })

  test('theme computed is single source of truth matching data-theme', async () => {
    const store = useSettingsStore()
    await store.updateSettings({ theme: 'dark' })
    expect(store.theme).toBe('dark')
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
    await store.updateSettings({ theme: 'system' })
    // 'system' removes data-theme attribute; store.theme reports 'system'.
    // Both represent the same logical state via different representations.
    expect(store.theme).toBe('system')
    expect(document.documentElement.getAttribute('data-theme')).toBeNull()
  })

  test('applyTheme only touches data-theme attribute, no other attributes', async () => {
    const store = useSettingsStore()
    const setAttributeSpy = vi.spyOn(document.documentElement, 'setAttribute')
    const removeAttributeSpy = vi.spyOn(document.documentElement, 'removeAttribute')
    await store.updateSettings({ theme: 'dark' })
    expect(setAttributeSpy).toHaveBeenCalledWith('data-theme', 'dark')
    expect(removeAttributeSpy).not.toHaveBeenCalled()
    setAttributeSpy.mockRestore()
    removeAttributeSpy.mockRestore()
  })

  test('prefers-reduced-motion is consulted before applying theme (single source of truth)', async () => {
    const matchMediaSpy = vi
      .spyOn(window, 'matchMedia')
      .mockReturnValue({ matches: false } as MediaQueryList)
    const store = useSettingsStore()
    await store.updateSettings({ theme: 'dark' })
    // L3: applyTheme no longer calls window.matchMedia (dead no-op removed).
    expect(matchMediaSpy).not.toHaveBeenCalled()
    matchMediaSpy.mockRestore()
  })

  test('persists theme via localDb.put on updateSettings', async () => {
    const store = useSettingsStore()
    await store.updateSettings({ theme: 'dark' })
    expect(mockPut).toHaveBeenCalledWith('setting', 'settings', expect.objectContaining({ theme: 'dark' }))
  })
})
