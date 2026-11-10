<template>
  <nav class="sidebar-nav">
    <ul class="nav-list">
      <li v-for="item in mainNavItems" :key="item.path">
        <router-link
          :to="item.path"
          :class="['nav-link', { active: isActive(item.path) }, 'text-zinc-700 dark:text-white font-medium font-semibold']"
          @click="closeMobileSidebar"
        >
          <component :is="item.icon" class="nav-icon" :style="{ color: item.color }" />
          <span>{{ item.label }}</span>
        </router-link>
      </li>

      <li class="nav-group">
        <button
          @click="toggleAnalytics"
          :class="['nav-link nav-group-toggle w-full text-zinc-700 dark:text-white font-medium font-semibold', { active: isAnalyticsActive && !isAnalyticsOpen }]"
        >
          <AnalyticsIcon class="nav-icon" style="color:#6366f1" />
          <span class="flex-1 text-left">{{ t('nav.analytics') }}</span>
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            :class="['transition-transform duration-200', { 'rotate-180': isAnalyticsOpen }]"
          >
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>

        <div class="expandable-wrapper" :style="{ height: isAnalyticsOpen ? analyticsHeight + 'px' : '0px' }">
          <ul class="sub-nav-list" ref="subNavRef">
            <li v-for="sub in analyticsItems" :key="sub.path">
              <router-link
                :to="sub.path"
                :class="['sub-nav-link', { active: route.path === sub.path }, 'text-zinc-600 dark:text-zinc-300 font-medium']"
                @click="closeMobileSidebar"
              >
                <component :is="(sub as any).icon" class="sub-nav-icon" v-if="(sub as any).icon" :style="{ color: sub.color }" />
                <span>{{ sub.label }}</span>
              </router-link>
            </li>
          </ul>
        </div>
      </li>

      <li v-for="item in goalsItems" :key="item.path">
        <router-link
          :to="item.path"
          :class="['nav-link', { active: isActive(item.path) }, 'text-zinc-700 dark:text-white font-medium font-semibold']"
          @click="closeMobileSidebar"
        >
          <component :is="item.icon" class="nav-icon" :style="{ color: item.color }" />
          <span>{{ item.label }}</span>
        </router-link>
      </li>

      <li v-for="item in bottomMainNavItems" :key="item.path">
        <router-link
          :to="item.path"
          :class="['nav-link', { active: isActive(item.path) }, 'text-zinc-700 dark:text-white font-medium font-semibold']"
          @click="closeMobileSidebar"
        >
          <component :is="item.icon" class="nav-icon" :style="{ color: item.color }" />
          <span>{{ item.label }}</span>
        </router-link>
      </li>
    </ul>

    <div class="nav-divider" />

    <ul class="nav-list">
      <li v-for="item in settingsItems" :key="item.path">
        <router-link
          :to="item.path"
          :class="['nav-link', { active: isActive(item.path) }, 'text-zinc-700 dark:text-white font-medium font-semibold']"
          @click="closeMobileSidebar"
        >
          <span>{{ item.label }}</span>
        </router-link>
      </li>
    </ul>
  </nav>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import AnalyticsIcon from '@/components/icons/AnalyticsIcon.vue'
import { useAnalyticsNavigation } from '@/composables/useAnalyticsNavigation'
import { useSidebarNav } from '@/composables/useSidebarNav'

const subNavRef = ref<HTMLElement | null>(null)
const route = useRoute()
const {
  isAnalyticsOpen, analyticsHeight, isAnalyticsActive, toggleAnalytics
} = useAnalyticsNavigation(subNavRef)
const {
  mainNavItems, analyticsItems, goalsItems, bottomMainNavItems, settingsItems,
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
  color: var(--text-secondary);
  font: inherit;
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
  @apply mx-[0.125rem] my-[0.125rem] w-[calc(100%-1rem)];
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
</style>
