<template>
  <div ref="root" class="color-picker">
    <button
      type="button"
      class="color-trigger"
      :style="{ '--selected-color': selectedColor }"
      :aria-label="t('colorPicker.open')"
      :aria-expanded="isOpen"
      :disabled="disabled"
      @click="toggle"
    >
      <span class="color-preview" />
    </button>
  </div>
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="color-popover"
      role="listbox"
      :aria-label="t('colorPicker.options')"
    >
      <ColorSwatch
        v-for="color in paletteColors"
        :key="color"
        :color="color"
        :color-active="selectedColor.toLowerCase() === color"
        @update:color="selectColor"
      />
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from '@/locales'
import ColorSwatch from './ColorSwatch.vue'

const props = defineProps<{
  modelValue?: string
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { t } = useI18n()
const root = ref<HTMLElement | null>(null)
const isOpen = ref(false)
const selectedColor = computed(() => props.modelValue || '#3b82f6')
const hues = [0, 20, 40, 55, 75, 100, 125, 145, 165, 185, 200, 215, 235]
const shades = [35, 45, 55, 65]
const paletteColors = hues.flatMap(hue => shades.map(lightness => `hsl(${hue} 85% ${lightness}%)`))

function toggle() {
  if (props.disabled) return
  isOpen.value = !isOpen.value
}

function selectColor(color: string) {
  emit('update:modelValue', color)
  isOpen.value = false
}

function closeOnOutsideClick(event: MouseEvent) {
  if (
    isOpen.value &&
    event.target instanceof Node &&
    !root.value?.contains(event.target) &&
    !(event.target instanceof Element && event.target.closest('.color-popover'))
  ) {
    isOpen.value = false
  }
}

function closeOnEscape(event: KeyboardEvent) {
  if (event.key === 'Escape') isOpen.value = false
}

onMounted(() => {
  document.addEventListener('click', closeOnOutsideClick)
  document.addEventListener('keydown', closeOnEscape)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', closeOnOutsideClick)
  document.removeEventListener('keydown', closeOnEscape)
})
</script>

<style scoped>
@reference "tailwindcss";

.color-picker {
  @apply inline-flex;
}

.color-trigger {
  @apply inline-flex items-center justify-center rounded cursor-pointer w-[2.75rem] h-[2.75rem] bg-[var(--bg-card)] border-[1px] border-[var(--border-primary)];
}

.color-trigger:focus-visible {
  @apply outline-none outline-2 outline-offset-2;
  outline-color: var(--border-focus);
}

.color-preview {
  @apply rounded w-[1.5rem] h-[1.5rem] bg-[var(--selected-color)];
}

.color-popover {
  @apply fixed flex flex-col overflow-y-auto z-[99999] gap-[1px] w-[80vw] max-h-[55vh] rounded-[var(--radius)] bg-[var(--bg-card)] border-[1px] border-[var(--border-primary)] shadow-[var(--shadow-md)] p-[0.5rem];
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.color-popover .color-swatch {
  @apply flex shrink-0 items-center justify-center cursor-pointer w-full h-[3rem] rounded-[var(--radius-sm)] border-[2px] border-transparent;
}

.color-popover .color-swatch.active {
  @apply border-[var(--accent)];
}

.color-popover .color-swatch:focus-visible {
  @apply outline-none outline-2 outline-offset-[-2px];
  outline-color: var(--border-focus);
}

.color-popover .color-swatch:disabled {
  @apply cursor-not-allowed opacity-50;
}
</style>