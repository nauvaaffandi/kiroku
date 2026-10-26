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
  @apply flex flex-col gap-5;
}

.form-row {
  @apply grid grid-cols-2 gap-4;
}

.form-group {
  @apply flex flex-col gap-1.5;
}

.form-checkbox {
  @apply flex-row items-center;
}

.form-checkbox .form-label {
  @apply mb-0;
}

.form-checkbox-input {
  accent-color: var(--accent);
  @apply w-4 h-4;
}

.form-input {
  color: var(--text-primary);
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  @apply px-3 py-2 text-sm font-sans rounded transition-[border-color,box-shadow] duration-150 ease-in-out;
}

.form-input:focus {
  @apply outline-none;
  border-color: var(--accent);
}

input[type="color"].form-input {
  @apply w-12 h-10 p-0 cursor-pointer rounded;
}

.modal-actions {
  @apply flex justify-end gap-3 mt-4;
}
</style>