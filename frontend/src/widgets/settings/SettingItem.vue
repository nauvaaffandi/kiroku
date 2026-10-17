<template>
  <div class="setting-item">
    <div class="setting-info">
      <span class="setting-label">{{ label }}</span>
      <span class="setting-desc" v-if="description">{{ description }}</span>
    </div>
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  description?: string
}>()
</script>

<style scoped>
.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  gap: 1rem;
}

.setting-info {
  flex: 1;
  min-width: 0;
}

.setting-label {
  display: block;
  font-weight: 500;
  color: var(--color-text);
  margin-bottom: 0.25rem;
}

.setting-desc {
  font-size: 0.8125rem;
  color: var(--color-text-muted);
}

@media (max-width: 600px) {
  .setting-item {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>