<template>
  <button
    :class="['btn', variant, size, { disabled: disabled, loading: loading }]"
    :disabled="disabled || loading"
    :type="type"
    @click="$emit('click', $event)"
  >
    <span v-if="loading" class="btn-spinner" aria-hidden="true"></span>
    <slot v-else />
  </button>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost' | 'link'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  type?: 'button' | 'submit' | 'reset'
}

withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false,
  loading: false,
  type: 'button'
})

defineEmits<{
  click: [event: MouseEvent]
}>()
</script>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  font-weight: 500;
  border-radius: var(--radius);
  border: none;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
  white-space: nowrap;
}

.btn:disabled,
.btn.loading {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--color-primary);
  color: var(--color-primary-contrast);
}

.btn-primary:hover:not(:disabled):not(.loading) {
  background: var(--color-primary-hover);
}

.btn-secondary {
  background: var(--color-surface);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover:not(:disabled):not(.loading) {
  background: var(--color-surface-hover);
}

.btn-danger {
  background: var(--color-danger);
  color: var(--color-danger-contrast);
}

.btn-danger:hover:not(:disabled):not(.loading) {
  background: var(--color-danger-hover);
}

.btn-ghost {
  background: transparent;
  color: var(--color-text);
}

.btn-ghost:hover:not(:disabled):not(.loading) {
  background: var(--color-surface);
}

.btn-link {
  background: transparent;
  color: var(--color-primary);
  padding: 0;
}

.btn-link:hover:not(:disabled):not(.loading) {
  text-decoration: underline;
}

.btn-sm {
  padding: 0.375rem 0.75rem;
  font-size: 0.8125rem;
  height: 32px;
}

.btn-md {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  height: 40px;
}

.btn-lg {
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  height: 48px;
}

.btn-spinner {
  width: 1em;
  height: 1em;
  border: 2px solid currentColor;
  border-right-color: transparent;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>