<template>
  <div class="upload-card">
    <div class="dropzone" @click="triggerFileInput" @dragover.prevent @drop.prevent="handleDrop" role="button" tabindex="0" :aria-label="t('receipt.dropzoneTitle')">
      <Upload :size="48" stroke-width="1.75" aria-hidden="true" />
      <h3>{{ t('receipt.dropzoneTitle') }}</h3>
      <p>{{ t('receipt.dropzoneDesc') }}</p>
      <input ref="fileInput" type="file" accept="image/*" class="hidden-input" @change="handleFileSelect" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from '@/locales'
import { Upload } from '@lucide/vue'

const fileInput = ref<HTMLInputElement | null>(null)

const emit = defineEmits(['file-selected'])

const { t } = useI18n()

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
  @apply bg-[var(--bg-card)] border border-[var(--border-primary)] rounded-[14px] p-6;
}

.dropzone {
  @apply border-2 border-dashed border-[var(--border-primary)] rounded-[10px] p-8 text-center cursor-pointer transition-colors duration-120 ease-out flex flex-col items-center gap-4;
}

.dropzone:hover {
  @apply border-[var(--accent)] bg-[var(--bg-tertiary)];
}

.dropzone:focus-visible {
  @apply ring-2 ring-offset-2 rounded-[10px];
  ring-offset-color: var(--bg-card);
  ring-color: var(--accent);
}

.hidden-input {
  @apply hidden;
}
</style>
