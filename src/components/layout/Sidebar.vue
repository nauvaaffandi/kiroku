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
      <span class="toggle-badge" :class="{ 'badge-expanded': uiStore.sidebarOpen }">
        <ChevronRight class="toggle-arrow" :class="{ 'arrow-rotated': uiStore.sidebarOpen }" />
      </span>
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
import { ChevronRight } from '@lucide/vue'
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

.toggle-badge {
  @apply flex items-center justify-center w-8 h-8 rounded-md border border-[var(--border-primary)] bg-[var(--bg-hover)] transition-all duration-200;
  flex-shrink: 0;
}

.badge-expanded {
  @apply border-[var(--border-primary)];
}

.toggle-arrow {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
  color: currentColor;
  transition: transform 0.2s ease-out;
}

.arrow-rotated {
  @apply rotate-180;
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
