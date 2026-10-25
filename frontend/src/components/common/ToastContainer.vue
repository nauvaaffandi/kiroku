<template>
  <div class="fixed bottom-[calc(72px+env(safe-area-inset-bottom))] right-6 flex flex-col gap-2 z-[400] pointer-events-none md:bottom-8">
    <div
      v-for="toast in uiStore.toasts"
      :key="toast.id"
      class="relative flex items-center gap-3 py-3 px-4 pl-10 rounded-xl shadow-lg min-w-[18rem] max-w-[24rem] pointer-events-auto bg-[#fcfbf9] dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 text-zinc-900 dark:text-zinc-100 animate-[slideIn_0.3s_ease]"
    >
      <div
        class="absolute top-0 left-6 w-3.5 h-7 shadow-xs transition-colors z-10 [clip-path:polygon(0_0,100%_0,100%_100%,50%_75%,0_100%)]"
        :class="{
          'bg-emerald-500 dark:bg-emerald-400': toast.type === 'success',
          'bg-rose-500 dark:bg-rose-400': toast.type === 'error',
          'bg-amber-500 dark:bg-amber-400': toast.type === 'warning',
          'bg-sky-500 dark:bg-sky-400': toast.type === 'info'
        }"
      ></div>

      <span class="flex-1 text-sm font-semibold overflow-hidden text-ellipsis whitespace-nowrap pl-4">{{ toast.message }}</span>

      <button
        @click="removeToast(toast.id)"
        class="p-1 rounded-md bg-transparent border-none text-current opacity-60 hover:opacity-100 cursor-pointer transition-opacity"
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
</style>
