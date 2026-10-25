<template>
  <div class="toast-container">
    <div
      v-for="toast in uiStore.toasts"
      :key="toast.id"
      :class="['toast', toast.type]"
    >
      <div class="toast-icon">
        <component :is="getIcon(toast.type)" />
      </div>
      <span class="toast-message">{{ toast.message }}</span>
      <button @click="removeToast(toast.id)" class="toast-close" aria-label="Dismiss">
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

import SuccessIcon from '@/components/icons/SuccessIcon.vue'
import ErrorIcon from '@/components/icons/ErrorIcon.vue'
import WarningIcon from '@/components/icons/WarningIcon.vue'
import InfoIcon from '@/components/icons/InfoIcon.vue'

const uiStore = useUIStore()

const getIcon = (type: 'success' | 'error' | 'warning' | 'info') => {
  switch (type) {
    case 'success': return SuccessIcon
    case 'error': return ErrorIcon
    case 'warning': return WarningIcon
    default: return InfoIcon
  }
}

const removeToast = (id: number) => {
  uiStore.toasts = uiStore.toasts.filter((t: { id: number }) => t.id !== id)
}
</script>

<style scoped>
.toast-container {
  position: fixed;
  bottom: calc(72px + env(safe-area-inset-bottom));
  right: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  z-index: 400;
  pointer-events: none;
}

.toast {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: var(--radius);
  box-shadow: var(--shadow-lg);
  min-width: 18rem;
  max-width: 24rem;
  animation: slideIn 0.3s ease;
  pointer-events: auto;
}

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

.toast.success {
  background-color: #dcfce7;
  border-left: 4px solid #16a34a;
  color: #166534;
}

:where([data-theme="dark"], [data-theme="dark"] *) .toast.success {
  background-color: #064e3b;
  color: #a7f3d0;
}

.toast.error {
  background-color: #fef2f2;
  border-left: 4px solid #dc2626;
  color: #991b1b;
}

:where([data-theme="dark"], [data-theme="dark"] *) .toast.error {
  background-color: #7f1d1d;
  color: #fecaca;
}

.toast.warning {
  background-color: #fffbeb;
  border-left: 4px solid #f59e0b;
  color: #92400e;
}

:where([data-theme="dark"], [data-theme="dark"] *) .toast.warning {
  background-color: #78350f;
  color: #fde68a;
}

.toast.info {
  background-color: #eff6ff;
  border-left: 4px solid #3b82f6;
  color: #1e40af;
}

:where([data-theme="dark"], [data-theme="dark"] *) .toast.info {
  background-color: #1e3a8a;
  color: #bfdbfe;
}

.toast-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
}

.toast.success .toast-icon { color: #16a34a; }
.toast.error .toast-icon { color: #dc2626; }
.toast.warning .toast-icon { color: #f59e0b; }
.toast.info .toast-icon { color: #3b82f6; }

:where([data-theme="dark"], [data-theme="dark"] *) .toast.success .toast-icon { color: #34d399; }
:where([data-theme="dark"], [data-theme="dark"] *) .toast.error .toast-icon { color: #f87171; }
:where([data-theme="dark"], [data-theme="dark"] *) .toast.warning .toast-icon { color: #fbbf24; }
:where([data-theme="dark"], [data-theme="dark"] *) .toast.info .toast-icon { color: #60a5fa; }

.toast-message {
  flex: 1;
  font-size: 0.875rem;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.toast-close {
  padding: 0.25rem;
  border-radius: var(--radius-sm);
  background: none;
  border: none;
  color: currentColor;
  opacity: 0.6;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.toast-close:hover {
  opacity: 1;
}

@media (min-width: 768px) {
  .toast-container {
    bottom: 2rem;
  }
}
</style>
