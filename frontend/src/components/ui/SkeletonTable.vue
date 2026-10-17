<template>
  <div class="skeleton-table" :class="{ 'skeleton': loading }">
    <slot v-if="!loading" />
    <template v-else>
      <div class="skeleton-table-header">
        <div class="skeleton skeleton-text" style="width: 100px;" />
        <div class="skeleton skeleton-text" style="width: 80px;" />
        <div class="skeleton skeleton-text" style="width: 80px;" />
        <div class="skeleton skeleton-text" style="width: 60px;" />
      </div>
      <div class="skeleton-table-body">
        <div v-for="i in rows" :key="i" class="skeleton-table-row">
          <div class="skeleton skeleton-text" />
          <div class="skeleton skeleton-text" />
          <div class="skeleton skeleton-text" />
          <div class="skeleton skeleton-text" />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  loading: boolean
  rows?: number
  cols?: number
}

const props = defineProps<Props>()

const rows = computed(() => props.rows || 5)
</script>

<style scoped>
.skeleton-table {
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  overflow: hidden;
}

.skeleton-table-header {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 1rem;
  padding: 1rem;
  background-color: var(--bg-secondary);
  border-bottom: 1px solid var(--border-primary);
}

.skeleton-table-body {
  display: flex;
  flex-direction: column;
}

.skeleton-table-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 1rem;
  padding: 1rem;
  border-bottom: 1px solid var(--border-primary);
}

.skeleton-table-row:last-child {
  border-bottom: none;
}
</style>