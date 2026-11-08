<template>
  <section class="notifications-page">
    <header class="page-header">
      <h1>{{ t('notifications.title') }}</h1>
      <div class="header-actions">
        <button class="action-button" :disabled="store.unreadCount === 0" @click="store.markAllAsRead">
          {{ t('notifications.markAllRead') }}
        </button>
        <button class="action-button" @click="store.clearRead">
          {{ t('notifications.clearRead') }}
        </button>
      </div>
    </header>

    <div class="push-row">
      <span>{{ push.isSupported ? (push.isSubscribed ? t('notifications.pushEnabled') : t('notifications.pushDisabled')) : t('notifications.pushUnsupported') }}</span>
      <button v-if="push.isSupported" class="action-button" :disabled="push.loading" @click="togglePush">
        {{ push.isSubscribed ? t('notifications.disablePush') : t('notifications.enablePush') }}
      </button>
    </div>

    <nav class="filters" :aria-label="t('notifications.filters')">
      <button class="filter-button" :class="{ active: !unreadOnly }" @click="setFilter(false)">{{ t('notifications.all') }}</button>
      <button class="filter-button" :class="{ active: unreadOnly }" @click="setFilter(true)">{{ t('notifications.unread') }}</button>
    </nav>

    <div v-if="store.loading" class="empty-state">{{ t('common.loading') }}</div>
    <div v-else-if="store.notifications.length === 0" class="empty-state">{{ t('notifications.empty') }}</div>
    <ul v-else class="notification-list">
      <li v-for="item in store.notifications" :key="item.id" class="notification-row" :class="{ unread: !item.read_at }">
        <button class="notification-main" @click="openNotification(item)">
          <component :is="iconFor(item.type)" class="type-icon" aria-hidden="true" />
          <span class="notification-content">
            <span class="notification-title">{{ item.title }}<span v-if="!item.read_at" class="unread-dot" /></span>
            <span class="notification-body">{{ item.body }}</span>
            <time class="notification-time" :datetime="item.created_at">{{ dayjs(item.created_at).fromNow() }}</time>
          </span>
        </button>
        <button class="delete-button" :aria-label="t('notifications.delete')" @click="store.remove(item.id)">
          <Trash2 :size="18" aria-hidden="true" />
        </button>
      </li>
    </ul>
  </section>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/locales'
import { useNotificationStore } from '@/stores/notifications'
import { usePushStore } from '@/stores/push'
import type { NotificationItem } from '@/types'
import { Bell, AlertCircle, Info, CheckCircle, Trash2 } from '@lucide/vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

const { t } = useI18n()
const router = useRouter()
const store = useNotificationStore()
const push = usePushStore()
const unreadOnly = ref(false)

const setFilter = (value: boolean) => {
  unreadOnly.value = value
  store.loadNotifications(value)
}

const togglePush = async () => {
  if (push.isSubscribed) await push.unsubscribe()
  else await push.subscribe()
}

const iconFor = (type: string) => {
  if (/alert|warning|error/i.test(type)) return AlertCircle
  if (/success|complete|paid/i.test(type)) return CheckCircle
  if (/info/i.test(type)) return Info
  return Bell
}

const openNotification = async (item: NotificationItem) => {
  if (!item.read_at) await store.markAsRead(item.id)
  if (item.action_url) router.push(item.action_url)
}

onMounted(() => {
  store.loadNotifications()
  store.loadUnreadCount()
  push.checkSupport()
  push.initializePush()
})
onBeforeUnmount(() => store.stopPolling())
</script>

<style scoped>
@reference "tailwindcss";

.notifications-page {
  @apply mx-auto flex max-w-md flex-col gap-4 px-4 py-6;
}
.page-header {
  @apply flex flex-wrap items-center justify-between gap-3;
}
.page-header h1 {
  @apply m-0 text-xl font-semibold;
  color: var(--text-primary);
}
.header-actions,
.push-row,
.filters {
  @apply flex flex-wrap items-center gap-2;
}
.push-row {
  @apply justify-between rounded-lg border px-4 py-3 text-sm;
  border-color: var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-secondary);
}
.action-button,
.filter-button {
  @apply min-h-11 rounded-lg border px-3 text-sm font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-50;
  border-color: var(--border-primary);
  background-color: var(--bg-card);
  color: var(--text-primary);
}
.action-button:hover,
.filter-button:hover {
  background-color: var(--bg-hover);
}
.action-button:focus-visible,
.filter-button:focus-visible,
.delete-button:focus-visible {
  outline: 2px solid var(--border-focus);
  outline-offset: 2px;
}
.filters {
  @apply border-b;
  border-color: var(--border-primary);
}
.filter-button {
  @apply rounded-b-none border-0 bg-transparent;
}
.filter-button.active {
  @apply border-b-2;
  border-color: var(--accent);
  color: var(--accent);
}
.notification-list {
  @apply m-0 list-none p-0;
}
.notification-list > :not([hidden]) ~ :not([hidden]) {
  border-top: 1px solid var(--border-primary);
}
.notification-row {
  @apply flex items-center gap-2 py-3;
}
.notification-main {
  @apply flex min-w-0 flex-1 items-start gap-3 border-0 bg-transparent p-0 text-left;
}
.type-icon {
  @apply mt-1 h-5 w-5 shrink-0;
  color: var(--text-secondary);
}
.notification-content {
  @apply flex min-w-0 flex-col gap-1;
}
.notification-title {
  @apply flex items-center gap-2 font-medium;
  color: var(--text-primary);
}
.notification-body {
  @apply break-words text-sm;
  color: var(--text-secondary);
}
.notification-time {
  @apply text-xs;
  color: var(--text-muted);
}
.unread-dot {
  @apply h-2 w-2 shrink-0 rounded-full;
  background-color: var(--accent);
}
.delete-button {
  @apply flex min-h-11 min-w-11 items-center justify-center rounded-lg border-0 bg-transparent;
  color: var(--text-muted);
}
.delete-button:hover {
  background-color: var(--bg-hover);
  color: var(--error);
}
.empty-state {
  @apply py-12 text-center text-sm;
  color: var(--text-muted);
}
</style>
