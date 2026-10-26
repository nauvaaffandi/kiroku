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
@reference "tailwindcss";

.analytics-card {
  @apply flex flex-col rounded-2xl p-5 gap-4;
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
}

.card-header {
  @apply flex justify-between items-center;
}

.card-title {
  @apply text-base font-semibold m-0;
  color: var(--text-primary);
}

.card-body {
  @apply flex flex-col flex-1;
}
</style>
