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
    
    <FAB v-if="authStore.isAuthenticated" />
    
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
import { useAuthStore } from '@/stores/auth'

import Sidebar from './Sidebar.vue'
import Header from './Header.vue'
import FAB from '@/components/forms/FAB.vue'
import ToastContainer from '@/components/common/ToastContainer.vue'

const uiStore = useUIStore()
const authStore = useAuthStore()
</script>

<style scoped>
@reference "tailwindcss";
.app-shell {
  @apply flex flex-col m-auto w-full max-w-[480px];
  min-height: 100vh;
  position: relative;
  background-color: var(--bg-primary);
  box-shadow: var(--shadow-xl);
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
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

:where([data-theme="dark"], [data-theme="dark"] *) .app-shell {
  background-color: #121212;
}

:where([data-theme="dark"], [data-theme="dark"] *) .sidebar {
  background-color: #18181b;
  border-right-color: #27272a;
}

:where([data-theme="dark"], [data-theme="dark"] *) .app-header {
  background-color: #18181b;
  border-bottom-color: #27272a;
}

.app-main {
  @apply flex-1 p-4 w-full;
}
</style>
