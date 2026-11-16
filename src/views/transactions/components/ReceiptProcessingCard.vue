<template>
  <div class="processing-card">
    <div class="metrics-grid">
      <div class="metric-item">
        <span class="label">{{ t('receipt.originalSize') }}</span>
        <span class="value">{{ formatSize(originalSize) }}</span>
      </div>
      <div class="metric-item">
        <span class="label">{{ t('receipt.compressedSize') }}</span>
        <span class="value">{{ formatSize(compressedSize) }}</span>
      </div>
    </div>

    <div class="events-timeline">
      <div v-for="ev in events" :key="ev.id" :class="['timeline-item', ev.status]">
        <div class="timeline-dot" aria-hidden="true"></div>
        <div class="timeline-content">
          <span class="timeline-message">{{ ev.message || ev.event_type }}</span>
          <span class="timeline-time">{{ new Date(ev.created_at).toLocaleTimeString() }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'

interface EventItem {
  id: string
  status?: string
  message?: string
  event_type?: string
  created_at: string
}

defineProps<{
  originalSize: number
  compressedSize: number
  events: EventItem[]
}>()

const { t } = useI18n()

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}
</script>

<style scoped>
@reference "tailwindcss";

.processing-card {
  @apply bg-[var(--bg-card)] border border-[var(--border-primary)] rounded-[14px] p-6;
}

.metrics-grid {
  @apply grid grid-cols-2 gap-4 my-6;
}

.metric-item {
  @apply bg-[var(--bg-tertiary)] p-4 rounded-[6px] flex flex-col gap-1;
}

.metric-item .label {
  @apply text-[0.75rem] text-[var(--text-muted)];
}

.metric-item .value {
  @apply text-[1.125rem] font-semibold text-[var(--text-primary)];
}

.events-timeline {
  @apply flex flex-col gap-3 mt-6 border-l border-[var(--border-primary)] pl-4;
}

.timeline-item {
  @apply flex items-center justify-between text-sm;
}

.timeline-item.success { @apply text-[var(--success)]; }
.timeline-item.warning { @apply text-[var(--warning)]; }
.timeline-item.error { @apply text-[var(--error)]; }

.timeline-time {
  @apply text-[0.75rem] text-[var(--text-muted)];
}
</style>
