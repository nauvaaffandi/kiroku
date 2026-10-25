<template>
  <section class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
    <div class="bg-white dark:bg-zinc-900 rounded-xl shadow-sm p-5 flex items-center gap-4 border border-transparent dark:border-zinc-800">
      <div class="w-12 h-12 rounded-lg bg-emerald-50 dark:bg-emerald-950/60 text-emerald-600 dark:text-emerald-300 flex items-center justify-center shrink-0">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <line x1="12" y1="19" x2="12" y2="5"></line>
          <polyline points="19 12 12 5 5 12"></polyline>
        </svg>
      </div>
      <div>
        <p class="text-xs font-semibold tracking-wide uppercase text-zinc-500 dark:text-zinc-400">{{ t('dashboard.monthlyIncome') }}</p>
        <p class="text-xl font-black text-emerald-600 dark:text-emerald-400 mt-0.5">{{ formatCurrency(monthlyIncome) }}</p>
      </div>
    </div>

    <div class="bg-white dark:bg-zinc-900 rounded-xl shadow-sm p-5 flex items-center gap-4 border border-transparent dark:border-zinc-800">
      <div class="w-12 h-12 rounded-lg bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-300 flex items-center justify-center shrink-0">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <polyline points="5 12 12 19 19 12"></polyline>
        </svg>
      </div>
      <div>
        <p class="text-xs font-semibold tracking-wide uppercase text-zinc-500 dark:text-zinc-400">{{ t('dashboard.monthlyExpense') }}</p>
        <p class="text-xl font-black text-rose-600 dark:text-rose-400 mt-0.5">{{ formatCurrency(monthlyExpense) }}</p>
      </div>
    </div>

    <div class="bg-white dark:bg-zinc-900 rounded-xl shadow-sm p-5 flex items-center gap-4 border border-transparent dark:border-zinc-800">
      <div class="w-12 h-12 rounded-lg bg-blue-50 dark:bg-blue-950/60 text-blue-600 dark:text-blue-300 flex items-center justify-center shrink-0">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <line x1="18" y1="20" x2="6" y2="20"></line>
          <line x1="12" y1="16" x2="12" y2="8"></line>
          <line x1="6" y1="12" x2="18" y2="12"></line>
        </svg>
      </div>
      <div>
        <p class="text-xs font-semibold tracking-wide uppercase text-zinc-500 dark:text-zinc-400">{{ t('dashboard.net') }}</p>
        <p :class="['text-xl font-black mt-0.5', monthlyNet >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600 dark:text-rose-400']">
          {{ formatCurrency(monthlyNet) }}
        </p>
      </div>
    </div>

    <div class="bg-white dark:bg-zinc-900 rounded-xl shadow-sm p-5 flex items-center gap-4 border border-transparent dark:border-zinc-800">
      <div class="w-12 h-12 rounded-lg bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-300 flex items-center justify-center shrink-0">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <rect x="2" y="4" width="20" height="16" rx="2"></rect>
          <line x1="2" y1="10" x2="22" y2="10"></line>
        </svg>
      </div>
      <div>
        <p class="text-xs font-semibold tracking-wide uppercase text-zinc-500 dark:text-zinc-400">{{ t('dashboard.totalBalance') }}</p>
        <p class="text-xl font-black text-zinc-900 dark:text-zinc-50 mt-0.5">{{ formatCurrency(totalBalance) }}</p>
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