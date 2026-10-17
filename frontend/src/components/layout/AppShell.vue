<template>
  <div class="app-shell">
    <aside :class="['sidebar', { open: uiStore.sidebarOpen }]" v-show="uiStore.sidebarOpen || !isMobile">
      <Sidebar />
    </aside>
    
    <div class="main-wrapper">
      <header class="app-header" v-show="!isMobile || uiStore.sidebarOpen">
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
    
    <BottomNav v-show="isMobile" />
    
    <FAB v-if="showFAB" @click="onFABClick" />
    
    <div v-if="uiStore.sidebarOpen && isMobile" class="sidebar-overlay" @click="uiStore.toggleSidebar" />
    
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

const isMobile = ref(false)
const showFAB = ref(true)

const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
  if (!isMobile.value) {
    uiStore.sidebarOpen = false
  }
}

const updateFABVisibility = () => {
  const fabRoutes = ['/dashboard', '/transactions', '/accounts', '/categories', '/scheduled']
  showFAB.value = fabRoutes.some(r => route.path.startsWith(r))
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  updateFABVisibility()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', checkMobile)
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
}

.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  width: 280px;
  background-color: var(--color-neutral-100);
  border-right: 1px solid var(--color-neutral-200);
  z-index: 200;
  display: flex;
  flex-direction: column;
  transform: translateX(-100%);
  transition: transform 0.3s ease;
}

.sidebar.open {
  transform: translateX(0);
}

.main-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  margin-left: 0;
  transition: margin-left 0.3s ease;
}

.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background-color: var(--color-neutral-100);
  border-bottom: 1px solid var(--color-neutral-200);
}

.app-main {
  flex: 1;
  padding: 1.5rem;
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
  padding-top: calc(1.5rem + 64px);
}

@media (min-width: 768px) {
  .main-wrapper {
    margin-left: 280px;
  }
  
  .app-header {
    left: 280px;
    right: 0;
  }
  
  .sidebar {
    transform: translateX(0);
  }
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 150;
}
</style>
