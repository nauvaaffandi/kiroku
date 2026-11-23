<template>
  <div class="sidebar">
    <button
      class="sidebar-toggle"
      :class="{ 'toggle-expanded': uiStore.sidebarOpen }"
      @click="uiStore.toggleSidebar"
      :title="uiStore.sidebarOpen ? t('sidebar.collapse') : t('sidebar.expand')"
      :aria-label="uiStore.sidebarOpen ? t('sidebar.collapse') : t('sidebar.expand')"
      :aria-expanded="uiStore.sidebarOpen"
    >
      <ChevronLeft v-if="uiStore.sidebarOpen" class="toggle-icon" />
      <ChevronRight v-else class="toggle-icon" />
    </button>

    <div v-if="uiStore.sidebarOpen" class="sidebar-header">
      <router-link to="/" class="logo">
        <AppIcon class="app-icon-logo" />
        <span>{{ t('common.appName') }}</span>
      </router-link>
    </div>

    <SidebarNav :collapsed="!uiStore.sidebarOpen" />
    <SidebarFooter />
  </div>
</template>

<script setup lang="ts">
import { useUIStore } from '@/stores/ui'
import { ChevronLeft, ChevronRight } from '@lucide/vue'
import AppIcon from '@/components/icons/AppIcon.vue'
import SidebarNav from './SidebarNav.vue'
import SidebarFooter from './SidebarFooter.vue'
import { useI18n } from '@/locales'

const uiStore = useUIStore()
const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";

.sidebar {
  @apply flex flex-col h-full;
}

.sidebar-toggle {
  @apply flex items-center w-full px-4 py-3 rounded-sm border-0 bg-transparent cursor-pointer transition-colors text-left;
  @apply justify-center;
  @apply text-[var(--text-secondary)];
}

.toggle-expanded {
  @apply justify-start;
}

.sidebar-toggle:hover {
  @apply hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)];
}

.sidebar-toggle:focus-visible {
  @apply outline-none ring-2 ring-[var(--accent)];
}

.toggle-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
  color: currentColor;
}

.sidebar-header {
  @apply flex items-center gap-3 px-6 py-4 border-b;
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
</style>
