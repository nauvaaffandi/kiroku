import { describe, test, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { defineComponent, nextTick, ref } from 'vue'
import { mount } from '@vue/test-utils'
import { useUIStore } from '@/stores/ui'
import { useSidebarRailGesture } from '@/composables/useSidebarRailGesture'

const Host = defineComponent({
  setup() {
    const rail = ref<HTMLElement | null>(null)
    useSidebarRailGesture(rail)
    return { rail }
  },
  template: '<aside ref="rail"><a class="link" href="/x">go</a></aside>'
})

class Pointerish extends Event {
  pointerId: number
  clientX: number
  clientY: number
  button: number
  timeStamp: number
  bubbles: boolean
  cancelable: boolean
  constructor(type: string, id: number, x: number, y: number, time: number) {
    super(type, { bubbles: true, cancelable: true })
    this.pointerId = id
    this.clientX = x
    this.clientY = y
    this.button = 0
    this.timeStamp = time
  }
  preventDefault() {}
  stopPropagation() {}
  stopImmediatePropagation() {}
}

const ptr = (type: string, id: number, x: number, y: number, time: number) =>
  new Pointerish(type, id, x, y, time)

const gesture = (el: HTMLElement, dx: number, dy: number, holdMs: number) => {
  el.dispatchEvent(ptr('pointerdown', 1, 10, 10, 0))
  el.dispatchEvent(ptr('pointermove', 1, 10 + dx, 10 + dy, holdMs))
  el.dispatchEvent(ptr('pointerup', 1, 10 + dx, 10 + dy, holdMs))
}

describe('useSidebarRailGesture', () => {
  beforeEach(() => setActivePinia(createPinia()))

  const mountRail = async (open: boolean) => {
    const ui = useUIStore()
    ui.sidebarOpen = open
    const wrapper = mount(Host, { attachTo: document.body })
    await nextTick()
    return { ui, el: wrapper.element as HTMLElement, wrapper }
  }

  test('collapsed: hold 100ms + swipe right ≥48px opens', async () => {
    const { ui, el } = await mountRail(false)
    gesture(el, 48, 0, 100)
    expect(ui.sidebarOpen).toBe(true)
  })

  test('collapsed: 50ms hold does not open even past 48px', async () => {
    const { ui, el } = await mountRail(false)
    gesture(el, 60, 0, 50)
    expect(ui.sidebarOpen).toBe(false)
  })

  test('collapsed: hold but below 48px does not open', async () => {
    const { ui, el } = await mountRail(false)
    gesture(el, 20, 0, 120)
    expect(ui.sidebarOpen).toBe(false)
  })

  test('expanded: hold 100ms + swipe left ≥48px closes', async () => {
    const { ui, el } = await mountRail(true)
    gesture(el, -48, 0, 100)
    expect(ui.sidebarOpen).toBe(false)
  })

  test('expanded: swipe right does not close', async () => {
    const { ui, el } = await mountRail(true)
    gesture(el, 60, 0, 120)
    expect(ui.sidebarOpen).toBe(true)
  })

  test('successful gesture swallows the following click', async () => {
    const { el } = await mountRail(false)
    gesture(el, 60, 0, 120)
    const click = new MouseEvent('click', { bubbles: true, cancelable: true })
    const link = el.querySelector('.link') as HTMLElement
    link.dispatchEvent(click)
    expect(click.defaultPrevented).toBe(true)
  })
})
