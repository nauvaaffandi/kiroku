<template>
  <Transition name="modal">
    <div v-if="showCreateModal" class="modal-overlay" @click.self="$emit('close')">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>{{ t('shared.access.createShare') }}</h2>
          <button @click="emit('close')" class="modal-close" :aria-label="t('shared.access.cancel')">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
        <form @submit.prevent="emit('submit', form)" class="modal-form">
          <div class="form-group">
            <label for="account_id">{{ t('shared.access.account') }} *</label>
            <select id="account_id" v-model="form.account_id" required class="form-input">
              <option value="">{{ t('shared.access.select') }}</option>
              <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
                {{ acc.name }} ({{ acc.currency }})
              </option>
            </select>
          </div>
          <div class="form-group">
            <label for="expires_at">{{ t('shared.access.expiresAt') }} (optional)</label>
            <input
              id="expires_at"
              v-model="form.expires_at"
              type="datetime-local"
              class="form-input"
            />
            <p class="form-hint">{{ t('shared.access.expiresHint') }}</p>
          </div>
          <div class="modal-actions">
            <button type="button" @click="emit('close')" class="btn-secondary">{{ t('shared.access.cancel') }}</button>
            <button type="submit" :disabled="loading" class="btn-primary">
              <span v-if="loading">{{ t('shared.access.loading') }}</span>
              <span v-else>{{ t('shared.access.createShare') }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from '@/locales'
import type { Account } from './useShare'

const props = defineProps<{
  accounts: Account[]
  showCreateModal: boolean
  loading: boolean
}>()

const emit = defineEmits<{
  close: []
  submit: [form: { account_id: string; expires_at: string }]
}>()

const { t } = useI18n()

const form = ref<{ account_id: string; expires_at: string }>({ account_id: '', expires_at: '' })
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

.modal-close {
  @apply p-1 rounded-[var(--radius-sm)] bg-transparent border-none cursor-pointer transition-colors duration-200;
  color: var(--text-muted);
}

.modal-close:hover {
  @apply text-[var(--text-primary)];
  background-color: var(--bg-hover);
}

.modal-form {
  @apply p-6;
}

.form-group {
  @apply mb-4;
}

.form-group label {
  @apply block text-sm font-medium mb-1.5;
  color: var(--text-secondary);
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

.form-hint {
  @apply text-xs mt-1;
  color: var(--text-muted);
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

.btn-secondary {
  @apply py-2 px-4 rounded-[var(--radius-sm)] font-medium text-sm cursor-pointer transition-opacity duration-200 inline-flex items-center justify-center gap-2 border-none text-[var(--text-secondary)];
  background-color: var(--bg-tertiary);
}

.btn-secondary:hover:not(:disabled) {
  background-color: var(--bg-hover);
}

.btn-primary:disabled, .btn-secondary:disabled {
  @apply opacity-50 cursor-not-allowed;
}

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
</style>