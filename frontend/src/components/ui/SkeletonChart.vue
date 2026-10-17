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
.skeleton-chart {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  width: 100%;
}
</style>