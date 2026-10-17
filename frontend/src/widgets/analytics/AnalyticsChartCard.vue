<template>
  <div class="analytics-card">
    <div class="card-header" v-if="title || $slots.header">
      <h3 class="card-title">{{ title }}</h3>
      <div class="card-actions">
        <slot name="header-actions" />
      </div>
    </div>
    <div class="card-body">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
}>()
</script>

<style scoped>
.analytics-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.card-body {
  display: flex;
  flex-direction: column;
  flex: 1;
}
</style>