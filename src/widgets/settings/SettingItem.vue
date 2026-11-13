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
@reference "tailwindcss";
.setting-item {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  @apply flex justify-between items-center p-4 gap-4;
}

.setting-info {
  @apply flex-1 min-w-0;
}

.setting-label {
  color: var(--text-primary);
  @apply block font-medium mb-1;
}

.setting-desc {
  color: var(--text-muted);
  @apply text-[0.8125rem];
}

@media (max-width: 600px) {
  .setting-item {
    @apply flex-col items-start;
  }
}
</style>