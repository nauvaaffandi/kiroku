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
.skeleton-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
}

.skeleton-row-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
</style>