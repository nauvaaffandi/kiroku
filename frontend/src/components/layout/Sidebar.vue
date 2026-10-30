<template>
  <div class="sidebar">
    <div class="sidebar-header">
      <router-link to="/dashboard" class="logo">
        <AppIcon class="app-icon-logo" />
        <span>Kiroku</span>
      </router-link>
      <button v-show="isMobile" @click="$emit('close')" class="sidebar-close text-zinc-600 dark:text-white hover:text-zinc-900 dark:hover:text-white" aria-label="Close sidebar">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>
    
    <nav class="sidebar-nav">
      <ul class="nav-list">
        <li v-for="item in mainNavItems" :key="item.path">
          <router-link
            :to="item.path"
            :class="['nav-link', { active: isActive(item.path) }, 'text-zinc-700 dark:text-white font-medium font-semibold']"
            @click="closeMobileSidebar"
          >
            <component :is="item.icon" class="nav-icon" :style="{ color: item.color }" />
            <span>{{ item.label }}</span>
          </router-link>
        </li>

        <li class="nav-group">
          <button
            @click="toggleAnalytics"
            :class="['nav-link nav-group-toggle w-full text-zinc-700 dark:text-white font-medium font-semibold', { active: isAnalyticsActive && !isAnalyticsOpen }]"
          >
            <AnalyticsIcon class="nav-icon" style="color:#6366f1" />
            <span class="flex-1 text-left">{{ t('nav.analytics') }}</span>
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              :class="['transition-transform duration-200', { 'rotate-180': isAnalyticsOpen }]"
            >
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </button>
          
          <div class="expandable-wrapper" :style="{ height: isAnalyticsOpen ? analyticsHeight + 'px' : '0px' }" ref="expandableRef">
            <ul class="sub-nav-list" ref="subNavRef">
              <li v-for="sub in analyticsItems" :key="sub.path">
                <router-link
                  :to="sub.path"
                  :class="['sub-nav-link', { active: route.path === sub.path }, 'text-zinc-600 dark:text-zinc-300 font-medium']"
                  @click="closeMobileSidebar"
                >
                  <component :is="sub.icon" class="sub-nav-icon" v-if="sub.icon" :style="{ color: sub.color }" />
                  <span>{{ sub.label }}</span>
                </router-link>
              </li>
            </ul>
          </div>
        </li>

        <li v-for="item in bottomMainNavItems" :key="item.path">
          <router-link
            :to="item.path"
            :class="['nav-link', { active: isActive(item.path) }, 'text-zinc-700 dark:text-white font-medium font-semibold']"
            @click="closeMobileSidebar"
          >
            <component :is="item.icon" class="nav-icon" :style="{ color: item.color }" />
            <span>{{ item.label }}</span>
          </router-link>
        </li>
      </ul>
      
      <div class="nav-divider" />
      
      <ul class="nav-list">
        <li v-for="item in settingsItems" :key="item.path">
          <router-link
            :to="item.path"
            :class="['nav-link', { active: isActive(item.path) }, 'text-zinc-700 dark:text-white font-medium font-semibold']"
            @click="closeMobileSidebar"
          >
            <span>{{ item.label }}</span>
          </router-link>
        </li>
      </ul>
    </nav>
    
      <div class="sidebar-footer">
      <div class="user-info" v-if="authStore.user">
        <div class="user-avatar">
          {{ authStore.user.user_id.charAt(0).toUpperCase() }}
        </div>
        <div class="user-details">
          <span class="user-name text-zinc-900 dark:text-white font-bold">{{ authStore.user.user_id.slice(0, 8) }}</span>
          <span class="user-role text-zinc-500 dark:text-zinc-300 font-medium">User</span>
        </div>
      </div>
      <router-link to="/settings" class="footer-link text-zinc-600 dark:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800 font-semibold" @click="closeMobileSidebar">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="color: #64748b">
          <circle cx="12" cy="12" r="3"></circle>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06-.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
        </svg>
        <span>Settings</span>
      </router-link>
      <ThemeToggle />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import { useI18n } from '@/locales'

import DashboardIcon from '@/components/icons/DashboardIcon.vue'
import TransactionsIcon from '@/components/icons/TransactionsIcon.vue'
import CalendarIcon from '@/components/icons/CalendarIcon.vue'
import AccountsIcon from '@/components/icons/AccountsIcon.vue'
import CategoriesIcon from '@/components/icons/CategoriesIcon.vue'
import LabelsIcon from '@/components/icons/LabelsIcon.vue'
import ScheduledIcon from '@/components/icons/ScheduledIcon.vue'
import ReportsIcon from '@/components/icons/ReportsIcon.vue'
import AnalyticsIcon from '@/components/icons/AnalyticsIcon.vue'
import SharedAccessIcon from '@/components/icons/SharedAccessIcon.vue'
import BudgetIcon from '@/components/icons/BudgetIcon.vue'
import GoalsIcon from '@/components/icons/GoalsIcon.vue'
import AppIcon from '@/components/icons/AppIcon.vue'
import ThemeToggle from '@/components/common/ThemeToggle.vue'

const route = useRoute()
const authStore = useAuthStore()
const uiStore = useUIStore()
const { t } = useI18n()

const isMobile = computed(() => window.innerWidth < 768)

const mainNavItems = [
  { path: '/dashboard', label: t('nav.dashboard'), icon: DashboardIcon, color: '#3b82f6' }, // blue-500
  { path: '/transactions', label: t('nav.transactions'), icon: TransactionsIcon, color: '#10b981' }, // emerald-500
  { path: '/calendar', label: t('nav.calendar'), icon: CalendarIcon, color: '#f59e0b' } // amber-500
]

const analyticsItems = [
  { path: '/analytics/balance', label: t('nav.analyticsBalance'), color: '#6366f1' }, // indigo-500
  { path: '/analytics/outlook', label: t('nav.analyticsOutlook'), color: '#8b5cf6' }, // violet-500
  { path: '/analytics/cashflow', label: t('nav.analyticsCashflow'), color: '#ec4899' }, // pink-500
  { path: '/analytics/spending', label: t('nav.analyticsSpending'), color: '#f43f5e' }, // rose-500
  { path: '/analytics/budget', label: t('nav.analyticsBudget'), icon: BudgetIcon, color: '#14b8a6' }, // teal-500
  { path: '/analytics/goals', label: t('nav.analyticsGoals'), icon: GoalsIcon, color: '#eab308' } // yellow-500
]

const bottomMainNavItems = [
  { path: '/accounts', label: t('nav.accounts'), icon: AccountsIcon, color: '#06b6d4' }, // cyan-500
  { path: '/categories', label: t('nav.categories'), icon: CategoriesIcon, color: '#84cc16' }, // lime-500
  { path: '/labels', label: t('nav.labels'), icon: LabelsIcon, color: '#6366f1' }, // indigo-500
  { path: '/scheduled', label: t('nav.scheduled'), icon: ScheduledIcon, color: '#0ea5e9' }, // sky-500
  { path: '/reports', label: t('nav.reports'), icon: ReportsIcon, color: '#a855f7' }, // purple-500
  { path: '/shared-access', label: t('nav.sharedAccess'), icon: SharedAccessIcon, color: '#f97316' } // orange-500
]

const settingsItems = [
  { path: '/settings', label: t('nav.settings') }
]

const isAnalyticsActive = computed(() => route.path.startsWith('/analytics'))
const isAnalyticsOpen = ref(isAnalyticsActive.value)
const analyticsHeight = ref(0)
const subNavRef = ref<HTMLElement | null>(null)

const updateHeight = () => {
  if (subNavRef.value) {
    analyticsHeight.value = subNavRef.value.scrollHeight
  }
}

watch(isAnalyticsActive, (val) => {
  if (val) {
    isAnalyticsOpen.value = true
  }
})

watch(isAnalyticsOpen, () => {
  nextTick(() => {
    updateHeight()
  })
})

onMounted(() => {
  if (isAnalyticsActive.value) {
    isAnalyticsOpen.value = true
  }
  nextTick(() => {
    updateHeight()
  })
})

const toggleAnalytics = () => {
  isAnalyticsOpen.value = !isAnalyticsOpen.value
}

const isActive = (path: string) => {
  return route.path.startsWith(path)
}

const closeMobileSidebar = () => {
  if (isMobile.value) {
    uiStore.sidebarOpen = false
  }
}

defineEmits(['close'])
</script>

<style scoped>
@reference "tailwindcss";
.sidebar {
  @apply flex flex-col h-full;
}

.sidebar-header {
  @apply flex items-center justify-between px-6 py-4 border-b;
  border-color: var(--border-primary);
}

.logo {
  @apply flex items-center gap-3 font-bold text-lg no-underline;
  color: var(--accent);
}

.app-icon-logo {
  width: 2rem;
  height: 2rem;
  flex-shrink: 0;
  color: var(--accent);
}

.sidebar-close {
  @apply p-2 rounded-sm cursor-pointer transition-colors;
  color: var(--text-muted);
}

.sidebar-close:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.sidebar-nav {
  @apply flex-1 overflow-y-auto px-2 py-2;
}

.nav-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-link {
  @apply flex items-center gap-3 px-4 py-3 rounded-sm no-underline mx-[0.125rem] my-[0.125rem] w-auto bg-transparent border-0 cursor-pointer text-left transition-colors;
  color: var(--text-secondary);
  font: inherit;
}

.nav-link:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.nav-link.active {
  background-color: var(--accent-light);
  color: var(--accent);
  font-weight: 500;
}

.nav-group-toggle {
  @apply mx-[0.125rem] my-[0.125rem] w-[calc(100%-1rem)];
}

.expandable-wrapper {
  overflow: hidden;
  transition: height 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  @apply mx-[0.5rem];
}

.sub-nav-list {
  list-style: none;
  padding: 0.25rem 0 0.25rem 2rem;
  margin: 0;
  @apply flex flex-col gap-0.5;
}

.sub-nav-link {
  @apply flex items-center gap-2 px-3 py-2 rounded-sm text-sm transition-colors no-underline;
  color: var(--text-secondary);
}

.sub-nav-link:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.sub-nav-link.active {
  background-color: var(--accent-light);
  color: var(--accent);
  font-weight: 600;
}

.nav-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
  color: currentColor;
}

.sub-nav-icon {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
  color: currentColor;
}

.nav-divider {
  height: 1px;
  background-color: var(--border-primary);
  @apply my-[1rem] mx-[0.5rem];
}

.sidebar-footer {
  @apply px-6 py-4 border-t;
  border-color: var(--border-primary);
}

.user-info {
  @apply flex items-center gap-3 p-2 rounded-sm mb-4;
  background-color: var(--bg-card);
}

.user-avatar {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  background-color: var(--accent);
  color: white;
  @apply flex items-center justify-center font-semibold text-sm;
}

.user-details {
  @apply flex-1 min-w-0 flex flex-col;
}

.user-name {
  @apply text-sm font-medium truncate;
  color: var(--text-primary);
}

.user-role {
  @apply text-xs;
  color: var(--text-muted);
}

.footer-link {
  @apply flex items-center gap-3 px-4 py-3 rounded-sm no-underline transition-colors;
  color: var(--text-secondary);
}

.footer-link:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.footer-link svg {
  flex-shrink: 0;
}
</style>
