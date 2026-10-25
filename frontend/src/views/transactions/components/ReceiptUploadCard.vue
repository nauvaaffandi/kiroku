<template>
  <div class="upload-card">
    <div class="dropzone" @click="triggerFileInput" @dragover.prevent @drop.prevent="handleDrop">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
        <polyline points="17 8 12 3 7 8"></polyline>
        <line x1="12" y1="3" x2="12" y2="15"></line>
      </svg>
      <h3>Click or drop receipt image here</h3>
      <p>Supports PNG, JPEG. Images > 2MB will be compressed locally preserving original dimensions.</p>
      <input ref="fileInput" type="file" accept="image/*" class="hidden-input" @change="handleFileSelect" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const fileInput = ref<HTMLInputElement | null>(null)

const emit = defineEmits(['file-selected'])

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleDrop = (e: DragEvent) => {
  const file = e.dataTransfer?.files?.[0]
  if (file) emit('file-selected', file)
}

const handleFileSelect = (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) emit('file-selected', file)
}
</script>

<style scoped>
@reference "tailwindcss";
.upload-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  padding: 2rem;
  box-shadow: var(--shadow);
}

.dropzone {
  border: 2px dashed var(--color-border);
  border-radius: var(--radius);
  padding: 3rem 2rem;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s ease, background-color 0.2s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.dropzone:hover {
  border-color: var(--color-primary);
  background-color: var(--color-neutral-100);
}

.hidden-input {
  display: none;
}
</style>