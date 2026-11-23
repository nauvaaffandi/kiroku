import { describe, test, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useUIStore } from '@/stores/ui'

describe('ui store — sidebarOpen', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  test('sidebarOpen defaults to false (collapsed)', () => {
    const store = useUIStore()
    expect(store.sidebarOpen).toBe(false)
  })

  test('toggleSidebar flips sidebarOpen', () => {
    const store = useUIStore()
    expect(store.sidebarOpen).toBe(false)
    store.toggleSidebar()
    expect(store.sidebarOpen).toBe(true)
    store.toggleSidebar()
    expect(store.sidebarOpen).toBe(false)
  })

  test('true means expanded, false means collapsed', () => {
    const store = useUIStore()
    store.sidebarOpen = true
    expect(store.sidebarOpen).toBe(true)
    store.sidebarOpen = false
    expect(store.sidebarOpen).toBe(false)
  })
})
