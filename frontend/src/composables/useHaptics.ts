/**
 * Haptic feedback composable for mobile devices
 * Uses the Vibration API for tactile feedback
 */

type HapticPattern = number | number[]

export function useHaptics() {
  const isSupported = typeof navigator !== 'undefined' && 'vibrate' in navigator

  const vibrate = (pattern: HapticPattern) => {
    if (isSupported) {
      navigator.vibrate(pattern)
    }
  }

  const light = () => vibrate(10)
  const medium = () => vibrate(20)
  const heavy = () => vibrate([10, 50, 10])
  const success = () => vibrate([10, 50, 10, 50, 10])
  const error = () => vibrate([50, 50, 50])
  const warning = () => vibrate([20, 100, 20])

  return {
    light,
    medium,
    heavy,
    success,
    error,
    warning,
    isSupported
  }
}

// Global haptic functions for easy use
export const haptics = {
  light: () => navigator.vibrate?.(10),
  medium: () => navigator.vibrate?.(20),
  heavy: () => navigator.vibrate?.([10, 50, 10]),
  success: () => navigator.vibrate?.([10, 50, 10, 50, 10]),
  error: () => navigator.vibrate?.([50, 50, 50]),
  warning: () => navigator.vibrate?.([20, 100, 20])
}

// Directive for adding haptics to any element
export const vHaptic = {
  mounted(el: HTMLElement, binding: { value: keyof typeof haptics }) {
    const type = binding.value || 'light'
    const handler = () => haptics[type]?.()
    el.addEventListener('click', handler)
    ;(el as any)._hapticHandler = handler
  },
  unmounted(el: HTMLElement) {
    const handler = (el as any)._hapticHandler
    if (handler) {
      el.removeEventListener('click', handler)
      delete (el as any)._hapticHandler
    }
  }
}