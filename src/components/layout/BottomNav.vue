<template>
  <nav class="bottom-nav">
    <template v-for="(item, index) in navItems" :key="item.path">
      <router-link
        :to="item.path"
        :class="['bottom-nav-item', { active: isActive(item.path) }]"
      >
        <component :is="item.icon" class="bottom-nav-icon" />
        <span class="bottom-nav-label">{{ item.label }}</span>
      </router-link>

      <button v-if="index === 1" @click="handleFABClick" class="bottom-add-btn" aria-label="Add transaction">
        <div class="bottom-add-circle">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
        </div>
      </button>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import DashboardIcon from '@/components/icons/DashboardIcon.vue'
import TransactionsIcon from '@/components/icons/TransactionsIcon.vue'
import AccountsIcon from '@/components/icons/AccountsIcon.vue'
import ReportsIcon from '@/components/icons/ReportsIcon.vue'

const route = useRoute()
const router = useRouter()

const navItems = [
  { path: '/dashboard', label: 'Home', icon: DashboardIcon },
  { path: '/transactions', label: 'Txns', icon: TransactionsIcon },
  { path: '/accounts', label: 'Accts', icon: AccountsIcon },
  { path: '/reports', label: 'Reports', icon: ReportsIcon }
]

const isActive = (path: string) => {
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
  color: var(--accent);
}

.bottom-nav-icon {
  width: 1.25rem;
  height: 1.25rem;
  margin-bottom: 0.125rem;
}

.bottom-nav-label {
  font-size: 0.6875rem;
  font-weight: 600;
}

.bottom-add-btn {
  position: static;
  @apply flex items-center justify-center flex-1 h-full border-0 bg-transparent cursor-pointer;
}

.bottom-add-circle {
  width: 2.75rem;
  height: 2.75rem;
  border-radius: 50%;
  background-color: var(--accent);
  color: white;
  @apply flex items-center justify-center shadow-md;
}
</style>