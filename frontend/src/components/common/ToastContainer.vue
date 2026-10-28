<template>
  <div class="toast-container">
    <div
      v-for="toast in uiStore.toasts"
      :key="toast.id"
      class="toast"
      :class="`toast-${toast.type}`"
    >
      <div
        class="toast-accent"
        :class="`toast-accent-${toast.type}`"
      ></div>

      <span class="toast-message">{{ toast.message }}</span>

      <button
        @click="removeToast(toast.id)"
        class="toast-dismiss"
        aria-label="Dismiss"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUIStore } from '@/stores/ui'

const uiStore = useUIStore()

const removeToast = (id: number) => {
  uiStore.toasts = uiStore.toasts.filter((t: { id: number }) => t.id !== id)
}
</script>

<style scoped>
@reference "tailwindcss";
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(100%);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.toast-container {
  @apply fixed bottom-[calc(72px+env(safe-area-inset-bottom))] right-6 flex flex-col gap-2 z-[400] pointer-events-none md:bottom-8;
}

.toast {
  @apply relative flex items-center gap-3 py-3 px-4 pl-10 rounded-xl shadow-lg min-w-[18rem] max-w-[24rem] pointer-events-auto border animate-[slideIn_0.3s_ease];
  background-color: var(--bg-hover);
  color: var(--text-primary);
  border-color: var(--border-secondary);
}

.toast-accent {
  @apply absolute top-0 left-6 w-3.5 h-7 shadow-xs transition-colors z-10;
  clip-path: polygon(0 0, 100% 0, 100% 100%, 50% 75%, 0 100%);
}

.toast-accent-success {
  background-color: var(--success);
}

.toast-accent-error {
  background-color: var(--error);
}

.toast-accent-warning {
  background-color: var(--warning);
}

.toast-accent-info {
  background-color: var(--info);
}

.toast-message {
  @apply flex-1 text-sm font-semibold overflow-hidden text-ellipsis whitespace-nowrap pl-4;
}

.toast-dismiss {
  @apply p-1 rounded-md border-0 bg-transparent text-current opacity-60 hover:opacity-100 cursor-pointer transition-opacity;
}
</style>
