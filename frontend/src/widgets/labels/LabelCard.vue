<template>
  <div class="label-card">
    <div class="label-info">
      <span class="color-dot" :style="{ backgroundColor: label.color }"></span>
      <span class="label-icon" v-if="label.icon">{{ label.icon }}</span>
      <div class="label-details">
        <span class="label-name">{{ label.name }}</span>
        <span class="label-desc" v-if="label.description">{{ label.description }}</span>
      </div>
    </div>
    <div class="label-actions">
      <button @click="$emit('edit', label)" class="action-btn" aria-label="Edit label">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
        </svg>
      </button>
      <button @click="$emit('delete', label.id)" class="action-btn delete" aria-label="Delete label">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="3 6 5 6 21 6"></polyline>
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Label } from '@/types'

defineProps<{
  label: Label
}>()

defineEmits<{
  edit: [label: Label]
  delete: [id: string]
}>()
</script>

<style scoped>
@reference "tailwindcss";

.label-card {
  @apply flex items-center justify-between p-4 rounded-lg border mb-3 transition-colors;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}

.label-card:hover {
  background-color: var(--bg-hover);
}

.label-info {
  @apply flex items-center gap-3;
}

.color-dot {
  @apply w-4 h-4 rounded-full shrink-0;
}

.label-icon {
  @apply text-lg;
}

.label-details {
  @apply flex flex-col;
}

.label-name {
  @apply font-medium;
  color: var(--text-primary);
}

.label-desc {
  @apply text-xs;
  color: var(--text-muted);
}

.label-actions {
  @apply flex items-center gap-2;
}

.action-btn {
  @apply p-2 rounded border-0 bg-transparent cursor-pointer transition-colors;
  color: var(--text-muted);
}

.action-btn:hover {
  background-color: var(--bg-tertiary);
  color: var(--text-primary);
}

.action-btn.delete:hover {
  color: var(--error);
}
</style>
