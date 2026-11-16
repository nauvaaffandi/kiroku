<template>
  <div class="empty-state">
    <p class="empty-state-message">{{ message }}</p>
    <button
      v-if="action"
      type="button"
      class="empty-state-action"
      @click="$emit('action')"
    >
      {{ action }}
    </button>
  </div>
</template>

<script setup lang="ts">
interface Props {
  message: string
  action?: string
}

const props = defineProps<Props>()

defineEmits<{
  action: []
}>()
</script>

<style scoped>
@reference "tailwindcss";

.empty-state {
  @apply flex flex-col items-center gap-3 py-12 px-6 text-center;
}

.empty-state-message {
  @apply text-base text-[var(--text-secondary)];
}

.empty-state-action {
  @apply inline-flex items-center justify-center gap-2 px-4 py-2 rounded-[10px] font-medium text-sm transition-colors duration-120 ease-out;
  background-color: var(--accent);
  color: var(--text-inverse);
  border: none;
  cursor: pointer;
  min-height: 44px;
  min-width: 44px;
}

.empty-state-action:hover {
  background-color: var(--accent-hover);
}

.empty-state-action:focus-visible {
  @apply ring-2 ring-offset-2 rounded-[10px];
  ring-offset-color: var(--bg-card);
  ring-color: var(--accent);
}
</style>
