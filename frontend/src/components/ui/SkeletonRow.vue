<template>
  <div class="skeleton-row" :class="{ 'skeleton': loading }">
    <slot v-if="!loading" />
    <template v-else>
      <div class="skeleton skeleton-avatar" />
      <div class="skeleton-row-content">
        <div class="skeleton skeleton-text" style="width: 60%;" />
        <div class="skeleton skeleton-text" style="width: 40%;" />
      </div>
      <div class="skeleton skeleton-text" style="width: 80px;" />
    </template>
  </div>
</template>

<script setup lang="ts">
interface Props {
  loading: boolean
}

defineProps<Props>()
</script>

<style scoped>
@reference "tailwindcss";
.skeleton-row {
  @apply flex items-center gap-4 p-4 rounded border shadow-sm;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}

.skeleton-row-content {
  @apply flex-1 flex flex-col gap-2;
}
</style>