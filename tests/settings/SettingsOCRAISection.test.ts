import { describe, test, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'

vi.mock('@/storage/repository', () => ({
  localDb: {
    put: vi.fn().mockResolvedValue(undefined),
    get: vi.fn().mockResolvedValue(null)
  }
}))

import { useSettingsStore } from '@/stores/settings'
import SettingsOCRAISection from '@/components/settings/SettingsOCRAISection.vue'

describe('SettingsOCRAISection component (TDD red phase)', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  test('renders OCR & AI provider section heading', () => {
    const wrapper = mount(SettingsOCRAISection)
    expect(wrapper.text()).toContain('OCR')
    expect(wrapper.text()).toContain('AI')
  })

  test('renders OCR toggle and OCR API key input', () => {
    const wrapper = mount(SettingsOCRAISection)
    expect(wrapper.find('input[type="checkbox"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
  })

  test('renders AI toggle, provider URL, model, and API key inputs', () => {
    const wrapper = mount(SettingsOCRAISection)
    const inputs = wrapper.findAll('input')
    expect(inputs.length).toBeGreaterThanOrEqual(4)
  })

  test('renders local-only warning banner', () => {
    const wrapper = mount(SettingsOCRAISection)
    expect(wrapper.text()).toContain('locally only')
  })

  test('toggling OCR calls store setSetting', async () => {
    const store = useSettingsStore()
    const setSettingSpy = vi.spyOn(store, 'setSetting')
    const wrapper = mount(SettingsOCRAISection)
    await wrapper.find('input[type="checkbox"]').setValue(true)
    expect(setSettingSpy).toHaveBeenCalledWith('ocr_enabled', true)
  })
})
