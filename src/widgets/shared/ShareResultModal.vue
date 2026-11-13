<template>
  <Transition name="modal">
    <div v-if="share" class="modal-overlay" @click.self="emit('close')">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>{{ t('shared.access.createdSuccess') }}</h2>
        </div>
        <div class="modal-form">
          <div class="share-result">
            <p>{{ t('shared.access.createdSuccess') }}. {{ t('shared.access.warning') }}</p>
            <div class="share-url-display">
              <input :value="fullShareUrl" readonly class="form-input" />
              <button @click="copyFullShareUrl" class="btn-secondary">{{ t('shared.access.copy') }}</button>
            </div>
            <p class="share-warning">{{ t('shared.access.warning') }}</p>
          </div>
          <div class="modal-actions">
            <button @click="emit('close')" class="btn-primary">{{ t('shared.access.done') }}</button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import type { Share } from './useShare'

const props = defineProps<{
  share: Share | null
  fullShareUrl: string
  copyFullShareUrl: () => void
}>()

const emit = defineEmits<{
  close: []
}>()

const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";

.modal-overlay {
  @apply fixed inset-0 flex items-center justify-center p-4 z-[1000];
  background-color: rgba(0, 0, 0, 0.5);
  animation: fadeIn 0.2s ease;
}

.modal {
  @apply rounded-[var(--radius)] w-full max-w-[480px] max-h-[90vh] overflow-y-auto shadow-[var(--shadow-md)];
  background-color: var(--bg-card);
  animation: slideUp 0.2s ease;
}

.modal-header {
  @apply flex items-center justify-between py-4 px-6;
  border-bottom: 1px solid var(--border-primary);
}

.modal-header h2 {
  @apply text-lg font-semibold;
  color: var(--text-primary);
}

.modal-form {
  @apply p-6;
}

.share-result p {
  @apply mb-3;
  color: var(--text-secondary);
}

.share-url-display {
  @apply flex gap-2 mb-3;
}

.share-url-display .form-input {
  @apply flex-1;
  background-color: var(--bg-tertiary);
}

.form-input {
  @apply w-full py-2 px-3 text-sm rounded-[var(--radius-sm)] transition-all duration-200;
  border: 1px solid var(--border-primary);
  background-color: var(--bg-primary);
  color: var(--text-primary);
}

.form-input:focus {
  @apply outline-none border-transparent ring-2 ring-[var(--accent)];
}

.share-warning {
  @apply text-[0.8125rem] m-0;
  color: var(--warning);
}

.modal-actions {
  @apply flex justify-end gap-3 mt-6;
}

.btn-primary {
  @apply py-2 px-4 rounded-[var(--radius-sm)] font-medium text-sm cursor-pointer transition-opacity duration-200 inline-flex items-center justify-center gap-2 border-none text-white;
  background-color: var(--accent);
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--accent-hover);
}

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
</style>