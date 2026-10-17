<template>
  <div class="processing-card">
    <div class="metrics-grid">
      <div class="metric-item">
        <span class="label">Original Size</span>
        <span class="value">{{ formatSize(originalSize) }}</span>
      </div>
      <div class="metric-item">
        <span class="label">Compressed Size</span>
        <span class="value">{{ formatSize(compressedSize) }}</span>
      </div>
    </div>

    <div class="events-timeline">
      <div v-for="ev in events" :key="ev.id" :class="['timeline-item', ev.status]">
        <div class="timeline-dot"></div>
        <div class="timeline-content">
          <span class="timeline-message">{{ ev.message || ev.event_type }}</span>
          <span class="timeline-time">{{ new Date(ev.created_at).toLocaleTimeString() }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
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

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}
</script>

<style scoped>
.processing-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  padding: 2rem;
  box-shadow: var(--shadow);
}

.metrics-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin: 1.5rem 0;
}

.metric-item {
  background: var(--color-neutral-100);
  padding: 1rem;
  border-radius: var(--radius-sm);
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.metric-item .label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.metric-item .value {
  font-size: 1.125rem;
  font-weight: 600;
}

.events-timeline {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-top: 1.5rem;
  border-left: 2px solid var(--color-border);
  padding-left: 1rem;
}

.timeline-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 0.875rem;
}

.timeline-item.success { color: #16a34a; }
.timeline-item.warning { color: #f59e0b; }
.timeline-item.error { color: #dc2626; }

.timeline-time {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}
</style>