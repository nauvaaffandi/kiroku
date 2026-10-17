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
.scheduled-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.scheduled-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.scheduled-info {
  flex: 1;
}

.scheduled-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 0.25rem;
}

.scheduled-description {
  font-size: 0.8125rem;
  color: var(--color-text-muted);
  margin: 0;
}

.scheduled-badges {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  flex-wrap: wrap;
}

.scheduled-type-badge {
  font-size: 0.75rem;
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
  background: var(--color-background);
  color: var(--color-text-muted);
  text-transform: capitalize;
}

.scheduled-status-badge {
  font-size: 0.75rem;
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
}

.scheduled-status-badge.active {
  background: var(--color-success-bg);
  color: var(--color-success);
}

.scheduled-status-badge.inactive {
  background: var(--color-warning-bg);
  color: var(--color-warning);
}

.scheduled-details {
  display: flex;
  gap: 1.5rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--color-border);
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.detail-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-value {
  font-size: 0.875rem;
  color: var(--color-text);
  font-weight: 500;
}

.detail-value.amount {
  font-weight: 600;
  font-size: 1rem;
}

.scheduled-actions {
  display: flex;
  gap: 0.5rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--color-border);
  margin-top: 0.5rem;
}
</style>