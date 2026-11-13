<template>
  <CheckCircle2 class="lucide-icon" />
</template>
<script setup lang="ts">
import { CheckCircle2 } from '@lucide/vue'
</script>
<style scoped>
@reference "tailwindcss";
.lucide-icon {
  @apply inline-block h-[1.25em] w-[1.25em] shrink-0 align-middle text-current;
}
</style>
