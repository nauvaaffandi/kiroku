<template>
  <div class="scheduled-list">
    <div v-if="scheduled.length === 0" class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="empty-icon">
        <circle cx="12" cy="12" r="10"></circle>
        <polyline points="12 6 12 12 16 14"></polyline>
      </svg>
      <p>No scheduled transactions yet</p>
      <BaseButton @click="$emit('create')">Add Scheduled Transaction</BaseButton>
    </div>

    <div v-else class="scheduled-grid">
      <ScheduledCard
        v-for="s in scheduled"
        :key="s.id"
        :scheduled="s"
        @edit="$emit('edit', $event)"
        @toggle="$emit('toggle', $event)"
        @delete="$emit('delete', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/ui'
import ScheduledCard from './ScheduledCard.vue'
import type { ScheduledTransaction } from '@/types'

defineProps<{
  scheduled: ScheduledTransaction[]
}>()

defineEmits<{
  edit: [scheduled: ScheduledTransaction]
  toggle: [scheduled: ScheduledTransaction]
  delete: [id: string]
  create: []
}>()
</script>

<style scoped>
@reference "tailwindcss";
.scheduled-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1rem;
}

.empty-state {
  text-align: center;
  padding: 3rem 1.5rem;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  color: var(--color-text-muted);
}

.empty-icon {
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state p {
  margin-bottom: 1rem;
}
</style>