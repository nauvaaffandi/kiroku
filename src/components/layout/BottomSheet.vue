<template>
  <transition name="bs" appear>
    <div v-if="modelValue" class="bs-wrapper">
      <div class="bs-overlay" @click.self="handleScrimClick" />
      <div
        class="bs-sheet"
        @click.self
        @keydown.escape="close"
      >
      <header class="bs-header">
        <h2 class="bs-title">{{ title }}</h2>
        <button
          type="button"
          class="bs-close"
          @click="close"
          :aria-label="t('bottomSheet.close')"
        >
          <X :size="20" stroke-width="1.75" />
        </button>
      </header>
      <div class="bs-body">
        <slot />
      </div>
      <footer v-if="$slots.footer" class="bs-footer">
        <slot name="footer" />
      </footer>
    </div>
  </div>
  </transition>
</template>

<script setup lang="ts">
import { withDefaults } from 'vue'
import { X } from '@lucide/vue'
import { useI18n } from '@/locales'

interface Props {
  modelValue: boolean
  title?: string
  dirty?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  dirty: false
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const { t } = useI18n()

console.log('[BottomSheet] initialized', { title: props.title, dirty: props.dirty })

const close = () => {
  console.log('[BottomSheet] closing')
  emit('update:modelValue', false)
}

const handleScrimClick = () => {
  if (props.dirty) {
    console.log('[BottomSheet] scrim tap ignored (dirty)')
    return
  }
  close()
}
</script>

<style scoped>
@reference "tailwindcss";

.bs-overlay {
  @apply fixed inset-0 z-[500];
  background-color: rgb(0 0 0 / 0.4);
}

[data-theme="dark"] .bs-overlay {
  background-color: rgb(0 0 0 / 0.6);
}

.bs-sheet {
  @apply fixed left-0 right-0 bottom-0 z-[501] w-full max-h-[90vh] rounded-t-[14px] overflow-hidden flex flex-col;
  background-color: var(--bg-card);
  border-top: 1px solid var(--border-primary);
}

.bs-header {
  @apply flex items-center justify-between px-6 py-4 border-b;
  border-color: var(--border-primary);
}

.bs-title {
  @apply text-base font-semibold text-[var(--text-primary)];
}

.bs-close {
  @apply p-2.5 rounded-[10px] flex items-center justify-center min-w-[44px] min-h-[44px] transition-colors duration-120 ease-out;
  color: var(--text-muted);
  background: transparent;
  border: none;
  cursor: pointer;
}

.bs-close:hover {
  @apply bg-[var(--bg-hover)];
  color: var(--text-primary);
}

.bs-close:focus-visible {
  @apply ring-2 ring-offset-2 rounded-[10px];
  ring-offset-color: var(--bg-card);
  ring-color: var(--accent);
}

.bs-body {
  @apply flex-1 overflow-y-auto p-6;
}

.bs-footer {
  @apply p-6 border-t;
  border-color: var(--border-primary);
}

.bs-enter-active,
.bs-leave-active {
  @apply transition duration-240 ease-out;
}

.bs-enter-from {
  @apply translate-y-full;
  opacity: 0;
}

.bs-leave-to {
  @apply translate-y-full;
  opacity: 0;
}

@media (prefers-reduced-motion: reduce) {
  .bs-enter-active,
  .bs-leave-active {
    @apply transition-none duration-0;
  }
}
</style>
