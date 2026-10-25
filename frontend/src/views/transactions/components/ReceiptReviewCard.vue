<template>
  <div class="review-card" v-if="draft">
    <h3>Review Transaction Draft</h3>
    <p class="text-muted">Verify and edit extracted details before confirming final persistence.</p>

    <form @submit.prevent="$emit('confirm', draft)" class="transaction-form">
      <div class="form-grid">
        <div class="form-group">
          <label>Type</label>
          <select v-model="draft.type" class="form-input" required>
            <option value="expense">Expense</option>
            <option value="income">Income</option>
          </select>
        </div>

        <div class="form-group">
          <label>Date & Time</label>
          <input v-model="draft.date_time" type="datetime-local" class="form-input" required />
        </div>

        <div class="form-group">
          <label>Amount</label>
          <input v-model.number="draft.amount" type="number" step="1" class="form-input" required />
        </div>

        <div class="form-group">
          <label>Currency</label>
          <input v-model="draft.currency" type="text" class="form-input" required maxlength="3" />
        </div>

        <div class="form-group">
          <label>Account</label>
          <select v-model="draft.from_account_id" class="form-input" required>
            <option value="">Select Account</option>
            <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
              {{ acc.name }} ({{ acc.currency }})
            </option>
          </select>
        </div>

        <div class="form-group">
          <label>Category</label>
          <select v-model="draft.category_id" class="form-input">
            <option value="">Select Category</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label>Payee / Merchant</label>
          <input v-model="draft.payee_payer" type="text" class="form-input" placeholder="Merchant name" />
        </div>

        <div class="form-group">
          <label>Payment Type</label>
          <select v-model="draft.payment_type" class="form-input">
            <option value="cash">Cash</option>
            <option value="card">Card</option>
            <option value="qr">QR Code</option>
            <option value="bank">Bank Transfer</option>
            <option value="other">Other</option>
          </select>
        </div>

        <div class="form-group full-width">
          <label>Note</label>
          <input v-model="draft.note" type="text" class="form-input" placeholder="Optional note" />
        </div>
      </div>

      <div class="modal-actions">
        <BaseButton type="button" variant="secondary" @click="$emit('reset')">Discard & Scan Another</BaseButton>
        <BaseButton type="submit" :disabled="submitting">Confirm & Save Transaction</BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/ui'
import type { AINormalizationResult, Category, Account } from '@/types'

defineProps<{
  draft: AINormalizationResult
  accounts: Account[]
  categories: Category[]
  submitting: boolean
}>()

defineEmits(['confirm', 'reset'])
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