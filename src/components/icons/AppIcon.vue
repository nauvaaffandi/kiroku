<template>
  <img
    class="app-icon"
    src="/legacy_512.png"
    alt="Kiroku logo"
    aria-hidden="true"
  />
</template>

<script setup lang="ts">
// AppIcon - Kiroku (記録) brand logo
// Uses raster logo from public/legacy_512.png
</script>

<style scoped>
@reference "tailwindcss";

.app-icon {
  @apply inline-block h-full w-full shrink-0;
  width: 1em;
  height: 1em;
}

/* Size variants via class */
.app-icon--sm {
  width: 1rem;
  height: 1rem;
}

.app-icon--md {
  width: 1.5rem;
  height: 1.5rem;
}

.app-icon--lg {
  width: 2rem;
  height: 2rem;
}

.app-icon--xl {
  width: 3rem;
  height: 3rem;
}
</style>