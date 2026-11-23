import { onBeforeUnmount, watch, type Ref } from 'vue'
import { useUIStore } from '@/stores/ui'

const HOLD_MS = 100
const DISTANCE_PX = 48

export function useSidebarRailGesture(rail: Ref<HTMLElement | null>) {
  const uiStore = useUIStore()

  let pointerId = -1
  let startX = 0
  let startY = 0
  let startAt = 0
  let decided = false

  const reset = () => {
    pointerId = -1
  }

  const onPointerDown = (e: PointerEvent) => {
    if (e.button !== 0 || pointerId !== -1) return
    pointerId = e.pointerId
    startX = e.clientX
    startY = e.clientY
    startAt = e.timeStamp
    decided = false
    rail.value?.setPointerCapture?.(e.pointerId)
  }

  const onPointerMove = (e: PointerEvent) => {
    if (e.pointerId !== pointerId || decided) return
    const dx = e.clientX - startX
    const dy = e.clientY - startY
    const held = e.timeStamp - startAt >= HOLD_MS

    const open = uiStore.sidebarOpen
    const hit = held
      ? open
        ? dx <= -DISTANCE_PX && -dx > Math.abs(dy)
        : dx >= DISTANCE_PX && dx > Math.abs(dy)
      : false
    if (!hit) return
    uiStore.sidebarOpen = !open
    decided = true
  }

  const onPointerUp = (e: PointerEvent) => {
    if (e.pointerId !== pointerId) return
    reset()
  }

  const onClick = (e: MouseEvent) => {
    if (!decided) return
    decided = false
    e.preventDefault()
    e.stopPropagation()
  }

  const bind = (el: HTMLElement) => {
    el.addEventListener('pointerdown', onPointerDown)
    el.addEventListener('pointermove', onPointerMove)
    el.addEventListener('pointerup', onPointerUp)
    el.addEventListener('pointercancel', onPointerUp)
    el.addEventListener('click', onClick, true)
  }

  const unbind = (el: HTMLElement) => {
    el.removeEventListener('pointerdown', onPointerDown)
    el.removeEventListener('pointermove', onPointerMove)
    el.removeEventListener('pointerup', onPointerUp)
    el.removeEventListener('pointercancel', onPointerUp)
    el.removeEventListener('click', onClick, true)
  }

  watch(rail, (el, prev) => {
    if (prev) unbind(prev)
    if (el) bind(el)
  }, { immediate: true })

  onBeforeUnmount(() => {
    if (rail.value) unbind(rail.value)
  })
}
