<template>
  <BaseModal
    v-model="isOpen"
    :title="isEditing ? t('transactions.editTransaction') : t('transactions.createTransaction')"
    size="lg"
    @update:modelValue="$emit('update:modelValue', $event)"
  >
    <form @submit.prevent="handleSubmit" class="transaction-form">
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label">{{ t('transactions.labelType') }} *</label>
          <BaseSelect
            v-model="form.type"
            :options="typeOptions"
            required
          />
        </div>

        <div class="form-group">
          <label class="form-label">{{ t('transactions.labelDateTime') }} *</label>
          <BaseInput
            v-model="form.date_time"
            type="datetime-local"
            required
          />
        </div>

        <div class="form-group">
          <label class="form-label">{{ t('common.amount') }} *</label>
          <BaseInput
            v-model.number="form.amount"
            type="number"
            :min="0"
            :step="0.01"
            required
            :placeholder="t('transactions.placeholderAmount')"
          />
        </div>

        <div class="form-group">
          <label class="form-label">{{ t('common.currency') }} *</label>
          <BaseSelect
            v-model="form.currency"
            :options="currencyOptions"
            required
          />
        </div>

        <div class="form-group" v-if="form.type !== 'transfer'">
          <label class="form-label">{{ t('common.account') }} *</label>
          <BaseSelect
            v-model="form.from_account_id"
            :options="accountOptions"
            required
          />
        </div>

        <div class="form-group" v-if="form.type === 'transfer'">
          <label class="form-label">{{ t('common.fromAccount') }} *</label>
          <BaseSelect
            v-model="form.from_account_id"
            :options="accountOptions"
            required
          />
        </div>

        <div class="form-group" v-if="form.type === 'transfer'">
          <label class="form-label">{{ t('common.toAccount') }} *</label>
          <BaseSelect
            v-model="form.to_account_id"
            :options="toAccountOptions"
            required
          />
        </div>

        <div class="form-group">
          <label class="form-label">{{ t('common.category') }}</label>
          <BaseSelect
            v-model="form.category_id"
            :options="categoryOptions"
            placeholder="Select category (optional)"
          />
        </div>

        <div class="form-group">
          <label class="form-label">{{ t('common.paymentType') }}</label>
          <BaseInput
            v-model="form.payment_type"
            :placeholder="t('transactions.placeholderPaymentType')"
          />
        </div>

        <div class="form-group">
          <label class="form-label">{{ t('transactions.labelPayeePayer') }}</label>
          <BaseInput
            v-model="form.payee_payer"
            :placeholder="t('transactions.placeholderPayeePayer')"
          />
        </div>

        <div class="form-group">
          <label class="form-label">{{ t('transactions.labelLabels') }}</label>
          <BaseInput
            v-model="form.labels"
            :placeholder="t('transactions.placeholderLabels')"
          />
        </div>

        <div class="form-group">
          <label class="form-label">{{ t('common.note') }}</label>
          <textarea
            v-model="form.note"
            class="form-textarea"
            rows="3"
            :placeholder="t('common.optionalNote')"
          />
        </div>

        <div class="form-group">
          <label class="form-label">{{ t('transactions.labelPlace') }}</label>
          <BaseInput
            v-model="form.place_name"
            :placeholder="t('common.locationName')"
          />
        </div>

        <div class="form-group">
          <label class="form-label">{{ t('common.status') }} *</label>
          <BaseSelect
            v-model="form.status"
            :options="statusOptions"
            required
          />
        </div>

        <div class="form-group" v-if="form.type === 'expense' || form.type === 'income'">
          <label class="form-label">{{ t('common.warrantyUntil') }}</label>
          <BaseInput
            v-model="form.warranty_until"
            type="date"
          />
        </div>
      </div>

      <div class="modal-actions">
        <BaseButton type="button" variant="secondary" @click="isOpen = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton type="submit">
          {{ isEditing ? t('common.save') : t('transactions.createTransaction') }}
        </BaseButton>
      </div>
    </form>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { BaseModal, BaseInput, BaseSelect, BaseButton } from '@/components/ui'
import { useI18n } from '@/locales'
import type { Transaction, Category, Account, TransactionType, TransactionStatus } from '@/types'
import { useCurrencies } from '@/composables/useCurrencies'
import { useSettingsStore } from '@/stores/settings'

interface Props {
  modelValue: boolean
  transaction?: Transaction | null
  accounts: Account[]
  categories: Category[]
  initialDateTime?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: [data: Partial<Transaction>]
}>()

const { t } = useI18n()

const settingsStore = useSettingsStore()
const { currencies } = useCurrencies()
const currencyOptions = computed(() => currencies.value.map(c => ({ value: c, label: c })))

const defaultCurrency = computed(() => settingsStore.baseCurrency)

const isOpen = ref(props.modelValue)
const isEditing = computed(() => !!props.transaction)

interface FormData {
  type: TransactionType
  date_time: string
  amount: string
  currency: string
  from_account_id: string
  to_account_id: string
  category_id: string
  payment_type: string
  payee_payer: string
  labels: string
  note: string
  place_name: string
  warranty_until: string
  status: TransactionStatus
}

const defaultForm = (): FormData => ({
  type: 'expense',
  date_time: new Date().toISOString().slice(0, 16),
  amount: '0',
  currency: defaultCurrency.value,
  from_account_id: '',
  to_account_id: '',
  category_id: '',
  payment_type: '',
  payee_payer: '',
  labels: '',
  note: '',
  place_name: '',
  warranty_until: '',
  status: 'cleared'
})

const form = ref<FormData>(defaultForm())

const typeOptions = [
  { value: 'income', label: 'Income' },
  { value: 'expense', label: 'Expense' },
  { value: 'transfer', label: 'Transfer' }
]

const statusOptions = [
  { value: 'cleared', label: 'Cleared' },
  { value: 'pending', label: 'Pending' }
]

const accountOptions = computed(() => (props.accounts || []).map(a => ({
  value: a.id,
  label: a.name
})))

const toAccountOptions = computed(() => {
  const fromId = form.value.from_account_id
  return (props.accounts || [])
    .filter(a => a.id !== fromId)
    .map(a => ({ value: a.id, label: a.name }))
})

const categoryOptions = computed(() => {
  return (props.categories || [])
    .filter(c => !c.children || c.children.length === 0)
    .map(c => ({ value: c.id, label: c.path || c.name }))
})

watch(() => props.modelValue, (val) => {
  isOpen.value = val
  if (val && props.transaction) {
    form.value = {
      type: props.transaction.type,
      date_time: props.transaction.date_time,
      amount: String(props.transaction.amount),
      currency: props.transaction.currency,
      from_account_id: props.transaction.from_account_id || '',
      to_account_id: props.transaction.to_account_id || '',
      category_id: props.transaction.category_id || '',
      payment_type: props.transaction.payment_type || '',
      payee_payer: props.transaction.payee_payer || '',
       labels: Array.isArray(props.transaction.labels) ? props.transaction.labels.join(', ') : (props.transaction.labels || ''),
      note: props.transaction.note || '',
      place_name: props.transaction.place_name || '',
      warranty_until: props.transaction.warranty_until || '',
      status: props.transaction.status
    }
  } else if (val) {
    form.value = defaultForm()
    if (props.initialDateTime) {
      form.value.date_time = props.initialDateTime
    }
  }
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

watch(() => form.value.type, (newType) => {
  if (newType === 'transfer') {
    if (!form.value.to_account_id && form.value.from_account_id) {
      const match = (props.accounts || []).find(a => a.id !== form.value.from_account_id)
      form.value.to_account_id = match ? match.id : ''
    }
  }
})

const handleSubmit = () => {
  emit('submit', {
    ...form.value,
    amount: Number(form.value.amount),
    from_account_id: form.value.from_account_id || null,
    to_account_id: form.value.to_account_id || null,
    category_id: form.value.category_id || null,
    warranty_until: form.value.warranty_until || null
  })
  isOpen.value = false
}
</script>

<style scoped>
@reference "tailwindcss";
.transaction-form {
  @apply flex flex-col gap-6;
}

.form-grid {
  @apply grid gap-4;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
}

.form-group {
  @apply flex flex-col gap-1.5;
}

.form-label {
  color: var(--text-primary);
  @apply text-[0.8125rem] font-medium;
}

.form-textarea {
  color: var(--text-primary);
  background-color: var(--bg-card);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  @apply w-full px-3 py-2 text-sm font-sans resize-y transition-[border-color,box-shadow] duration-150 ease-in-out;
}

.form-textarea:focus {
  border-color: var(--border-focus);
  box-shadow: 0 0 0 3px var(--accent-light);
  @apply outline-none;
}

.form-textarea::placeholder {
  color: var(--text-muted);
}

.modal-actions {
  @apply flex justify-end gap-3 mt-4;
}
</style>