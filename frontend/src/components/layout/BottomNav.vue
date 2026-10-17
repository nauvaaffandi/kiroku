<template>
  <nav class="bottom-nav">
    <router-link
      v-for="item in navItems"
      :key="item.path"
      :to="item.path"
      :class="['bottom-nav-item', { active: isActive(item.path) }]"
    >
      <component :is="item.icon" class="bottom-nav-icon" />
      <span class="bottom-nav-label">{{ item.label }}</span>
    </router-link>
  </nav>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import DashboardIcon from '@/components/icons/DashboardIcon.vue'
import TransactionsIcon from '@/components/icons/TransactionsIcon.vue'
import AccountsIcon from '@/components/icons/AccountsIcon.vue'
import ReportsIcon from '@/components/icons/ReportsIcon.vue'

const route = useRoute()

const navItems = [
  { path: '/dashboard', label: 'Home', icon: DashboardIcon },
  { path: '/transactions', label: 'Txns', icon: TransactionsIcon },
  { path: '/accounts', label: 'Accts', icon: AccountsIcon },
  { path: '/reports', label: 'Reports', icon: ReportsIcon }
]

const isActive = (path: string) => {
  return route.path.startsWith(path)
}
</script>

<style scoped>
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background-color: var(--color-neutral-100);
  border-top: 1px solid var(--color-neutral-200);
  display: flex;
  align-items: center;
  justify-content: space-around;
  z-index: 100;
  padding-bottom: env(safe-area-inset-bottom);
}

.bottom-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  height: 100%;
  color: var(--color-neutral-500);
  text-decoration: none;
  transition: color 0.2s ease;
}

.bottom-nav-item.active {
  color: var(--color-primary);
}

.bottom-nav-icon {
  width: 1.25rem;
  height: 1.25rem;
  margin-bottom: 0.125rem;
}

.bottom-nav-label {
  font-size: 0.6875rem;
  font-weight: 500;
}
</style>