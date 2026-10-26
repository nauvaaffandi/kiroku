<template>
  <section class="settings-section">
    <h2 class="section-title">{{ title }}</h2>
    <slot />
  </section>
</template>

<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<style scoped>
@reference "tailwindcss";
.settings-section {
  @apply mb-8;
}

.section-title {
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-primary);
  @apply text-lg font-semibold pb-2 mb-4 mt-0;
}
</style>