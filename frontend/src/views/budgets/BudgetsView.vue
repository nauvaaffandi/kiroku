<template>
  <div class="budgets-view">
    <div class="page-header">
      <h1>{{ t('budgets.title') }}</h1>
      <p class="page-subtitle">{{ t('budgets.subtitle') }}</p>
      <button @click="openCreate" class="btn-primary">{{ t('budgets.add') }}</button>
    </div>

    <BudgetSummaryCardsWidget :loading="loading" :summary="summary" />

    <BudgetTableWidget
      :loading="loading"
      :budgets="rows"
      @export="exportCSV"
      @edit="openEdit"
      @delete="confirmDelete"
    />

    <BudgetFormModal :open="modalOpen" :id="editingId" @save="fetchData" @close="modalOpen = false" />

    <Teleport to="body">
      <div v-if="deleteConfirm" class="modal-overlay" @click.self="cancelDelete">
        <div class="modal-container confirm">
          <h3>{{ t('budgets.deleteConfirmTitle') }}</h3>
          <p>{{ t('budgets.deleteConfirmBody', { name: deletingName }) }}</p>
          <div class="confirm-actions">
            <button @click="cancelDelete" class="btn-cancel">{{ t('common.cancel') }}</button>
            <button @click="doDelete" class="btn-danger">{{ t('common.delete') }}</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'
import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import BudgetSummaryCardsWidget from '@/widgets/budgets/BudgetSummaryCardsWidget.vue'
import BudgetTableWidget from '@/widgets/budgets/BudgetTableWidget.vue'
import BudgetFormModal from '@/widgets/budgets/BudgetFormModal.vue'

interface BudgetResult {
  budget_id: string
  name: string
  period: string
  amount: string
  currency: string
  spent: string
  remaining: string
  percentage_used: number
  is_over_budget: boolean
}

const { t } = useI18n()
const uiStore = useUIStore()

const loading = ref(true)
const rows = ref<any[]>([])
const summary = ref({ totalBudgets: 0, overBudgets: 0, underBudgets: 0 })

const modalOpen = ref(false)
const editingId = ref<string | null>(null)

const deleteConfirm = ref(false)
const deletingId = ref<string | null>(null)
const deletingName = ref('')

async function fetchData() {
  loading.value = true
  try {
    const res = await api.get<BudgetResult[]>('/api/v1/budgets/')
    rows.value = (res.data || []).map(toRow)
    let over = 0, under = 0
    for (const b of rows.value) { if (b.is_over_budget) over++; else under++ }
    summary.value = { totalBudgets: rows.value.length, overBudgets: over, underBudgets: under }
  } catch (e) {
    console.error('Failed to fetch budgets:', e)
    uiStore.error(t('budgets.fetchError'))
  } finally {
    loading.value = false
  }
}

function toRow(b: BudgetResult) {
  return {
    id: b.budget_id, name: b.name, period: b.period,
    amount: parseNum(b.amount), currency: b.currency,
    spent: parseNum(b.spent), remaining: parseNum(b.remaining),
    percentageUsed: b.percentage_used || 0,
    status: b.is_over_budget ? 'over' : 'under'
  }
}

function parseNum(s: string): number { return parseFloat(s) || 0 }

function exportCSV() {
  const headers = [t('budgets.table.name'), t('budgets.table.period'), t('budgets.table.budget'), t('budgets.table.spent'), t('budgets.table.remaining'), t('budgets.table.usage'), t('budgets.table.status')]
  const body = rows.value.map(b => [b.name, t(`budgets.period.${b.period}`), b.amount + '', b.spent + '', b.remaining + '', b.percentageUsed + '%', t(`budgets.status.${b.status}`)])
  const csv = [headers, ...body].map(r => r.map(v => `"${v}"`).join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `budgets-${Date.now()}.csv`
  link.click()
}

function openCreate() { editingId.value = null; modalOpen.value = true }
function openEdit(id: string) { editingId.value = id; modalOpen.value = true }
function confirmDelete(row: any) { deletingId.value = row.id; deletingName.value = row.name; deleteConfirm.value = true }
function cancelDelete() { deleteConfirm.value = false; deletingId.value = null; deletingName.value = '' }

async function doDelete() {
  if (!deletingId.value) return
  try {
    await api.delete(`/api/v1/budgets/${deletingId.value}`)
    cancelDelete()
    await fetchData()
  } catch (e) {
    console.error('Failed to delete budget:', e)
    uiStore.error((e as any).response?.data?.error || t('budgets.deleteError'))
  }
}

onMounted(() => fetchData())
</script>

<style scoped>
@reference "tailwindcss";
.budgets-view { @apply flex flex-col gap-4 p-6; }
.page-header { @apply pb-4 flex items-center justify-between gap-4 w-full; border-bottom: 1px solid var(--border-primary); }
.page-header h1 { @apply text-xl font-semibold; color: var(--text-primary); }
.page-subtitle { @apply text-[0.875rem] mt-1; color: var(--text-secondary); }
.btn-primary { @apply inline-flex items-center justify-center px-4 py-2 rounded-sm font-medium text-sm no-underline; background-color: var(--accent); color: var(--text-inverse); transition: opacity 0.2s ease; }
.btn-primary:hover { @apply opacity-90; }
.btn-secondary { @apply inline-flex items-center gap-[0.375rem] border rounded-sm font-medium text-sm; border-color: var(--border-primary); background-color: var(--bg-card); color: var(--text-primary); }
.modal-overlay { @apply fixed inset-0 flex items-center justify-center p-4 z-[1000]; background-color: rgba(0, 0, 0, 0.5); }
.modal-container { @apply border rounded-xl w-full max-w-[560px] max-h-[90vh] overflow-auto; background-color: var(--bg-card); border-color: var(--border-primary); }
.confirm { @apply p-6 pt-4 gap-4 text-sm; }
.confirm-actions { @apply flex justify-end gap-2 mt-4; }
.btn-cancel { @apply px-4 py-2 rounded-sm font-medium text-sm border border-zinc-300 bg-white; color: var(--text-secondary); }
.btn-danger { @apply px-4 py-2 rounded-sm font-medium text-sm no-underline; background-color: var(--error); color: white; }
</style>
