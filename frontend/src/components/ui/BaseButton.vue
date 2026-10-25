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
@reference "tailwindcss";

.btn {
  @apply inline-flex items-center justify-center gap-2 font-medium rounded border-0 cursor-pointer transition-all font-sans whitespace-nowrap;
  font-family: inherit;
}

.btn:disabled,
.btn.loading {
  @apply opacity-60 cursor-not-allowed;
}

.btn-primary {
  background: var(--accent);
  color: var(--text-inverse);
}

.btn-primary:hover:not(:disabled):not(.loading) {
  background: var(--accent-hover);
}

.btn-secondary {
  background: var(--bg-card);
  color: var(--text-primary);
  border: 1px solid var(--border-primary);
}

.btn-secondary:hover:not(:disabled):not(.loading) {
  background: var(--bg-hover);
}

.btn-danger {
  background: var(--error);
  color: var(--text-inverse);
}

.btn-danger:hover:not(:disabled):not(.loading) {
  background: var(--error);
}

.btn-ghost {
  background: transparent;
  color: var(--text-primary);
}

.btn-ghost:hover:not(:disabled):not(.loading) {
  background: var(--bg-card);
}

.btn-link {
  background: transparent;
  color: var(--accent);
  @apply p-0;
}

.btn-link:hover:not(:disabled):not(.loading) {
  text-decoration: underline;
}

.btn-sm {
  @apply px-3 py-1 text-sm h-[32px];
}

.btn-md {
  @apply px-4 py-2 text-sm h-[40px];
}

.btn-lg {
  @apply px-6 py-3 text-base h-[48px];
}

.btn-spinner {
  width: 1em;
  height: 1em;
  border: 2px solid currentColor;
  border-right-color: transparent;
  @apply rounded-full;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>