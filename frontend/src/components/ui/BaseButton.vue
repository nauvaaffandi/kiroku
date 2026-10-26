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
  background-color: var(--accent);
  color: var(--text-inverse);
}

.btn-primary:hover:not(:disabled):not(.loading) {
  background-color: var(--accent-hover);
}

.btn-secondary {
  @apply border;
  background-color: var(--bg-card);
  color: var(--text-primary);
  border-color: var(--border-primary);
}

.btn-secondary:hover:not(:disabled):not(.loading) {
  background-color: var(--bg-hover);
}

.btn-danger {
  background-color: var(--error);
  color: var(--text-inverse);
}

.btn-danger:hover:not(:disabled):not(.loading) {
  background-color: var(--error);
}

.btn-ghost {
  @apply bg-transparent;
  color: var(--text-primary);
}

.btn-ghost:hover:not(:disabled):not(.loading) {
  background-color: var(--bg-card);
}

.btn-link {
  @apply bg-transparent p-0;
  color: var(--accent);
}

.btn-link:hover:not(:disabled):not(.loading) {
  @apply underline;
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
  @apply w-[1em] h-[1em] border-2 border-current border-r-transparent rounded-full;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>