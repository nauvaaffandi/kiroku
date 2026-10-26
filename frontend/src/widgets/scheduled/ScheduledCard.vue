<template>
  <div class="scheduled-card">
    <div class="scheduled-header">
      <div class="scheduled-info">
        <h3 class="scheduled-name">{{ scheduled.name }}</h3>
        <p class="scheduled-description">{{ scheduled.description }}</p>
      </div>
      <div class="scheduled-badges">
        <span class="scheduled-type-badge" :class="scheduled.type">{{ scheduled.type }}</span>
        <span v-if="scheduled.active" class="scheduled-status-badge active">Active</span>
        <span v-else class="scheduled-status-badge inactive">Inactive</span>
      </div>
    </div>
    <div class="scheduled-details">
      <div class="detail-item">
        <span class="detail-label">Frequency</span>
        <span class="detail-value">{{ formatFrequency(scheduled.frequency) }}</span>
      </div>
      <div class="detail-item">
        <span class="detail-label">Amount</span>
        <span class="detail-value amount">{{ formatCurrency(scheduled.amount) }} {{ scheduled.currency }}</span>
      </div>
      <div class="detail-item" v-if="scheduled.next_run">
        <span class="detail-label">Next Run</span>
        <span class="detail-value">{{ formatDate(scheduled.next_run) }}</span>
      </div>
    </div>
    <div class="scheduled-actions">
      <BaseButton variant="secondary" size="sm" @click="$emit('edit', scheduled)">
        Edit
      </BaseButton>
      <BaseButton variant="secondary" size="sm" @click="$emit('toggle', scheduled)">
        {{ scheduled.active ? 'Pause' : 'Resume' }}
      </BaseButton>
      <BaseButton variant="danger" size="sm" @click="$emit('delete', scheduled.id)">
        Delete
      </BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/ui'
import type { ScheduledTransaction } from '@/types'

defineProps<{
  scheduled: ScheduledTransaction
}>()

defineEmits<{
  edit: [scheduled: ScheduledTransaction]
  toggle: [scheduled: ScheduledTransaction]
  delete: [id: string]
}>()

const formatCurrency = (num: number) => {
  return new Intl.NumberFormat('ja-JP').format(num)
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('ja-JP')
}

const formatFrequency = (freq: string) => {
  const map: Record<string, string> = {
    daily: 'Daily',
    weekly: 'Weekly',
    monthly: 'Monthly',
    yearly: 'Yearly',
    custom: 'Custom'
  }
  return map[freq] || freq
}
</script>

<style scoped>
@reference "tailwindcss";
.scheduled-card {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  @apply p-5 flex flex-col gap-4;
}

.scheduled-header {
  @apply flex justify-between items-start;
}

.scheduled-info {
  @apply flex-1;
}

.scheduled-name {
  color: var(--text-primary);
  @apply text-base font-semibold mb-1 mt-0;
}

.scheduled-description {
  color: var(--text-muted);
  @apply text-[0.8125rem] m-0;
}

.scheduled-badges {
  @apply flex gap-2 items-center flex-wrap;
}

.scheduled-type-badge {
  background-color: var(--bg-primary);
  color: var(--text-muted);
  @apply text-xs px-2 py-0.5 rounded-full capitalize;
}

.scheduled-status-badge {
  @apply text-xs px-2 py-0.5 rounded-full;
}

.scheduled-status-badge.active {
  background-color: var(--success-light);
  color: var(--success);
}

.scheduled-status-badge.inactive {
  background-color: var(--warning-light);
  color: var(--warning);
}

.scheduled-details { border-top: 1px solid var(--border-primary); @apply flex gap-6 pt-3; }

.detail-item {
  @apply flex flex-col gap-1;
}

.detail-label {
  color: var(--text-muted);
  @apply text-xs uppercase tracking-wider;
}

.detail-value {
  color: var(--text-primary);
  @apply text-sm font-medium;
}

.detail-value.amount {
  @apply font-semibold text-base;
}

.scheduled-actions { border-top: 1px solid var(--border-primary); @apply flex gap-2 pt-3 mt-2; }
</style>