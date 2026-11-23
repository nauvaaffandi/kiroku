<template>
  <div class="app-shell">
    <div class="app-body">
      <aside
        ref="railEl"
        class="sidebar-rail"
        :class="{ 'sidebar-overlay': uiStore.sidebarOpen }"
        :style="{ width: sidebarWidth + 'px' }"
      >
        <Sidebar />
      </aside>

      <div v-if="uiStore.sidebarOpen" class="sidebar-backdrop" @click="uiStore.toggleSidebar"></div>

      <div class="main-wrapper" :style="{ marginLeft: uiStore.sidebarOpen ? sidebarWidth + 'px' : '0px' }">
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
    </div>

    <ToastContainer />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useUIStore } from '@/stores/ui'
import { useSidebarRailGesture } from '@/composables/useSidebarRailGesture'

import Sidebar from './Sidebar.vue'
import Header from './Header.vue'
import ToastContainer from '@/components/common/ToastContainer.vue'

const uiStore = useUIStore()
const railEl = ref<HTMLElement | null>(null)
const sidebarWidth = computed(() => (uiStore.sidebarOpen ? 280 : 64))
useSidebarRailGesture(railEl)
</script>

<style scoped>
@reference "tailwindcss";

.app-shell {
  @apply flex flex-col m-auto w-full max-w-[480px];
  min-height: 100vh;
  position: relative;
  background-color: var(--bg-primary);
  overflow-x: hidden;
}

.app-body {
  @apply flex flex-1 min-h-[100vh] min-w-0;
  position: relative;
  min-h: 0;
}

.sidebar-rail {
  @apply flex-shrink-0 border-r min-w-0;
  position: sticky;
  top: 0;
  height: 100svh;
  overflow-x: hidden;
  overflow-y: auto;
  touch-action: pan-y;
  transition: width 0.3s ease;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}

.sidebar-rail.sidebar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 40;
  height: 100%;
}

.sidebar-backdrop {
  position: absolute;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.4);
  z-index: 35;
}

.main-wrapper {
  @apply flex-1 flex flex-col min-h-0 min-w-0;
  transition: margin-left 0.3s ease;
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
  padding-bottom: calc(env(safe-area-inset-bottom) + 1rem);
}
</style>
