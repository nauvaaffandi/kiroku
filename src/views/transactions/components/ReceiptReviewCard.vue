<template>
  <div class="review-card" v-if="draft">
    <h3>{{ t('receipt.reviewTitle') }}</h3>
    <p class="text-muted">{{ t('receipt.reviewDesc') }}</p>

    <form @submit.prevent="$emit('confirm', draft)" class="transaction-form">
      <div class="form-grid">
        <div class="form-group">
          <label>{{ t('transactions.typeExpense') }}</label>
          <select v-model="draft.type" class="form-input" required>
            <option value="expense">{{ t('transactions.typeExpense') }}</option>
            <option value="income">{{ t('transactions.typeIncome') }}</option>
          </select>
        </div>

        <div class="form-group">
          <label>{{ t('receipt.dateTime') }}</label>
          <input v-model="draft.date_time" type="datetime-local" class="form-input" required />
        </div>

        <div class="form-group">
          <label>{{ t('common.amount') }}</label>
          <input v-model.number="draft.amount" type="number" step="1" class="form-input" required />
        </div>

        <div class="form-group">
          <label>{{ t('common.currency') }}</label>
          <input v-model="draft.currency" type="text" class="form-input" required maxlength="3" />
        </div>

        <div class="form-group">
          <label>{{ t('common.account') }}</label>
          <select v-model="draft.from_account_id" class="form-input" required>
            <option value="">{{ t('common.selectAccount') }}</option>
            <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
              {{ acc.name }} ({{ acc.currency }})
            </option>
          </select>
        </div>

        <div class="form-group">
          <label>{{ t('common.category') }}</label>
          <select v-model="draft.category_id" class="form-input">
            <option value="">{{ t('common.selectCategory') }}</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label>{{ t('receipt.payeeMerchant') }}</label>
          <input v-model="draft.payee_payer" type="text" class="form-input" :placeholder="t('receipt.merchantPlaceholder')" />
        </div>

        <div class="form-group">
          <label>{{ t('common.paymentType') }}</label>
          <select v-model="draft.payment_type" class="form-input">
            <option value="cash">{{ t('common.cash') }}</option>
            <option value="card">{{ t('common.card') }}</option>
            <option value="qr">{{ t('common.qrCode') }}</option>
            <option value="bank">{{ t('common.bankTransfer') }}</option>
            <option value="other">{{ t('common.other') }}</option>
          </select>
        </div>

        <div class="form-group full-width">
          <label>{{ t('common.note') }}</label>
          <input v-model="draft.note" type="text" class="form-input" :placeholder="t('common.optionalNote')" />
        </div>
      </div>

      <div class="modal-actions">
        <BaseButton type="button" variant="secondary" @click="$emit('reset')">{{ t('receipt.discardScan') }}</BaseButton>
        <BaseButton type="submit" :disabled="submitting">{{ t('receipt.confirmSave') }}</BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { BaseButton } from '@/components/ui'
import type { AINormalizationResult, Category, Account } from '@/types'

defineProps<{
  draft: AINormalizationResult
  accounts: Account[]
  categories: Category[]
  submitting: boolean
}>()

defineEmits(['confirm', 'reset'])

const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";
.review-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  padding: 2rem;
  box-shadow: var(--shadow);
}

.transaction-form {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.form-group.full-width {
  grid-column: span 2;
}

.form-group label {
  font-size: 0.8125rem;
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background-color: var(--color-surface);
  color: var(--color-text);
  font-size: 0.875rem;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
}
</style>