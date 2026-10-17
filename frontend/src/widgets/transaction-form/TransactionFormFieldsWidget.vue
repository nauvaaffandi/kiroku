<template>
  <form @submit.prevent="handleSubmit" class="transaction-form">
    <!-- Template Selector -->
    <div class="form-group" v-if="templates.length > 0 && !isEditing">
      <label>Load from Template</label>
      <select @change="applyTemplate($event)" class="form-input">
        <option value="">Select a template...</option>
        <option v-for="tmpl in templates" :key="tmpl.id" :value="tmpl.id">
          {{ tmpl.name }}
        </option>
      </select>
    </div>

    <div class="form-row">
      <div class="form-group">
        <label>Type</label>
        <select v-model="form.type" class="form-input" required>
          <option value="expense">Expense</option>
          <option value="income">Income</option>
          <option value="transfer">Transfer</option>
        </select>
      </div>
      <div class="form-group">
        <label>Date & Time</label>
        <input v-model="form.date_time" type="datetime-local" class="form-input" required />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group amount-group">
        <label>Amount</label>
        <div class="amount-input-wrapper">
          <input v-model.number="form.amount" type="number" step="1" min="1" class="form-input" required />
          <button type="button" @click="showCalculator = true" class="calc-btn" title="Open Calculator">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="4" y="2" width="16" height="20" rx="2"></rect>
              <line x1="8" y1="6" x2="16" y2="6"></line>
              <line x1="8" y1="10" x2="10" y2="10"></line>
              <line x1="14" y1="10" x2="16" y2="10"></line>
              <line x1="8" y1="14" x2="10" y2="14"></line>
              <line x1="14" y1="14" x2="16" y2="14"></line>
              <line x1="8" y1="18" x2="10" y2="18"></line>
              <line x1="14" y1="18" x2="16" y2="18"></line>
            </svg>
          </button>
        </div>
      </div>
      <div class="form-group">
        <label>Currency</label>
        <input v-model="form.currency" type="text" class="form-input" required maxlength="3" />
      </div>
    </div>

    <div class="form-group">
      <label>Category</label>
      <select v-model="form.category_id" class="form-input">
        <option value="">Select Category</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id" :disabled="cat.children && cat.children.length > 0">
          {{ cat.path }}
        </option>
      </select>
    </div>

    <div class="form-row" v-if="form.type === 'expense' || form.type === 'transfer'">
      <div class="form-group">
        <label>From Account</label>
        <select v-model="form.from_account_id" class="form-input" required>
          <option value="">Select Account</option>
          <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
            {{ acc.name }} ({{ acc.currency }}) - Bal: ¥{{ acc.balance }}
          </option>
        </select>
      </div>
    </div>

    <div class="form-row" v-if="form.type === 'income' || form.type === 'transfer'">
      <div class="form-group">
        <label>To Account</label>
        <select v-model="form.to_account_id" class="form-input" required>
          <option value="">Select Account</option>
          <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
            {{ acc.name }} ({{ acc.currency }}) - Bal: ¥{{ acc.balance }}
          </option>
        </select>
      </div>
    </div>

    <div class="form-group">
      <label>Note</label>
      <input
        v-model="form.note"
        type="text"
        class="form-input"
        placeholder="Optional note"
        list="notes-autocomplete"
      />
      <datalist id="notes-autocomplete">
        <option v-for="note in noteSuggestions" :key="note" :value="note" />
      </datalist>
    </div>

    <div class="form-group">
      <label>Payee / Payer</label>
      <input
        v-model="form.payee_payer"
        type="text"
        class="form-input"
        placeholder="Optional payee or payer"
        list="payee-autocomplete"
      />
      <datalist id="payee-autocomplete">
        <option v-for="payee in payeeSuggestions" :key="payee" :value="payee" />
      </datalist>
    </div>

    <div class="form-row">
      <div class="form-group">
        <label>Payment Type</label>
        <select v-model="form.payment_type" class="form-input">
          <option value="">Select payment type</option>
          <option value="cash">Cash</option>
          <option value="card">Card</option>
          <option value="qr">QR Code</option>
          <option value="bank">Bank Transfer</option>
        </select>
      </div>
      <div class="form-group">
        <label>Status</label>
        <select v-model="form.status" class="form-input">
          <option value="cleared">Cleared</option>
          <option value="pending">Pending</option>
          <option value="reconciled">Reconciled</option>
        </select>
      </div>
    </div>

    <div class="form-row">
      <div class="form-group">
        <label>Warranty Until</label>
        <input v-model="form.warranty_until" type="date" class="form-input" />
      </div>
      <div class="form-group">
        <label>Place Name</label>
        <div class="place-input-wrapper">
          <input v-model="form.place_name" type="text" class="form-input" placeholder="Location name" />
          <button type="button" @click="showMapModal = true" class="btn-secondary map-btn">Map</button>
        </div>
      </div>
    </div>

    <div class="form-group" v-if="!isEditing">
      <label class="checkbox-label">
        <input type="checkbox" v-model="saveAsTemplate" />
        Save as template
      </label>
      <input v-if="saveAsTemplate" v-model="templateName" type="text" class="form-input mt-2" placeholder="Template name" required />
    </div>

    <div class="modal-actions">
      <button type="button" @click="$emit('cancel')" class="btn-secondary">Cancel</button>
      <button type="submit" :disabled="loading" class="btn-primary">
        {{ loading ? 'Saving...' : (isEditing ? 'Update Transaction' : 'Create Transaction') }}
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { api } from '@/api/client'
import type { Category, Account } from '@/types'

const props = defineProps<{
  initialData?: any
  isEditing?: boolean
  categories: Category[]
  accounts: Account[]
}>()

const emit = defineEmits(['save', 'cancel'])

const loading = ref(false)
const templates = ref<any[]>([])
const saveAsTemplate = ref(false)
const templateName = ref('')

const noteSuggestions = ref<string[]>([])
const payeeSuggestions = ref<string[]>([])

const form = ref({
  type: 'expense',
  date_time: new Date().toISOString().slice(0, 16),
  amount: 0,
  currency: 'JPY',
  from_account_id: '',
  to_account_id: '',
  category_id: '',
  labels: '',
  note: '',
  payee_payer: '',
  payment_type: 'cash',
  warranty_until: null,
  status: 'cleared',
  place_name: '',
  place_lat: null as number | null,
  place_lng: null as number | null,
  attachments: '',
  template_id: null
})

// Calculator state
const showCalculator = ref(false)

// Map Modal state
const showMapModal = ref(false)

const fetchAutocomplete = async () => {
  try {
    const res = await api.get('/api/v1/transactions/autocomplete', { params: { q: '' } })
    if (res.data) {
      if (res.data.notes) noteSuggestions.value = res.data.notes
      if (res.data.payees) payeeSuggestions.value = res.data.payees
    }
  } catch (e) {
    // ignore
  }
}

const fetchTemplates = async () => {
  try {
    const res = await api.get('/api/v1/templates')
    templates.value = res.data || []
  } catch (e) {
    // ignore
  }
}

const applyTemplate = (e: Event) => {
  const id = (e.target as HTMLSelectElement).value
  if (!id) return
  const tmpl = templates.value.find(t => t.id === id)
  if (tmpl && tmpl.transaction_template_json) {
    const tjson = tmpl.transaction_template_json
    form.value.type = tjson.type || form.value.type
    form.value.amount = tjson.amount || form.value.amount
    form.value.currency = tjson.currency || form.value.currency
    form.value.from_account_id = tjson.from_account_id || form.value.from_account_id
    form.value.to_account_id = tjson.to_account_id || form.value.to_account_id
    form.value.category_id = tjson.category_id || form.value.category_id
    form.value.note = tjson.note || form.value.note
    form.value.payee_payer = tjson.payee_payer || form.value.payee_payer
    form.value.template_id = tmpl.id
  }
}

const handleSubmit = async () => {
  loading.value = true
  try {
    let tid = form.value.template_id
    if (saveAsTemplate.value && templateName.value) {
      const tres = await api.post('/api/v1/templates', {
        name: templateName.value,
        transaction_template_json: form.value
      })
      tid = tres.data.id
    }

    const payload = { ...form.value, template_id: tid }
    emit('save', payload)
  } catch (e) {
    // handled by parent
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  if (props.initialData) {
    form.value = { ...props.initialData }
    if (form.value.date_time) {
      form.value.date_time = new Date(form.value.date_time).toISOString().slice(0, 16)
    }
  }
  await Promise.all([
    fetchAutocomplete(),
    fetchTemplates()
  ])
})

watch(() => props.initialData, (newData) => {
  if (newData) {
    form.value = { ...newData }
    if (form.value.date_time) {
      form.value.date_time = new Date(form.value.date_time).toISOString().slice(0, 16)
    }
  }
}, { deep: true })
</script>

<style scoped>
.transaction-form { display: flex; flex-direction: column; gap: 1rem; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.form-group { display: flex; flex-direction: column; gap: 0.375rem; }
.form-group label { font-size: 0.875rem; font-weight: 500; color: var(--color-neutral-700); }
.form-input { width: 100%; padding: 0.5rem 0.75rem; border: 1px solid var(--color-neutral-300); border-radius: var(--radius-sm); background-color: var(--color-neutral-100); color: var(--color-neutral-900); font-size: 0.875rem; }
.form-input:focus { outline: none; box-shadow: 0 0 0 2px var(--color-primary); }
.amount-input-wrapper { display: flex; gap: 0.5rem; }
.calc-btn { padding: 0.5rem 0.75rem; background-color: var(--color-neutral-200); border-radius: var(--radius-sm); color: var(--color-neutral-700); display: flex; align-items: center; justify-content: center; cursor: pointer; transition: background-color 0.2s ease; }
.calc-btn:hover { background-color: var(--color-neutral-300); }
.place-input-wrapper { display: flex; gap: 0.5rem; }
.map-btn { white-space: nowrap; }
.checkbox-label { display: flex; align-items: center; gap: 0.5rem; font-size: 0.875rem; font-weight: 500; color: var(--color-neutral-700); cursor: pointer; }
.mt-2 { margin-top: 0.5rem; }
.modal-actions { display: flex; justify-content: flex-end; gap: 0.75rem; padding-top: 1rem; border-top: 1px solid var(--color-neutral-200); }
</style>