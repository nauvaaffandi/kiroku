<template>
  <div class="analytics-spending-view">
    <div class="page-header">
      <h1>{{ t('analytics.spending.title') }}</h1>
      <p class="page-subtitle">{{ t('analytics.spending.subtitle') }}</p>
    </div>

    <DateFilterBar @change="onDateChange" />

    <div class="widgets-grid">
      <SpendingByCategoryWidget
        :loading="categoryLoading"
        :data="categoryData"
        v-model:groupBy="categoryGroupBy"
        @update:groupBy="fetchSpendingByCategory"
      />

      <SpendingByLabelWidget
        :loading="labelLoading"
        :data="labelData"
      />

      <SpendingByNatureWidget
        :loading="natureLoading"
        :data="natureData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from '@/locales'
import { localDb } from '@/storage/repository'
import { useCategoriesStore } from '@/stores/categories'
import DateFilterBar from '@/components/common/DateFilterBar.vue'
import SpendingByCategoryWidget from '@/widgets/analytics/SpendingByCategoryWidget.vue'
import SpendingByLabelWidget from '@/widgets/analytics/SpendingByLabelWidget.vue'
import SpendingByNatureWidget from '@/widgets/analytics/SpendingByNatureWidget.vue'

const { t } = useI18n()

const dateFrom = ref('')
const dateTo = ref('')
const categoryGroupBy = ref('category')

const categoryLoading = ref(true)
const labelLoading = ref(true)
const natureLoading = ref(true)

const categoryData = ref<any[]>([])
const labelData = ref<any[]>([])
const natureData = ref<Record<string, any>>({})

function onDateChange(payload: { dateFrom: string; dateTo: string }) {
  dateFrom.value = payload.dateFrom
  dateTo.value = payload.dateTo
  fetchAll()
}

async function fetchAll() {
  await Promise.all([
    fetchSpendingByCategory(),
    fetchSpendingByLabel(),
    fetchSpendingByNature(),
  ])
}

async function fetchSpendingByCategory() {
  categoryLoading.value = true
  try {
    const categoriesStore = useCategoriesStore()
    await categoriesStore.fetchCategories()
    const allTransactions = await localDb.list<any>('transaction')
    const expenses = allTransactions.filter((t) => t.type === 'expense')
    const cats: any[] = (categoriesStore.categories as any)?.value ?? []
    const byCategory = new Map<string, any>()
    let totalExpense = 0
    for (const t of expenses) {
      const catId = t.category_id ?? 'uncategorized'
      const cat = (cats as any[]).find((c: any) => c.id === catId) ?? { name: 'Uncategorized', icon: '', color: '' }
      const entry = byCategory.get(catId) ?? {
        name: cat.name,
        icon: cat.icon,
        color: cat.color,
        amount: 0,
        count: 0,
        percentage: 0,
        children: []
      }
      entry.amount += t.amount
      entry.count += 1
      byCategory.set(catId, entry)
      totalExpense += t.amount
    }
    const groups = [...byCategory.entries()].map(([category_id, entry]) => ({
      id: category_id,
      category_id,
      category_name: entry.name,
      category_icon: entry.icon,
      category_color: entry.color,
      kind: entry.kind ?? 'expense',
      amount: entry.amount,
      count: entry.count,
      percentage: totalExpense ? (entry.amount / totalExpense) * 100 : 0,
      children: []
    }))
    categoryData.value = groups
  } catch (e) {
    console.error('Failed to fetch spending by category:', e)
  } finally {
    categoryLoading.value = false
  }
}

async function fetchSpendingByLabel() {
  labelLoading.value = true
  try {
    const allTransactions = await localDb.list<any>('transaction')
    const labelsMap = new Map<string, any>()
    let totalExpense = 0
    for (const t of allTransactions) {
      if (t.type !== 'expense') continue
      const labels = Array.isArray(t.labels) ? t.labels : String(t.labels ?? '').split(',')
      const valid = labels.filter((l: string) => l?.trim())
      if (valid.length) {
        for (const label of valid) {
          const entry = labelsMap.get(label) ?? { label, amount: 0, count: 0 }
          entry.amount += t.amount
          entry.count += 1
          labelsMap.set(label, entry)
          totalExpense += t.amount
        }
      } else {
        const entry = labelsMap.get('uncategorized') ?? { label: 'Uncategorized', amount: 0, count: 0 }
        entry.amount += t.amount
        entry.count += 1
        labelsMap.set('uncategorized', entry)
        totalExpense += t.amount
      }
    }
    const items = [...labelsMap.entries()].map(([label, entry]) => ({
      label,
      amount: entry.amount,
      count: entry.count,
      percentage: totalExpense ? (entry.amount / totalExpense) * 100 : 0
    }))
    labelData.value = items.slice(0, 10)
  } catch (e) {
    console.error('Failed to fetch spending by label:', e)
  } finally {
    labelLoading.value = false
  }
}

async function fetchSpendingByNature() {
  natureLoading.value = true
  try {
    const categories = await localDb.list<any>('category')
    const allTransactions = await localDb.list<any>('transaction')
    const nature: Record<string, { amount: number; count: number }> = {
      must: { amount: 0, count: 0 },
      have: { amount: 0, count: 0 },
      want: { amount: 0, count: 0 },
      unclassified: { amount: 0, count: 0 }
    }
    for (const t of allTransactions) {
      if (t.type === 'expense') continue
      const cat = categories.find((c) => c.id === t.category_id)
      const kind = cat?.kind ?? 'uncategorized'
      if (!['must', 'have', 'want', 'uncategorized'].includes(kind)) continue
      nature[kind].amount += t.amount
      nature[kind].count += 1
    }
    natureData.value = nature
  } catch (e) {
    console.error('Failed to fetch spending by nature:', e)
  } finally {
    natureLoading.value = false
  }
}

onMounted(async () => {
  const now = new Date()
  dateTo.value = now.toISOString().split('T')[0]
  const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1)
  dateFrom.value = lastMonth.toISOString().split('T')[0]
  await fetchAll()
})
</script>

<style scoped>
@reference "tailwindcss";
.analytics-spending-view {
  @apply max-w-[80rem] mx-auto p-4 flex flex-col gap-4;
}
.page-header { 
  @apply flex flex-col gap-1; 
}
.page-header h1 { 
  @apply text-xl font-bold; 
  color: var(--text-primary); 
}
.page-subtitle { 
  color: var(--text-muted); 
}
.widgets-grid { 
  @apply grid gap-4; 
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); 
}
</style>