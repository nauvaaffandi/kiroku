<template>
  <BaseModal
    v-model="isOpen"
    :title="isEditing ? 'Edit Account' : 'New Account'"
    size="lg"
    @update:modelValue="$emit('update:modelValue', $event)"
  >
    <form @submit.prevent="handleSubmit" class="account-form">
      <div class="form-group">
        <label class="form-label">Name *</label>
        <BaseInput
          v-model="formData.name"
          type="text"
          required
          placeholder="Account name"
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
          <label class="form-label">Icon</label>
          <BaseInput
            v-model="formData.icon"
            type="text"
            placeholder="Emoji or leave blank"
            maxlength="4"
          />
        </div>
        <div class="form-group">
          <label class="form-label">Color</label>
          <input
            v-model="formData.color"
            type="color"
            class="form-input"
          />
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">Initial Balance</label>
        <BaseInput
          v-model="formData.balance"
          type="number"
          :step="0.01"
          placeholder="0.00"
        />
      </div>
      <div class="form-group form-checkbox">
        <label class="form-label">Archived</label>
        <input
          v-model="formData.archived"
          type="checkbox"
          class="form-checkbox-input"
        />
      </div>
      <div class="modal-actions">
        <BaseButton type="button" variant="secondary" @click="isOpen = false">
          Cancel
        </BaseButton>
        <BaseButton type="submit">
          {{ isEditing ? 'Save Changes' : 'Create Account' }}
        </BaseButton>
      </div>
    </form>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { BaseModal, BaseInput, BaseSelect, BaseButton } from '@/components/ui'
import type { Account, AccountType } from '@/types'

interface Props {
  modelValue: boolean
  account?: Account | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: [data: Partial<Account>]
}>()

const isOpen = ref(props.modelValue)
const isEditing = computed(() => !!props.account)

interface FormData {
  name: string
  type: AccountType
  currency: string
  icon: string
  color: string
  balance: string
  archived: boolean
}

const defaultForm = (): FormData => ({
  name: '',
  type: 'checking',
  currency: 'JPY',
  icon: '',
  color: '#2563eb',
  balance: '0',
  archived: false
})

const formData = ref<FormData>(defaultForm())

const typeOptions = [
  { value: 'cash', label: 'Cash' },
  { value: 'checking', label: 'Checking' },
  { value: 'savings', label: 'Savings' },
  { value: 'credit', label: 'Credit' },
  { value: 'investment', label: 'Investment' },
  { value: 'other', label: 'Other' }
]

watch(() => props.modelValue, (val) => {
  isOpen.value = val
  if (val && props.account) {
    formData.value = {
      name: props.account.name,
      type: props.account.type,
      currency: props.account.currency,
      icon: props.account.icon || '',
      color: props.account.color || '#2563eb',
      balance: String(props.account.balance || 0),
      archived: props.account.archived || false
    }
  } else if (val) {
    formData.value = defaultForm()
  }
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = () => {
  emit('submit', {
    name: formData.value.name,
    type: formData.value.type,
    currency: formData.value.currency,
    icon: formData.value.icon,
    color: formData.value.color,
    balance: Number(formData.value.balance),
    archived: formData.value.archived
  })
  isOpen.value = false
}
</script>

<style scoped>
@reference "tailwindcss";
.account-form {
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

.form-input {
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  font-family: inherit;
  color: var(--color-text);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-focus);
}

input[type="color"].form-input {
  width: 3rem;
  height: 2.5rem;
  padding: 0;
  border-radius: var(--radius);
  cursor: pointer;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
}
</style>