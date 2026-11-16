import { describe, test, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useSettingsStore } from '@/stores/settings'

vi.mock('@/storage/repository', () => ({
  localDb: {
    put: vi.fn().mockResolvedValue(undefined),
    get: vi.fn().mockResolvedValue(undefined),
    list: vi.fn().mockResolvedValue([])
  }
}))

describe('settings applyTheme (L3): no window.matchMedia call', () => {
  let matchMediaSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    setActivePinia(createPinia())
    matchMediaSpy = vi.fn()
    ;(window as any).matchMedia = matchMediaSpy
  })

  test('applyTheme does NOT call window.matchMedia when theme changes to dark', async () => {
    const store = useSettingsStore()
    await store.updateSettings({ theme: 'dark' })

    expect(matchMediaSpy).not.toHaveBeenCalled()
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
  })

  test('applyTheme does NOT call window.matchMedia when theme changes to system', async () => {
    const store = useSettingsStore()
    await store.updateSettings({ theme: 'system' })

    expect(matchMediaSpy).not.toHaveBeenCalled()
    expect(document.documentElement.hasAttribute('data-theme')).toBe(false)
  })
})
