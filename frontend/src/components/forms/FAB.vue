<template>
  <div class="fab-wrapper">
    <div v-if="isOpen" class="fab-backdrop" @click="toggleMenu"></div>

    <div class="fab-children-container" :class="{ open: isOpen }">
      <button
        v-for="(child, index) in children"
        :key="child.id"
        class="fab-child-btn"
        :style="getChildStyle(index, children.length)"
        @click.stop="handleChildClick(child)"
        :title="child.label"
      >
        <component :is="child.icon" class="fab-child-icon" />
      </button>
    </div>

    <button
      class="fab-main"
      :class="{ open: isOpen }"
      @click="toggleMenu"
      aria-label="Toggle menu"
    >
      <div class="fab-icon-wrapper">
        <svg v-if="!isOpen" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </div>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import TransactionsIcon from '@/components/icons/TransactionsIcon.vue'
import CalendarIcon from '@/components/icons/CalendarIcon.vue'
import AccountsIcon from '@/components/icons/AccountsIcon.vue'
import ReportsIcon from '@/components/icons/ReportsIcon.vue'
import GoalsIcon from '@/components/icons/GoalsIcon.vue'

const router = useRouter()
const isOpen = ref(false)

// Exactly 5 child buttons, fully dynamic array ready for 6+ items in the future
const children = computed(() => [
  { id: 'add', label: 'Add Transaction', icon: TransactionsIcon, action: () => router.push('/transactions?create=true') },
  { id: 'scan', label: 'Scan Receipt', icon: CalendarIcon, action: () => router.push('/transactions/scan') },
  { id: 'accounts', label: 'Accounts', icon: AccountsIcon, action: () => router.push('/accounts') },
  { id: 'reports', label: 'Reports', icon: ReportsIcon, action: () => router.push('/reports') },
  { id: 'goals', label: 'Goals', icon: GoalsIcon, action: () => router.push('/analytics/goals') }
])

const toggleMenu = () => {
  if (navigator.vibrate) {
    navigator.vibrate(10)
  }
  isOpen.value = !isOpen.value
}

const handleChildClick = (child: any) => {
  if (navigator.vibrate) {
    navigator.vibrate(15)
  }
  isOpen.value = false
  child.action()
}

// Compute radial coordinates for "water splash" explosion effect in an upper semi-circle arc
const getChildStyle = (index: number, total: number) => {
  if (!isOpen.value) {
    return {
      transform: 'translate(0px, 0px) scale(0)',
      opacity: '0',
      pointerEvents: 'none' as const,
      transitionDelay: '0s'
    }
  }

  const radius = 110 // distance from center
  const startAngle = -Math.PI * 0.92 // ~ -165 deg (left upward)
  const endAngle = -Math.PI * 0.08   // ~ -15 deg (right upward)
  const angleStep = total > 1 ? (endAngle - startAngle) / (total - 1) : 0
  const angle = startAngle + index * angleStep

  const x = Math.cos(angle) * radius
  const y = Math.sin(angle) * radius

  // Staggered delay for cascading splash effect
  const delay = index * 0.04

  return {
    transform: `translate(${x}px, ${y}px) scale(1)`,
    opacity: '1',
    pointerEvents: 'auto' as const,
    transitionDelay: `${delay}s`
  }
}
</script>

<style scoped>
@reference "tailwindcss";
.fab-wrapper {
  @apply fixed bottom-[2rem] right-[2rem] z-[999];
}

.fab-backdrop {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(2px);
  z-index: 998;
  animation: fadeIn 0.25s ease forwards;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.fab-main {
  @apply relative w-[3.75rem] h-[3.75rem] rounded-full text-white border-0 cursor-pointer flex items-center justify-center z-[1000] transition-transform transition-colors;
  background-color: var(--accent);
  box-shadow: 0 8px 24px rgba(74, 144, 217, 0.4);
}

.fab-main:hover {
  @apply scale-[1.08];
  background-color: var(--accent-hover);
}

.fab-main:active {
  @apply scale-[0.92];
}

.fab-main.open {
  background-color: var(--error);
  box-shadow: 0 8px 24px rgba(220, 53, 69, 0.4);
  transform: rotate(90deg);
}

.fab-icon-wrapper {
  @apply flex items-center justify-center;
  transition: transform 0.3s ease;
}

.fab-children-container {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 0;
  height: 0;
  @apply flex items-center justify-center z-[999];
}

.fab-child-btn {
  position: absolute;
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  background-color: var(--bg-card);
  color: var(--accent);
  border: 1px solid var(--border-primary);
  @apply flex items-center justify-center cursor-pointer shadow-md transition-all duration-300 ease-out;
  box-shadow: var(--shadow-md);
  transform: translate(0, 0) scale(0);
  opacity: 0;
}

.fab-child-btn:hover {
  background-color: var(--accent);
  color: white;
  transform: scale(1.15) !important;
}

.fab-child-icon {
  width: 1.25rem;
  height: 1.25rem;
}
</style>