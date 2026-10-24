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
    
    <FAB v-if="showFAB" @click="onFABClick" />
    
    <div v-if="uiStore.sidebarOpen" class="sidebar-overlay" @click="uiStore.toggleSidebar">
      <aside class="sidebar open" @click.stop>
        <Sidebar @close="uiStore.toggleSidebar" />
      </aside>
    </div>
    
    <ToastContainer />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUIStore } from '@/stores/ui'

import Sidebar from './Sidebar.vue'
import Header from './Header.vue'
import BottomNav from './BottomNav.vue'
import FAB from '@/components/forms/FAB.vue'
import ToastContainer from '@/components/common/ToastContainer.vue'

const router = useRouter()
const route = useRoute()
const uiStore = useUIStore()

const showFAB = ref(true)

const updateFABVisibility = () => {
  const fabRoutes = ['/dashboard', '/transactions', '/accounts', '/categories', '/scheduled']
  showFAB.value = fabRoutes.some(r => route.path.startsWith(r))
}

onMounted(() => {
  updateFABVisibility()
})

watch(() => route.path, updateFABVisibility)

const onFABClick = () => {
  if (route.path === '/transactions') {
    router.push({ name: 'Transactions', query: { create: 'true' } })
  } else {
    router.push('/transactions?create=true')
  }
}
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
  background-color: var(--bg-primary);
  max-width: 480px;
  margin: 0 auto;
  width: 100%;
  box-shadow: var(--shadow-xl);
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 200;
  display: flex;
  justify-content: flex-start;
}

.sidebar {
  width: 280px;
  height: 100vh;
  background-color: var(--bg-card);
  border-right: 1px solid var(--border-primary);
  display: flex;
  flex-direction: column;
  animation: slideRight 0.3s ease;
}

@keyframes slideRight {
  from { transform: translateX(-100%); }
  to { transform: translateX(0); }
}

.main-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  padding-bottom: 70px; /* space for bottom nav */
}

.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background-color: var(--bg-card);
  border-bottom: 1px solid var(--border-primary);
}

.app-main {
  flex: 1;
  padding: 1rem;
  width: 100%;
}
</style>
