<template>
  <div class="analytics-spending-view">
    <div class="page-header">
      <h1>支出分析</h1>
      <p class="page-subtitle">カテゴリ別・ラベル別・性質別（Must/Have/Want）の支出内訳</p>
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
import { api } from '@/api/client'
import DateFilterBar from '@/components/common/DateFilterBar.vue'
import SpendingByCategoryWidget from '@/widgets/analytics/SpendingByCategoryWidget.vue'
import SpendingByLabelWidget from '@/widgets/analytics/SpendingByLabelWidget.vue'
import SpendingByNatureWidget from '@/widgets/analytics/SpendingByNatureWidget.vue'

interface CategorySpending {
  id: string
  name: string
  parent_id: string
  color: string
  kind: string
  amount: number
  count: number
  percentage: number
  children: any[]
}

interface LabelSpending {
  label: string
  amount: number
  count: number
  percentage: number
}

interface NatureSpending {
  amount: number
  count: number
  percentage: number
}

const dateFrom = ref('')
const dateTo = ref('')
const categoryGroupBy = ref('category')

const categoryLoading = ref(true)
const labelLoading = ref(true)
const natureLoading = ref(true)

const categoryData = ref<CategorySpending[]>([])
const labelData = ref<LabelSpending[]>([])
const natureData = ref<Record<string, NatureSpending>>({})

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
    const res = await api.get<{ categories: CategorySpending[]; total: number }>(
      `/analytics/spending/by-category?date_from=${dateFrom.value}&date_to=${dateTo.value}&group_by=${categoryGroupBy.value}`
    )
    categoryData.value = res.data.categories
  } catch (e) {
    console.error('Failed to fetch spending by category:', e)
  } finally {
    categoryLoading.value = false
  }
}

async function fetchSpendingByLabel() {
  labelLoading.value = true
  try {
    const res = await api.get<{ labels: LabelSpending[]; total: number }>(
      `/analytics/spending/by-label?date_from=${dateFrom.value}&date_to=${dateTo.value}`
    )
    labelData.value = res.data.labels.slice(0, 10)
  } catch (e) {
    console.error('Failed to fetch spending by label:', e)
  } finally {
    labelLoading.value = false
  }
}

async function fetchSpendingByNature() {
  natureLoading.value = true
  try {
    const res = await api.get<{
      must: NatureSpending
      have: NatureSpending
      want: NatureSpending
      unclassified: NatureSpending
      total: number
    }>(
      `/analytics/spending/by-nature?date_from=${dateFrom.value}&date_to=${dateTo.value}`
    )
    natureData.value = {
      must: res.data.must,
      have: res.data.have,
      want: res.data.want,
      unclassified: res.data.unclassified,
    }
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