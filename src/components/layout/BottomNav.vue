<template>
  <nav class="bottom-nav">
    <template v-for="(item, index) in navItems" :key="item.path">
      <router-link
        :to="item.path"
        :class="['bottom-nav-item', { active: isActive(item.path) }]"
        :aria-label="t(item.labelKey)"
      >
        <component :is="item.icon" class="bottom-nav-icon" />
        <span class="bottom-nav-label">{{ t(item.labelKey) }}</span>
      </router-link>

      <button
        v-if="index === 1"
        @click="handleFABClick"
        class="bottom-add-btn"
        :aria-label="t('common.addTransaction')"
      >
        <div class="bottom-add-circle">
          <Plus class="bottom-add-icon" />
        </div>
      </button>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from '@/locales'
import { Plus } from '@lucide/vue'
import DashboardIcon from '@/components/icons/DashboardIcon.vue'
import TransactionsIcon from '@/components/icons/TransactionsIcon.vue'
import AccountsIcon from '@/components/icons/AccountsIcon.vue'
import ReportsIcon from '@/components/icons/ReportsIcon.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const navItems = [
  { path: '/', labelKey: 'nav.dashboard', icon: DashboardIcon },
  { path: '/transactions', labelKey: 'nav.transactions', icon: TransactionsIcon },
  { path: '/accounts', labelKey: 'nav.accounts', icon: AccountsIcon },
  { path: '/reports', labelKey: 'nav.reports', icon: ReportsIcon }
]

const isActive = (path: string): boolean => {
  // '/' must match exactly (it matches every route under startsWith)
  if (path === '/') {
    return route.path === '/'
  }
  return route.path.startsWith(path)
}

const handleFABClick = () => {
  if (navigator.vibrate) {
    navigator.vibrate(10)
  }
  if (route.path === '/transactions') {
    router.push({ name: 'Transactions', query: { create: 'true' } })
  } else {
    router.push('/transactions?create=true')
  }
}
</script>

<style scoped>
@reference "tailwindcss";

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  width: 100%;
  height: 60px;
  @apply flex items-center justify-around z-[100] pb-[env(safe-area-inset-bottom)];
  background-color: var(--bg-card);
  border-top: 1px solid var(--border-primary);
}

.bottom-nav-item {
  position: static;
  @apply flex flex-col items-center justify-center flex-1 h-full;
  color: var(--text-secondary);
  text-decoration: none;
}

.bottom-nav-item.active {
  color: var(--text-primary);
  border-top: 2px solid var(--accent);
}

.bottom-nav-icon {
  width: 1.25rem;
  height: 1.25rem;
  margin-bottom: 0.125rem;
}

.bottom-nav-label {
  @apply text-[11px] font-medium;
}

.bottom-add-btn {
  position: static;
  @apply flex items-center justify-center border-0 bg-transparent cursor-pointer;
}

.bottom-add-circle {
  width: 52px;
  height: 52px;
  border-radius: var(--radius-lg);
  background-color: var(--text-primary);
  color: var(--text-inverse);
  @apply flex items-center justify-center shadow-md;
}

.bottom-add-icon {
  width: 1.125rem;
  height: 1.125rem;
}
</style>
