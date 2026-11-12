<template>
  <BaseModal
    v-model="isOpen"
    :title="isEditing ? t('accounts.formTitleEdit') : t('accounts.formTitleCreate')"
    size="lg"
    @update:modelValue="$emit('update:modelValue', $event)"
  >
    <form @submit.prevent="handleSubmit" class="account-form">
      <div class="form-group">
        <label class="form-label">{{ t('accounts.labelName') }} *</label>
        <BaseInput
          v-model="formData.name"
          type="text"
          required
          :placeholder="t('accounts.placeholderName')"
        />
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('accounts.labelType') }} *</label>
          <BaseSelect
            v-model="formData.type"
            :options="typeOptions"
            required
          />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('accounts.labelCurrency') }} *</label>
          <BaseInput
            v-model="formData.currency"
            type="text"
            required
            maxlength="3"
            :placeholder="t('accounts.placeholderCurrency')"
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('accounts.labelIcon') }}</label>
          <BaseInput
            v-model="formData.icon"
            type="text"
            :placeholder="t('accounts.placeholderIcon')"
            maxlength="4"
          />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('accounts.labelColor') }}</label>
          <ColorPalette v-model="formData.color" />
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">{{ t('accounts.labelInitialBalance') }}</label>
        <BaseInput
          v-model="formData.balance"
          type="number"
          :step="0.01"
          :placeholder="t('accounts.placeholderBalance')"
        />
      </div>
      <div class="form-group form-checkbox">
        <label class="form-label">{{ t('accounts.archived') }}</label>
        <input
          v-model="formData.archived"
          type="checkbox"
          class="form-checkbox-input"
        />
      </div>
      <div class="modal-actions">
        <BaseButton type="button" variant="secondary" @click="isOpen = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton type="submit">
          {{ isEditing ? t('common.save') : t('accounts.btnCreateAccount') }}
        </BaseButton>
      </div>
    </form>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { BaseModal, BaseInput, BaseSelect, BaseButton } from '@/components/ui'
import ColorPalette from '@/components/ui/ColorPalette.vue'
import { useI18n } from '@/locales'
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

const { t } = useI18n()

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

.color-native {
  @apply w-8 h-8 p-0 border-none cursor-pointer rounded-full;
  background-color: transparent;
}

.modal-actions {
  @apply flex justify-end gap-3 mt-4;
}
</style>