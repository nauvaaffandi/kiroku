<template>
  <Teleport to="body">
    <div v-if="visible" class="modal-overlay" @click.self="emit('close')">
      <div class="modal-container">
        <div class="modal-header">
          <h2 class="modal-title">{{ t('goals.addGoal') }}</h2>
          <button @click="emit('close')" class="modal-close" aria-label="Close">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
        <form @submit.prevent="submit" class="modal-body">
          <div class="form-field">
            <label class="form-label">{{ t('goals.name') }}</label>
            <input v-model="form.name" class="form-input" placeholder="" />
          </div>
          <div class="form-field">
            <label class="form-label">{{ t('goals.targetAmount') }}</label>
            <input v-model.number="form.target_amount" type="number" class="form-input" placeholder="" />
          </div>
          <div class="form-field">
            <label class="form-label">{{ t('goals.currency') }}</label>
            <input v-model="form.currency" class="form-input" placeholder="USD" />
          </div>
          <div class="form-field">
            <label class="form-label">{{ t('goals.desiredDate') }}</label>
            <input v-model="form.desired_date" type="date" class="form-input" />
          </div>
          <div class="form-field">
            <label class="form-label">{{ t('goals.colorTheme') }}</label>
            <input v-model="form.color_theme" type="color" class="form-color" />
          </div>
          <div class="form-actions">
            <button type="button" class="btn-secondary" @click="emit('close')">{{ t('goals.cancel') }}</button>
            <button type="submit" class="btn-primary" :disabled="submitting">{{ submitting ? t('goals.saving') : t('goals.create') }}</button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { api } from '@/api/client'
import { useUIStore } from '@/stores/ui'
import { useI18n } from '@/locales'
import type { CreateGoalInput } from '@/types/goal'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  close: []
  created: []
}>()

const { t } = useI18n()
const uiStore = useUIStore()
const submitting = ref(false)

const form = ref<CreateGoalInput>({
  name: '',
  target_amount: 0,
  currency: 'USD',
  icon: '',
  desired_date: '',
  color_theme: '',
  note: ''
})

async function submit() {
  submitting.value = true
  try {
    await api.post('/api/v1/goals/', form.value)
    uiStore.success(t('goals.createSuccess'))
    emit('created')
  } catch (e) {
    uiStore.error(t('goals.fetchError'))
  } finally {
    submitting.value = false
  }
}

watch(() => props.visible, (val) => {
  if (val) {
    form.value = { name: '', target_amount: 0, currency: 'USD', icon: '', desired_date: '', color_theme: '', note: '' }
  }
})
</script>

<style scoped>
@reference "tailwindcss";

.modal-overlay {
  @apply fixed inset-0 flex items-center justify-center p-4 z-[1000];
  background-color: rgba(0, 0, 0, 0.5);
}
.modal-container {
  @apply border rounded-xl w-full max-w-[400px] max-h-[90vh] overflow-hidden;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}
.modal-header {
  @apply flex items-center justify-between py-4 px-6 border-b;
  border-color: var(--border-primary);
}
.modal-title {
  @apply text-lg font-semibold;
  color: var(--text-primary);
}
.modal-close {
  @apply p-2 rounded-sm bg-transparent border-none;
  color: var(--text-secondary);
}
.modal-body {
  @apply p-6;
}
.form-field {
  @apply mb-4;
}
.form-label {
  @apply block text-xs uppercase tracking-wider mb-1.5;
  color: var(--text-secondary);
}
.form-input {
  @apply w-full px-3 py-2 rounded-sm border;
  border-color: var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
}
.form-color {
  @apply w-full h-8 rounded-sm border-0 cursor-pointer bg-transparent;
}
.form-actions {
  @apply flex justify-end gap-2 mt-6;
}
.btn-primary {
  @apply inline-flex items-center gap-2 px-4 py-2 rounded-sm font-medium text-sm no-underline;
  background-color: var(--accent);
  color: var(--text-inverse);
  border: none;
  cursor: pointer;
  transition: opacity 0.2s ease;
}
.btn-primary:disabled {
  @apply opacity-50 cursor-not-allowed;
}
.btn-secondary {
  @apply inline-flex items-center gap-2 px-4 py-2 rounded-sm font-medium text-sm no-underline;
  background-color: var(--bg-hover);
  color: var(--text-primary);
  border: 1px solid var(--border-primary);
  cursor: pointer;
  transition: opacity 0.2s ease;
}
</style>