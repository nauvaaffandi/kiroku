<template>
  <Teleport to="body">
    <div v-if="open" class="modal-overlay" @click.self="$emit('close')">
      <div class="modal-container">
        <div class="modal-header">
          <h3>{{ title }}</h3>
          <button @click="$emit('close')" class="modal-close">{{ t('common.close') }}</button>
        </div>
        <form @submit.prevent="save" class="modal-body">
          <label class="form-row">
            <span class="form-label">{{ t('budgets.field.name') }}</span>
            <input v-model="form.name" type="text" :required="true" maxlength="255" class="form-input" />
          </label>
          <div class="form-row">
            <span class="form-label">{{ t('budgets.field.period') }}</span>
            <select v-model="form.period" class="form-input">
              <option value="weekly">{{ t('budgets.period.weekly') }}</option>
              <option value="monthly">{{ t('budgets.period.monthly') }}</option>
              <option value="yearly">{{ t('budgets.period.yearly') }}</option>
              <option value="one_time">{{ t('budgets.period.oneTime') }}</option>
            </select>
          </div>
          <div class="form-row">
            <span class="form-label">{{ t('budgets.field.amount') }}</span>
            <input v-model.number="form.amount" type="number" step="0.01" min="0.01" required class="form-input" />
          </div>
          <div class="form-row">
            <span class="form-label">{{ t('common.currency') }}</span>
            <input v-model="form.currency" type="text" maxlength="3" pattern="[A-Z]{3}" class="form-input" spellcheck="false" />
          </div>
          <div class="form-row">
            <span class="form-label">{{ t('accounts.type') }}</span>
            <select v-model="form.account_id" class="form-input">
              <option value="" disabled>{{ t('reports.allAccounts') }}</option>
              <option v-for="a in accountsStore.accounts" :key="a.id" :value="a.id">{{ a.name }} ({{ a.currency }})</option>
            </select>
          </div>
          <div class="form-row form-row--checkboxes">
            <span class="form-label">{{ t('categories.labelIcon') }}</span>
            <select v-model="form.category_ids" multiple class="form-input form-input--multi">
              <option v-for="c in categoriesStore.categories" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
          </div>
          <div class="form-row form-row--checkboxes">
            <span class="form-label">{{ t('labels.name') }}</span>
            <select v-model="form.label_ids" multiple class="form-input form-input--multi">
              <option v-for="l in labelsStore.labels" :key="l.id" :value="l.id">{{ l.name }}</option>
            </select>
          </div>
          <div class="form-row form-row--checkboxes">
            <label class="toggle"><input type="checkbox" v-model="form.trending_over" /><span>{{ t('budgets.toggle.trendingOver') }}</span></label>
            <label class="toggle"><input type="checkbox" v-model="form.over_budget" /><span>{{ t('budgets.toggle.overBudget') }}</span></label>
          </div>
          <div v-if="form.period === 'one_time'" class="form-row form-row--checkboxes">
            <div class="inline-fields">
              <div class="form-row">
                <span class="form-label">{{ t('budgets.field.startsAt') }}</span>
                <input v-model="form.starts_at" type="date" class="form-input" />
              </div>
              <div class="form-row">
                <span class="form-label">{{ t('budgets.field.endsAt') }}</span>
                <input v-model="form.ends_at" type="date" class="form-input" />
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" @click="cancel" class="btn-cancel">{{ t('common.cancel') }}</button>
            <button type="submit" class="btn-primary">{{ title }}</button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, ref, reactive, watch } from 'vue'
import { localDb } from '@/storage/repository'
import { useI18n } from '@/locales'
import { useCategoriesStore } from '@/stores/categories'
import { useAccountsStore } from '@/stores/accounts'
import { useLabelsStore } from '@/stores/labels'

interface FormState {
  id?: string
  name: string
  period: string
  amount: number
  currency: string
  account_id?: string
  category_ids: string[]
  label_ids: string[]
  trending_over: boolean
  over_budget: boolean
  starts_at?: string
  ends_at?: string
}

const props = defineProps<{
  open: boolean
  id?: string | null
}>()

const emit = defineEmits<{
  close: []
  save: [state: FormState]
}>()

const categoriesStore = useCategoriesStore()
const accountsStore = useAccountsStore()
const labelsStore = useLabelsStore()

const { t } = useI18n()

const loading = ref(false)
const form = reactive<FormState>({
  name: '',
  period: 'monthly',
  amount: 0,
  currency: 'USD',
  category_ids: [],
  label_ids: [],
  trending_over: false,
  over_budget: false
})

watch(() => [props.open, props.id], async ([openVal]) => {
  if (!openVal) return
  form.id = props.id ?? undefined
  loading.value = true
  try {
    await Promise.all([categoriesStore.fetchCategories(), accountsStore.fetchAccounts(), labelsStore.fetchLabels()])
  } catch (e) { /* stores surface their own errors */ }
  if (props.id) await fetchCurrent()
  loading.value = false
}, { immediate: true })

async function fetchCurrent() {
  if (!props.open || !form.id) return
  try {
    const b = await localDb.get<any>('budget', form.id)
    if (!b) return
    form.name = b.name
    form.period = b.period
    form.amount = b.amount ?? 0
    form.currency = b.currency ?? 'USD'
    form.account_id = b.account_id ?? ''
    form.category_ids = (b.category_ids ?? []).map((c: any) => c)
    form.label_ids = (b.label_ids ?? []).map((l: any) => l)
    form.trending_over = !!b.trending_over
    form.over_budget = !!b.over_budget
    if (b.starts_at) form.starts_at = sliceDate(b.starts_at)
    if (b.ends_at) form.ends_at = sliceDate(b.ends_at)
  } catch (e) {
    console.error('Failed to fetch budget:', e)
    loading.value = false
    return
  }
}

function sliceDate(iso: string): string {
  return iso.slice(0, 10)
}

function resetForm() {
  Object.assign(form, {
    id: undefined, name: '', period: 'monthly', amount: 0, currency: 'USD',
    account_id: undefined, category_ids: [], label_ids: [],
    trending_over: false, over_budget: false, starts_at: undefined, ends_at: undefined
  })
}

async function save() {
  if (!form.name.trim() || !form.amount || !form.currency) {
    return
  }
  try {
    const payload = {
      name: form.name.trim(),
      period: form.period,
      amount: form.amount,
      currency: form.currency.toUpperCase(),
      account_id: form.account_id || undefined,
      category_ids: form.category_ids,
      label_ids: form.label_ids,
      trending_over: form.trending_over,
      over_budget: form.over_budget,
      ...(form.starts_at ? { starts_at: form.starts_at } : {}),
      ...(form.ends_at ? { ends_at: form.ends_at } : {})
    }
    if (form.id) {
      const current = await localDb.get<any>('budget', form.id)
      await localDb.put('budget', form.id, { ...current, ...payload })
    } else {
      await localDb.put('budget', crypto.randomUUID(), payload)
    }
    emit('save', { ...form })
    emit('close')
  } catch (e) {
    console.error('Failed to save budget:', e)
    return
  }
}

function cancel() {
  resetForm()
  emit('close')
}

const title = computed(() => (form.id ? t('common.edit') : t('budgets.add')))
</script>

<style scoped>
@reference "tailwindcss";
.modal-overlay {
  @apply fixed inset-0 flex items-center justify-center p-4 z-[1000];
  background-color: rgba(0, 0, 0, 0.5);
}
.modal-container {
  @apply border rounded-xl w-full max-w-[560px] max-h-[90vh] overflow-auto;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}
.modal-header {
  @apply flex items-center justify-between py-4 px-6 border-b;
  border-color: var(--border-primary);
}
.modal-close {
  @apply p-2 rounded-sm bg-transparent border-none text-zinc-700 transition-colors;
}
.modal-close:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}
.modal-body {
  @apply p-6 pt-2 grid grid-cols-[1fr_repeat(2,minmax(0,1fr))] gap-4;
}
.form-row {
  @apply flex flex-col gap-1;
}
.form-row--checkboxes {
  @apply flex items-center gap-2;
}
.form-label {
  @apply text-xs uppercase tracking-wider text-zinc-500;
}
.form-input {
  @apply block w-full rounded-sm border bg-transparent px-3 py-2 text-sm;
  border-color: var(--border-focus);
  color: var(--text-primary);
}
.form-input--multi {
  min-height: 96px;
}
.inline-fields {
  @apply flex flex-col gap-2;
}
.toggle {
  @apply inline-flex items-center gap-2 cursor-pointer;
  color: var(--text-secondary);
}
.toggle input { width: 1rem; height: 1rem; }
.modal-footer {
  @apply flex justify-end gap-2 mt-4;
}
.btn-cancel {
  @apply px-4 py-2 rounded-sm font-medium text-sm border border-zinc-300 bg-white;
  color: var(--text-secondary);
}
.btn-primary {
  @apply inline-flex items-center justify-center px-4 py-2 rounded-sm font-medium text-sm no-underline;
  background-color: var(--accent);
  color: var(--text-inverse);
  transition: opacity 0.2s ease;
}
.btn-primary:hover {
  @apply opacity-90;
}
</style>
