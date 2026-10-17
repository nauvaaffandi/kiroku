<template>
  <section class="stats-grid">
    <div class="stat-card income">
      <div class="stat-icon">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="19" x2="12" y2="5"></line>
          <polyline points="19 12 12 5 5 12"></polyline>
        </svg>
      </div>
      <div class="stat-info">
        <p class="stat-label">{{ t('dashboard.monthlyIncome') }}</p>
        <p class="stat-value text-success">{{ formatCurrency(monthlyIncome) }}</p>
      </div>
    </div>

    <div class="stat-card expense">
      <div class="stat-icon">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <polyline points="5 12 12 19 19 12"></polyline>
        </svg>
      </div>
      <div class="stat-info">
        <p class="stat-label">{{ t('dashboard.monthlyExpense') }}</p>
        <p class="stat-value text-danger">{{ formatCurrency(monthlyExpense) }}</p>
      </div>
    </div>

    <div class="stat-card net">
      <div class="stat-icon">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="20" x2="6" y2="20"></line>
          <line x1="12" y1="16" x2="12" y2="8"></line>
          <line x1="6" y1="12" x2="18" y2="12"></line>
        </svg>
      </div>
      <div class="stat-info">
        <p class="stat-label">{{ t('dashboard.net') }}</p>
        <p :class="['stat-value', monthlyNet >= 0 ? 'text-success' : 'text-danger']">
          {{ formatCurrency(monthlyNet) }}
        </p>
      </div>
    </div>

    <div class="stat-card balance">
      <div class="stat-icon">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="2" y="4" width="20" height="16" rx="2"></rect>
          <line x1="2" y1="10" x2="22" y2="10"></line>
        </svg>
      </div>
      <div class="stat-info">
        <p class="stat-label">{{ t('dashboard.totalBalance') }}</p>
        <p class="stat-value">{{ formatCurrency(totalBalance) }}</p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'

const { t } = useI18n()

defineProps<{
  monthlyIncome: number
  monthlyExpense: number
  monthlyNet: number
  totalBalance: number
}>()

const formatCurrency = (num: number) => {
  return new Intl.NumberFormat('ja-JP').format(Math.abs(num))
}
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}
.stat-card {
  background-color: var(--color-neutral-100);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 1.25rem;
  display: flex;
  align-items: center;
  gap: 1rem;
}
.stat-icon {
  width: 3rem;
  height: 3rem;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
}
.stat-card.income .stat-icon { background-color: rgba(40, 167, 69, 0.15); color: var(--color-success); }
.stat-card.expense .stat-icon { background-color: rgba(220, 53, 69, 0.15); color: var(--color-danger); }
.stat-card.net .stat-icon { background-color: rgba(23, 162, 184, 0.15); color: var(--color-info); }
.stat-card.balance .stat-icon { background-color: rgba(74, 144, 217, 0.15); color: var(--color-primary); }
.stat-label { font-size: 0.875rem; color: var(--color-neutral-500); margin-bottom: 0.25rem; }
.stat-value { font-size: 1.25rem; font-weight: 600; color: var(--color-neutral-900); }
.text-success { color: var(--color-success); }
.text-danger { color: var(--color-danger); }
</style>