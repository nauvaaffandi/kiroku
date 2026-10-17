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
.skeleton-card {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}
</style>