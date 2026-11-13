<template>
  <div class="skeleton-chart" :class="{ 'skeleton': loading }" :style="{ height: height }">
    <slot v-if="!loading" />
  </div>
</template>

<script setup lang="ts">
interface Props {
  loading: boolean
  height?: string
}

defineProps<Props>()
</script>

<style scoped>
@reference "tailwindcss";
.skeleton-chart {
  @apply w-full rounded border;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}
</style>