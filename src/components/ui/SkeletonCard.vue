<template>
  <div class="skeleton-card" :class="{ 'skeleton': loading }">
    <slot v-if="!loading" />
    <template v-else>
      <div class="skeleton skeleton-title" />
    <template v-else>
      <div class="skeleton skeleton-text skeleton-text-md" />
      <div class="skeleton skeleton-text skeleton-text-lg" />
      <div class="skeleton skeleton-text skeleton-text-sm" />
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
.skeleton-card {
  @apply rounded border shadow-sm;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}

.skeleton-text-md { @apply w-[80%]; }
.skeleton-text-lg { @apply w-[60%]; }
.skeleton-text-sm { @apply w-[40%]; }
</style>