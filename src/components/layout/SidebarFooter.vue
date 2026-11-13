<template>
  <div class="sidebar-footer" />
</template>

<script setup lang="ts">
defineEmits(['navigate'])
</script>

<style scoped>
@reference "tailwindcss";

.sidebar-footer {
  @apply px-6 py-4 border-t;
  border-color: var(--border-primary);
}

.footer-link {
  @apply flex items-center gap-3 px-4 py-3 rounded-sm no-underline transition-colors;
  color: var(--text-secondary);
}

.footer-link:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.footer-link svg {
  flex-shrink: 0;
}
</style>
