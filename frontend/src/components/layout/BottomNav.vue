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

      <!-- Insert the Add Transaction button right in the exact center (after index 1) -->
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
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  width: 100%;
  height: 60px;
  background-color: var(--color-neutral-100);
  border-top: 1px solid var(--color-neutral-200);
  display: flex;
  align-items: center;
  justify-content: space-around;
  z-index: 100;
  padding-bottom: env(safe-area-inset-bottom);
}

:where([data-theme="dark"], [data-theme="dark"] *) .bottom-nav {
  background-color: #18181b;
  border-top-color: #27272a;
}

.bottom-nav-item {
  position: static;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  height: 100%;
  color: var(--color-neutral-500);
  text-decoration: none;
}

:where([data-theme="dark"], [data-theme="dark"] *) .bottom-nav-item {
  color: #a1a1aa;
}

.bottom-nav-item.active {
  color: var(--color-primary);
}

:where([data-theme="dark"], [data-theme="dark"] *) .bottom-nav-item.active {
  color: var(--color-primary);
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
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  height: 100%;
  background: none;
  border: none;
  cursor: pointer;
}

.bottom-add-circle {
  width: 2.75rem;
  height: 2.75rem;
  border-radius: 50%;
  background-color: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-md);
}
</style>