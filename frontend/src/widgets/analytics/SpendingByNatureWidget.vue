<template>
  <div class="widget-card">
    <div class="widget-header">
      <h2>支出の性質</h2>
    </div>
    <div class="widget-content" v-if="!loading">
      <div class="nature-charts">
        <div class="nature-donut" v-for="nature in ['must', 'have', 'want']" :key="nature">
          <div class="donut-container" ref="donutRefs[nature]" />
          <div class="nature-summary">
            <h3 class="nature-title" :class="nature">{{ natureLabels[nature] }}</h3>
            <p class="nature-amount">{{ formatCurrencyShort(data[nature]?.amount || 0) }}</p>
            <p class="nature-pct">{{ data[nature]?.percentage?.toFixed(1) || 0 }}%</p>
            <p class="nature-count">{{ data[nature]?.count || 0 }} 件</p>
          </div>
        </div>
        <div class="nature-donut unclassified">
          <div class="donut-container" ref="donutRefs.unclassified" />
          <div class="nature-summary">
            <h3 class="nature-title unclassified">未分類</h3>
            <p class="nature-amount">{{ formatCurrencyShort(data.unclassified?.amount || 0) }}</p>
            <p class="nature-pct">{{ data.unclassified?.percentage?.toFixed(1) || 0 }}%</p>
            <p class="nature-count">{{ data.unclassified?.count || 0 }} 件</p>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="widget-content">
      <div class="skeleton skeleton-chart" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'

interface NatureSpending {
  amount: number
  count: number
  percentage: number
}

const props = defineProps<{
  loading: boolean
  data: Record<string, NatureSpending>
}>()

const donutRefs = ref<Record<string, HTMLElement | null>>({
  must: null,
  have: null,
  want: null,
  unclassified: null,
})
const donutCharts: Record<string, echarts.ECharts | null> = {
  must: null,
  have: null,
  want: null,
  unclassified: null,
}

const natureLabels: Record<string, string> = {
  must: 'Must (必需)',
  have: 'Have (必要)',
  want: 'Want (欲求)',
  unclassified: '未分類',
}

const natureColors: Record<string, string> = {
  must: '#28a745',
  have: '#ffc107',
  want: '#dc3545',
  unclassified: '#6c757d',
}

function renderDonutCharts() {
  Object.keys(donutRefs.value).forEach(nature => {
    const el = donutRefs.value[nature]
    if (!el) return
    if (donutCharts[nature]) donutCharts[nature]?.dispose()

    const data = props.data[nature]
    if (!data) return

    donutCharts[nature] = echarts.init(el)
    const amount = data.amount || 0
    const pct = data.percentage || 0

    donutCharts[nature].setOption({
      series: [
        {
          type: 'pie',
          radius: ['60%', '80%'],
          avoidLabelOverlap: false,
          data: [
            { value: amount, name: '支出', itemStyle: { color: natureColors[nature] } },
            { value: 100 - pct, name: '', itemStyle: { color: '#eee' }, label: { show: false } },
          ],
          label: { show: false },
        },
      ],
      graphic: {
        type: 'text',
        left: 'center',
        top: 'center',
        style: {
          text: pct.toFixed(0) + '%',
          fontSize: 18,
          fontWeight: 'bold',
          fill: natureColors[nature],
        },
      },
    })
  })
}

function formatCurrencyShort(val: number): string {
  if (Math.abs(val) >= 1e8) return (val / 1e8).toFixed(1) + '億'
  if (Math.abs(val) >= 1e4) return (val / 1e4).toFixed(1) + '万'
  return val.toLocaleString('ja-JP')
}

watch(() => props.data, () => {
  nextTick(() => renderDonutCharts())
}, { deep: true })

onMounted(() => {
  nextTick(() => renderDonutCharts())
  window.addEventListener('resize', () => {
    Object.values(donutCharts).forEach(c => c?.resize())
  })
})
</script>

<style scoped>
@reference "tailwindcss";
.widget-card { background-color: var(--bg-card); border: 1px solid var(--border-primary); border-radius: var(--radius); overflow: hidden; box-shadow: var(--shadow); }
.widget-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem; border-bottom: 1px solid var(--border-primary); background-color: var(--bg-secondary); }
.widget-header h2 { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
.widget-content { padding: 1rem; }
.nature-charts { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; }
.nature-donut { display: flex; flex-direction: column; align-items: center; gap: 0.5rem; }
.donut-container { width: 120px; height: 120px; }
.nature-summary { text-align: center; }
.nature-title { font-size: 0.875rem; font-weight: 600; margin: 0 0 0.25rem; }
.nature-title.must { color: #28a745; }
.nature-title.have { color: #ffc107; }
.nature-title.want { color: #dc3545; }
.nature-title.unclassified { color: #6c757d; }
.nature-amount { font-size: 1.125rem; font-weight: 700; color: var(--text-primary); margin: 0; }
.nature-pct { font-size: 0.875rem; color: var(--text-secondary); margin: 0; }
.nature-count { font-size: 0.75rem; color: var(--text-muted); margin: 0.25rem 0 0; }
</style>