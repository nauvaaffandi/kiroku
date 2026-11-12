<template>
  <div class="skeleton-table" :class="{ 'skeleton': loading }">
    <slot v-if="!loading" />
    <template v-else>
      <div class="skeleton-table-header">
        <div class="skeleton skeleton-text skeleton-text-100px" />
        <div class="skeleton skeleton-text skeleton-text-80px" />
        <div class="skeleton skeleton-text skeleton-text-80px" />
        <div class="skeleton skeleton-text skeleton-text-60px" />
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
@reference "tailwindcss";
.skeleton-table {
  @apply rounded border overflow-hidden shadow-sm;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}

.skeleton-table-header {
  @apply grid grid-cols-[2fr_1fr_1fr_1fr] gap-4 p-4 border-b;
  background-color: var(--bg-secondary);
  border-color: var(--border-primary);
}

.skeleton-table-row {
  @apply grid grid-cols-[2fr_1fr_1fr_1fr] gap-4 p-4 border-b;
  border-color: var(--border-primary);
}

.skeleton-text-100px { @apply w-[100px]; }
.skeleton-text-80px { @apply w-[80px]; }
.skeleton-text-60px { @apply w-[60px]; }
</style>