<template>
  <div class="skeleton-card" :class="{ 'skeleton': loading }">
    <slot v-if="!loading" />
    <template v-else>
      <div class="skeleton skeleton-title" />
      <div class="skeleton skeleton-text" style="width: 80%;" />
      <div class="skeleton skeleton-text" style="width: 60%;" />
      <div class="skeleton skeleton-text" style="width: 40%;" />
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
</style>