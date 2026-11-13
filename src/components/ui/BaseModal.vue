<template>
  <Transition name="modal">
    <div v-if="modelValue" class="modal-overlay" @click.self="handleOverlayClick">
      <div class="modal" :class="size" role="dialog" aria-modal="true" :aria-labelledby="title ? `modal-title-${id}` : undefined">
        <div v-if="title || closable" class="modal-header">
          <h2 v-if="title" :id="`modal-title-${id}`" class="modal-title">{{ title }}</h2>
          <button v-if="closable" @click="close" class="modal-close" aria-label="Close">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <slot />
        </div>
        <div v-if="$slots.actions" class="modal-actions">
          <slot name="actions" />
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
interface Props {
  modelValue: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  closable?: boolean
  persistent?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  closable: true,
  persistent: false,
  size: 'md'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const id = Math.random().toString(36).slice(2)

const close = () => {
  if (!props.persistent) {
    emit('update:modelValue', false)
  }
}

const handleOverlayClick = () => {
  close()
}
</script>

<style scoped>
@reference "tailwindcss";

.modal-overlay {
  @apply fixed inset-0 flex items-center justify-center p-4 z-[500];
  background-color: rgba(0, 0, 0, 0.5);
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal {
  background-color: var(--bg-card);
  @apply rounded-lg flex flex-col w-full max-h-[90vh];
  box-shadow: var(--shadow-xl);
  animation: slideUp 0.2s ease;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(1rem);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-sm { @apply max-w-[320px]; }
.modal-md { @apply max-w-[480px]; }
.modal-lg { @apply max-w-[640px]; }
.modal-xl { @apply max-w-[800px]; }
.modal-full { @apply max-w-[95vw] max-h-[95vh]; }

.modal-header {
  @apply flex items-center justify-between px-5 py-4 border-b;
  border-color: var(--border-primary);
}

.modal-title {
  @apply m-0 text-lg font-semibold;
  color: var(--text-primary);
}

.modal-close {
  @apply flex items-center justify-center w-8 h-8 border-0 bg-transparent rounded cursor-pointer transition-colors;
  color: var(--text-muted);
}

.modal-close:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.modal-body {
  @apply p-5 overflow-y-auto flex-1;
}

.modal-actions {
  @apply flex justify-end gap-3 px-5 py-4 border-t;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
  border-radius: 0 0 var(--radius-lg) var(--radius-lg);
}
</style>