<template>
  <div class="app-shell">
    <div class="main-wrapper">
      <header class="app-header">
        <Header @toggle-sidebar="uiStore.toggleSidebar" />
      </header>

      <main class="app-main">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <BottomNav />

    <div v-if="uiStore.sidebarOpen" class="sidebar-overlay" @click="uiStore.toggleSidebar">
      <aside class="sidebar open" @click.stop>
        <Sidebar @close="uiStore.toggleSidebar" />
      </aside>
    </div>

    <ToastContainer />
  </div>
</template>

<script setup lang="ts">
import { useUIStore } from '@/stores/ui'

import Sidebar from './Sidebar.vue'
import Header from './Header.vue'
import BottomNav from './BottomNav.vue'
import ToastContainer from '@/components/common/ToastContainer.vue'

const uiStore = useUIStore()
</script>

<style scoped>
@reference "tailwindcss";

.app-shell {
  @apply flex flex-col m-auto w-full max-w-[480px];
  min-height: 100vh;
  position: relative;
  background-color: var(--bg-primary);
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.4);
  @apply z-[200] flex justify-start;
}

.sidebar {
  @apply w-[280px] h-full flex flex-col border-r;
  height: 100vh;
  background-color: var(--bg-card);
  border-right-color: var(--border-primary);
  animation: slideRight 0.3s ease;
}

@keyframes slideRight {
  from { transform: translateX(-100%); }
  to { transform: translateX(0); }
}

.main-wrapper {
  @apply flex-1 flex flex-col min-h-[100vh];
}

.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background-color: var(--bg-card);
  border-bottom: 1px solid var(--border-primary);
}

.app-main {
  @apply flex-1 p-4 w-full;
  /* Keep content clear of the fixed bottom nav + safe area */
  padding-bottom: calc(60px + env(safe-area-inset-bottom) + 1rem);
}
</style>
