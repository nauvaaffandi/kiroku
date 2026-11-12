<template>
  <BaseModal
    v-model="isOpen"
    :title="isEditing ? t('scheduled.formTitleEdit') : t('scheduled.formTitleCreate')"
    size="lg"
    @update:modelValue="$emit('update:modelValue', $event)"
  >
    <form @submit.prevent="handleSubmit" class="scheduled-form">
      <div class="form-group">
        <label class="form-label">{{ t('scheduled.labelName') }} *</label>
        <BaseInput
          v-model="formData.name"
          type="text"
          required
          :placeholder="t('scheduled.placeholderName')"
        />
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('scheduled.labelType') }} *</label>
          <BaseSelect
            v-model="formData.type"
            :options="typeOptions"
            required
          />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('scheduled.labelFrequency') }} *</label>
          <BaseSelect
            v-model="formData.frequency"
            :options="frequencyOptions"
            required
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('scheduled.labelStartDate') }} *</label>
          <BaseInput
            v-model="formData.start_date"
            type="date"
            required
          />
        </div>
        <div class="form-group" v-if="formData.frequency === 'custom'">
          <label class="form-label">{{ t('scheduled.labelCron') }}</label>
          <BaseInput
            v-model="formData.cron_expression"
            type="text"
            :placeholder="t('scheduled.placeholderCron')"
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('scheduled.labelEndDate') }}</label>
          <BaseInput
            v-model="formData.end_date"
            type="date"
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('common.amount') }} *</label>
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
          <label class="form-label">{{ t('common.currency') }} *</label>
          <BaseInput
            v-model="formData.currency"
            type="text"
            required
            maxlength="3"
            :placeholder="t('scheduled.placeholderCurrency')"
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('common.fromAccount') }} *</label>
          <BaseSelect
            v-model="formData.from_account_id"
            :options="accountOptions"
            required
          />
        </div>
        <div class="form-group" v-if="formData.type === 'transfer'">
          <label class="form-label">{{ t('common.toAccount') }} *</label>
          <BaseSelect
            v-model="formData.to_account_id"
            :options="accountOptions"
            required
          />
        </div>
        <div class="form-group" v-else>
          <label class="form-label">{{ t('common.category') }}</label>
          <BaseSelect
            v-model="formData.category_id"
            :options="categoryOptions"
            :placeholder="t('scheduled.placeholderCategory')"
          />
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">{{ t('common.description') }}</label>
        <BaseInput
          v-model="formData.description"
          type="text"
          :placeholder="t('scheduled.placeholderDescription')"
        />
      </div>
      <div class="form-group form-checkbox">
        <label class="form-label">{{ t('scheduled.active') }}</label>
        <input
          v-model="formData.active"
          type="checkbox"
          class="form-checkbox-input"
        />
      </div>
      <div class="modal-actions">
        <BaseButton type="button" variant="secondary" @click="isOpen = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton type="submit">
          {{ isEditing ? t('common.save') : t('scheduled.btnCreate') }}
        </BaseButton>
      </div>
    </form>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { BaseModal, BaseInput, BaseSelect, BaseButton } from '@/components/ui'
import { useI18n } from '@/locales'
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

const { t } = useI18n()

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
@reference "tailwindcss";
.scheduled-form { @apply flex flex-col gap-5; }
.form-row { @apply grid grid-cols-2 gap-4; }
.form-group { @apply flex flex-col gap-1.5; }
.form-checkbox { @apply flex-row items-center; }
.form-checkbox .form-label { @apply mb-0; }
.form-checkbox-input { accent-color: var(--accent); @apply w-4 h-4; }
.form-label { color: var(--text-primary); @apply text-[0.8125rem] font-medium; }
.modal-actions { @apply flex justify-end gap-3 mt-4; }
</style>