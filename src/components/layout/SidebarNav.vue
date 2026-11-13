<template>
  <nav class="sidebar-nav">
    <ul class="nav-list">
      <li v-for="item in mainNavItems" :key="item.path">
        <router-link
          :to="item.path"
          :class="['nav-link', { active: isActive(item.path) }]"
          @click="closeMobileSidebar"
        >
          <component :is="item.icon" class="nav-icon" :style="{ color: item.color }" />
          <span>{{ item.label }}</span>
        </router-link>
      </li>

      <li class="nav-group">
        <button
          @click="toggleAnalytics"
          :class="['nav-link nav-group-toggle', { active: isAnalyticsActive && !isAnalyticsOpen }]"
        >
          <AnalyticsIcon class="nav-icon nav-icon-analytics" />
          <span class="nav-label">{{ t('nav.analytics') }}</span>
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            :class="['nav-expand-icon', { 'rotate-180': isAnalyticsOpen }]"
          >
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>

        <div class="expandable-wrapper" :style="{ height: isAnalyticsOpen ? analyticsHeight + 'px' : '0px' }">
          <ul class="sub-nav-list" ref="subNavRef">
            <li v-for="sub in analyticsItems" :key="sub.path">
              <router-link
                :to="sub.path"
                :class="['sub-nav-link', { active: route.path === sub.path }]"
                @click="closeMobileSidebar"
              >
                <component :is="(sub as any).icon" class="sub-nav-icon" v-if="(sub as any).icon" :style="{ color: sub.color }" />
                <span>{{ sub.label }}</span>
              </router-link>
            </li>
          </ul>
        </div>
      </li>
      <li class="nav-divider" aria-hidden="true" />

      <li v-for="item in goalsItems" :key="item.path">
        <router-link
          :to="item.path"
          :class="['nav-link', { active: isActive(item.path) }]"
          @click="closeMobileSidebar"
        >
          <component :is="item.icon" class="nav-icon" :style="{ color: item.color }" />
          <span>{{ item.label }}</span>
        </router-link>
      </li>
      <li class="nav-divider" aria-hidden="true" />

      <li v-for="item in section3NavItems" :key="item.path">
        <router-link
          :to="item.path"
          :class="['nav-link', { active: isActive(item.path) }]"
          @click="closeMobileSidebar"
        >
          <component :is="item.icon" class="nav-icon" :style="{ color: item.color }" />
          <span>{{ item.label }}</span>
        </router-link>
      </li>
      <li class="nav-divider" aria-hidden="true" />

      <li v-for="item in section4NavItems" :key="item.path">
        <router-link
          :to="item.path"
          :class="['nav-link', { active: isActive(item.path) }]"
          @click="closeMobileSidebar"
        >
          <component :is="item.icon" class="nav-icon" :style="{ color: item.color }" />
          <span>{{ item.label }}</span>
        </router-link>
      </li>
      <li class="nav-divider" aria-hidden="true" />

      <li class="nav-group">
        <router-link
          to="/settings"
          :class="['nav-link', { active: isActive('/settings') }]"
          @click="closeMobileSidebar"
        >
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="nav-icon nav-icon-settings">
            <circle cx="12" cy="12" r="3"></circle>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06-.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
          </svg>
          <span>{{ t('nav.settings') }}</span>
        </router-link>
        <ThemeToggle />
      </li>
    </ul>

    <div class="nav-divider" />

  </nav>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import AnalyticsIcon from '@/components/icons/AnalyticsIcon.vue'
import { useAnalyticsNavigation } from '@/composables/useAnalyticsNavigation'
import { useSidebarNav } from '@/composables/useSidebarNav'
import ThemeToggle from '@/components/common/ThemeToggle.vue'

const subNavRef = ref<HTMLElement | null>(null)
const route = useRoute()
const {
  isAnalyticsOpen, analyticsHeight, isAnalyticsActive, toggleAnalytics
} = useAnalyticsNavigation(subNavRef)
const {
  mainNavItems, analyticsItems, goalsItems, section3NavItems, section4NavItems,
  isActive, closeMobileSidebar, t
} = useSidebarNav()
</script>

<style scoped>
@reference "tailwindcss";

.sidebar-nav {
  @apply flex-1 overflow-y-auto px-2 py-2;
}

.nav-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-link {
  @apply flex items-center gap-3 px-4 py-3 rounded-sm no-underline mx-[0.125rem] my-[0.125rem] w-auto bg-transparent border-0 cursor-pointer text-left transition-colors;
  @apply text-zinc-700 dark:text-white font-medium font-semibold;
}

.nav-link:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.nav-link.active {
  background-color: var(--accent-light);
  color: var(--accent);
  font-weight: 500;
}

.nav-group-toggle {
  @apply mx-[0.125rem] my-[0.125rem] w-full text-zinc-700 dark:text-white font-medium font-semibold;
}

.expandable-wrapper {
  overflow: hidden;
  transition: height 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  @apply mx-[0.5rem];
}

.sub-nav-list {
  list-style: none;
  padding: 0.25rem 0 0.25rem 2rem;
  margin: 0;
  @apply flex flex-col gap-0.5;
}

.sub-nav-link {
  @apply flex items-center gap-2 px-3 py-2 rounded-sm text-sm transition-colors no-underline;
  @apply text-zinc-600 dark:text-zinc-300 font-medium;
  color: var(--text-secondary);
}

.sub-nav-link:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.sub-nav-link.active {
  background-color: var(--accent-light);
  color: var(--accent);
  font-weight: 600;
}

.nav-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
  color: currentColor;
}

.nav-icon-analytics {
  color: var(--accent);
}

.nav-icon-settings {
  color: var(--text-secondary);
}

.sub-nav-icon {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
  color: currentColor;
}

.nav-divider {
  height: 1px;
  background-color: var(--border-primary);
  @apply my-[1rem] mx-[0.5rem];
}

.nav-label {
  @apply flex-1 text-left;
}

.nav-expand-icon {
  @apply transition-transform duration-200;
}
</style>
