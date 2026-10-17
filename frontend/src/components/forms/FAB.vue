<template>
  <button @click="handleClick" class="fab" aria-label="Add transaction">
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
      <line x1="12" y1="5" x2="12" y2="19"></line>
      <line x1="5" y1="12" x2="19" y2="12"></line>
    </svg>
  </button>
</template>

<script setup lang="ts">
const emit = defineEmits(['click'])

const handleClick = (e: MouseEvent) => {
  if (navigator.vibrate) {
    navigator.vibrate(10)
  }
  emit('click', e)
}
</script>

<style scoped>
.fab {
  position: fixed;
  bottom: calc(72px + env(safe-area-inset-bottom));
  right: 1.5rem;
  width: 3.5rem;
  height: 3.5rem;
  border-radius: 50%;
  background-color: var(--color-primary);
  color: white;
  border: none;
  box-shadow: var(--shadow-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 90;
  transition: transform 0.2s ease, background-color 0.2s ease;
}

.fab:hover {
  background-color: var(--color-primary-hover);
  transform: scale(1.05);
}

.fab:active {
  transform: scale(0.95);
}

@media (min-width: 768px) {
  .fab {
    bottom: 2rem;
  }
}
</style>