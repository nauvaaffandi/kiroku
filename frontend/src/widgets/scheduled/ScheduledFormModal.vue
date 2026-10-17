<template>
  <BaseModal
    v-model="isOpen"
    :title="isEditing ? 'Edit Scheduled Transaction' : 'New Scheduled Transaction'"
    size="lg"
    @update:modelValue="$emit('update:modelValue', $event)"
  >
    <form @submit.prevent="handleSubmit" class="scheduled-form">
      <div class="form-group">
        <label class="form-label">Name *</label>
        <BaseInput
          v-model="formData.name"
          type="text"
          required
          placeholder="e.g., Monthly Rent"
        />
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">Type *</label>
          <BaseSelect
            v-model="formData.type"
            :options="typeOptions"
            required
          />
        </div>
        <div class="form-group">
          <label class="form-label">Frequency *</label>
          <BaseSelect
            v-model="formData.frequency"
            :options="frequencyOptions"
            required
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">Start Date *</label>
          <BaseInput
            v-model="formData.start_date"
            type="date"
            required
          />
        </div>
        <div class="form-group" v-if="formData.frequency === 'custom'">
          <label class="form-label">Cron Expression</label>
          <BaseInput
            v-model="formData.cron_expression"
            type="text"
            placeholder="e.g., 0 9 1 * *"
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">End Date</label>
          <BaseInput
            v-model="formData.end_date"
            type="date"
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">Amount *</label>
          <BaseInput
            v-model.number="formData.amount"
            type="number"
            required
            :min="0"
            :step="0.01"
            placeholder="0.00"
          />
        </div>
        <div class="form-group">
          <label class="form-label">Currency *</label>
          <BaseInput
            v-model="formData.currency"
            type="text"
            required
            maxlength="3"
            placeholder="JPY"
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">From Account *</label>
          <BaseSelect
            v-model="formData.from_account_id"
            :options="accountOptions"
            required
          />
        </div>
        <div class="form-group" v-if="formData.type === 'transfer'">
          <label class="form-label">To Account *</label>
          <BaseSelect
            v-model="formData.to_account_id"
            :options="accountOptions"
            required
          />
        </div>
        <div class="form-group" v-else>
          <label class="form-label">Category</label>
          <BaseSelect
            v-model="formData.category_id"
            :options="categoryOptions"
            placeholder="Select category (optional)"
          />
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">Description</label>
        <BaseInput
          v-model="formData.description"
          type="text"
          placeholder="Optional description"
        />
      </div>
      <div class="form-group form-checkbox">
        <label class="form-label">Active</label>
        <input
          v-model="formData.active"
          type="checkbox"
          class="form-checkbox-input"
        />
      </div>
      <div class="modal-actions">
        <BaseButton type="button" variant="secondary" @click="isOpen = false">
          Cancel
        </BaseButton>
        <BaseButton type="submit">
          {{ isEditing ? 'Save Changes' : 'Create Scheduled Transaction' }}
        </BaseButton>
      </div>
    </form>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { BaseModal, BaseInput, BaseSelect, BaseButton } from '@/components/ui'
import type { ScheduledTransaction, ScheduledType, ScheduledFrequency } from '@/types'

interface Props {
  modelValue: boolean
  scheduled?: ScheduledTransaction | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: [data: Partial<ScheduledTransaction>]
}>()

const isOpen = ref(props.modelValue)
const isEditing = computed(() => !!props.scheduled)

interface FormData {
  name: string
  type: ScheduledType
  frequency: ScheduledFrequency
  start_date: string
  end_date: string
  cron_expression: string
  amount: string
  currency: string
  from_account_id: string
  to_account_id: string
  category_id: string
  description: string
  active: boolean
}

const defaultForm = (): FormData => ({
  name: '',
  type: 'expense',
  frequency: 'monthly',
  start_date: new Date().toISOString().split('T')[0],
  end_date: '',
  cron_expression: '',
  amount: '0',
  currency: 'JPY',
  from_account_id: '',
  to_account_id: '',
  category_id: '',
  description: '',
  active: true
})

const formData = ref<FormData>(defaultForm())

const typeOptions = [
  { value: 'income', label: 'Income' },
  { value: 'expense', label: 'Expense' },
  { value: 'transfer', label: 'Transfer' }
]

const frequencyOptions = [
  { value: 'daily', label: 'Daily' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'monthly', label: 'Monthly' },
  { value: 'yearly', label: 'Yearly' },
  { value: 'custom', label: 'Custom (Cron)' }
]

const accountOptions = computed(() => (props.scheduled?.from_account_id
  ? [] // Would need accounts passed as prop
  : []))

const categoryOptions = computed(() => [])

watch(() => props.modelValue, (val) => {
  isOpen.value = val
  if (val && props.scheduled) {
    formData.value = {
      name: props.scheduled.name,
      type: props.scheduled.type,
      frequency: props.scheduled.frequency,
      start_date: props.scheduled.start_date,
      end_date: props.scheduled.end_date || '',
      cron_expression: props.scheduled.cron_expression || '',
      amount: String(props.scheduled.amount),
      currency: props.scheduled.currency,
      from_account_id: props.scheduled.from_account_id || '',
      to_account_id: props.scheduled.to_account_id || '',
      category_id: props.scheduled.category_id || '',
      description: props.scheduled.description || '',
      active: props.scheduled.active
    }
  } else if (val) {
    formData.value = defaultForm()
  }
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

watch(() => formData.value.frequency, (newFreq) => {
  if (newFreq !== 'custom') {
    formData.value.cron_expression = ''
  }
})

const handleSubmit = () => {
  emit('submit', {
    ...formData.value,
    amount: Number(formData.value.amount),
    from_account_id: formData.value.from_account_id || null,
    to_account_id: formData.value.to_account_id || null,
    category_id: formData.value.category_id || null,
    end_date: formData.value.end_date || null,
    cron_expression: formData.value.cron_expression || null
  })
  isOpen.value = false
}
</script>

<style scoped>
.scheduled-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.form-checkbox {
  flex-direction: row;
  align-items: center;
}

.form-checkbox .form-label {
  margin-bottom: 0;
}

.form-checkbox-input {
  width: 1rem;
  height: 1rem;
  accent-color: var(--color-primary);
}

.form-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
}
</style>