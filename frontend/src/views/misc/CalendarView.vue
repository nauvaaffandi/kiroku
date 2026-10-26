<template>
  <div class="calendar-view">
    <CalendarHeaderWidget
      :year="year"
      :month="month"
      v-model:selectedAccount="selectedAccount"
      :accounts="accounts"
      @prev="prevMonth"
      @next="nextMonth"
      @today="goToday"
      @update:selectedAccount="fetchCalendar"
    />

    <CalendarGridWidget
      v-if="!loading"
      :calendarDays="calendarDays"
      @open-detail="openDayDetail"
    />

    <div v-else class="calendar-loading">
      <div class="skeleton skeleton-card" style="height: 600px;" />
    </div>

    <DayDetailModalWidget
      :activeDay="activeDay"
      @close="activeDay = null"
      @add-transaction="openAddTransaction"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountsStore } from '@/stores/accounts'
import { api } from '@/api/client'
import CalendarHeaderWidget from '@/widgets/calendar/CalendarHeaderWidget.vue'
import CalendarGridWidget from '@/widgets/calendar/CalendarGridWidget.vue'
import DayDetailModalWidget from '@/widgets/calendar/DayDetailModalWidget.vue'
import type { CalendarDayBackend, DisplayDay } from '@/types'

const router = useRouter()
const accountsStore = useAccountsStore()

const loading = ref(true)
const currentDate = ref(new Date())
const selectedAccount = ref('')
const calendarData = ref<Map<string, CalendarDayBackend>>(new Map())
const activeDay = ref<DisplayDay | null>(null)

const year = computed(() => currentDate.value.getFullYear())
const month = computed(() => currentDate.value.getMonth() + 1)
const accounts = computed(() => accountsStore.accounts)

const calendarDays = computed(() => {
  const y = year.value
  const m = month.value
  const firstDayOfMonth = new Date(y, m - 1, 1)
  const lastDayOfMonth = new Date(y, m, 0)
  
  const days: DisplayDay[] = []
  const todayStr = new Date().toISOString().split('T')[0]

  const firstWeekday = firstDayOfMonth.getDay()
  for (let i = firstWeekday - 1; i >= 0; i--) {
    const d = new Date(y, m - 1, -i)
    const dateStr = d.toISOString().split('T')[0]
    const backendDay = calendarData.value.get(dateStr)
    days.push({
      date: dateStr,
      dayNumber: d.getDate(),
      isCurrentMonth: false,
      isToday: dateStr === todayStr,
      transactions: backendDay?.transactions || [],
      scheduled: backendDay?.scheduled || [],
      projectedBalance: backendDay?.projected_balance || 0
    })
  }

  for (let i = 1; i <= lastDayOfMonth.getDate(); i++) {
    const d = new Date(y, m - 1, i)
    const dateStr = d.toISOString().split('T')[0]
    const backendDay = calendarData.value.get(dateStr)
    days.push({
      date: dateStr,
      dayNumber: i,
      isCurrentMonth: true,
      isToday: dateStr === todayStr,
      transactions: backendDay?.transactions || [],
      scheduled: backendDay?.scheduled || [],
      projectedBalance: backendDay?.projected_balance || 0
    })
  }

  const remainingCells = 42 - days.length
  for (let i = 1; i <= remainingCells; i++) {
    const d = new Date(y, m, i)
    const dateStr = d.toISOString().split('T')[0]
    const backendDay = calendarData.value.get(dateStr)
    days.push({
      date: dateStr,
      dayNumber: d.getDate(),
      isCurrentMonth: false,
      isToday: dateStr === todayStr,
      transactions: backendDay?.transactions || [],
      scheduled: backendDay?.scheduled || [],
      projectedBalance: backendDay?.projected_balance || 0
    })
  }

  return days
})

function prevMonth() {
  currentDate.value = new Date(year.value, month.value - 2, 1)
  fetchCalendar()
}

function nextMonth() {
  currentDate.value = new Date(year.value, month.value, 1)
  fetchCalendar()
}

function goToday() {
  currentDate.value = new Date()
  fetchCalendar()
}

async function fetchCalendar() {
  loading.value = true
  try {
    const params = new URLSearchParams({
      year: year.value.toString(),
      month: month.value.toString()
    })
    if (selectedAccount.value) params.append('account_id', selectedAccount.value)

    const res = await api.get<{ days: CalendarDayBackend[] }>(`/api/v1/analytics/calendar?${params}`)
    const map = new Map<string, CalendarDayBackend>()
    res.data.days.forEach(d => map.set(d.date, d))
    calendarData.value = map
  } catch (e) {
    console.error('Failed to fetch calendar:', e)
  } finally {
    loading.value = false
  }
}

function openDayDetail(day: DisplayDay) {
  activeDay.value = day
}

function openAddTransaction(date: string) {
  router.push({ path: '/transactions', query: { action: 'add', date } })
}

onMounted(async () => {
  await accountsStore.fetchAccounts()
  await fetchCalendar()
})
</script>

<style scoped>
@reference "tailwindcss";
.calendar-view {
  @apply mx-auto p-4 flex flex-col gap-4 max-w-[80rem];
}
</style>