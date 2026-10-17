<template>
  <div></div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  if (import.meta.env.VITE_USE_ERUDA === 'true') {
    const script = document.createElement('script')
    script.src = 'https://cdn.jsdelivr.net/npm/eruda'
    script.onload = () => {
      // @ts-ignore
      if (window.eruda) {
        // @ts-ignore
        window.eruda.init()
      }
    }
    document.body.appendChild(script)
  }
})
</script>