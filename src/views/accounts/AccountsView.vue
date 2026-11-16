<template>
  <div class="accounts-page">
    <header class="page-header">
      <div>
        <h1>{{ t('accounts.title') }}</h1>
        <p class="page-subtitle">{{ t('accounts.subtitle') }}</p>
      </div>
      <BaseButton @click="openCreateModal">
        <Plus :size="20" stroke-width="1.75" aria-hidden="true" />
        {{ t('accounts.addAccount') }}
      </BaseButton>
    </header>

    <AccountFilters
      v-model="localFilters"
      :searchQuery="searchQuery"
      @change="handleFilterChange"
      @search="handleSearch"
    />

    <AccountList
      :accounts="accounts"
      @edit="handleEdit"
      @delete="handleDelete"
      @create="openCreateModal"
    />

    <AccountFormModal
      v-model="showFormModal"
      :account="editingAccount"
      @submit="handleFormSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useAccountsStore } from '@/stores/accounts'
import { useI18n } from '@/locales'
import type { Account, AccountFilters as AccountFiltersType } from '@/types'
import { AccountList, AccountFormModal, AccountFilters } from '@/widgets/accounts'
import { BaseButton } from '@/components/ui'
import { Plus } from '@lucide/vue'

const accountsStore = useAccountsStore()
const { t } = useI18n()

const accounts = computed(() => accountsStore.accounts)

const showFormModal = ref(false)
const editingAccount = ref<Account | null>(null)

const localFilters = ref<AccountFiltersType>({
  q: '',
  type: '',
  currency: '',
  archived: '',
  sort: 'name',
  order: 'ASC'
})

const searchQuery = ref('')

let searchTimeout: ReturnType<typeof setTimeout>
watch(searchQuery, (value) => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    accountsStore.filters.q = value
    fetchAccounts(1)
  }, 300)
})

const fetchAccounts = async (page = 1) => {
  console.log('[AccountsView] fetching accounts', { page })
  await accountsStore.fetchAccounts(page)
}

onMounted(() => {
  fetchAccounts(1)
})

const handleFilterChange = () => {
  accountsStore.filters = { ...localFilters.value }
  fetchAccounts(1)
}

const handleSearch = (value: string) => {
  searchQuery.value = value
}

const handleEdit = (account: Account) => {
  editingAccount.value = account
  showFormModal.value = true
}

const handleDelete = (id: string) => {
  const acc = accounts.value.find(a => a.id === id)
  if (confirm(t('accounts.deleteConfirm', { name: acc?.name ?? '' }))) {
    if (import.meta.env.DEV) console.log('[AccountsView] deleting account', { id })
    accountsStore.deleteAccount(id)
  }
}

const openCreateModal = () => {
  editingAccount.value = null
  showFormModal.value = true
}

const handleFormSubmit = async (data: Partial<Account>) => {
  if (import.meta.env.DEV) console.log('[AccountsView] saving account', { id: editingAccount.value?.id, data })
  if (editingAccount.value) {
    await accountsStore.updateAccount(editingAccount.value.id, data)
  } else {
    await accountsStore.createAccount(data)
  }
  showFormModal.value = false
  editingAccount.value = null
  await fetchAccounts(1)
}
</script>

<style scoped>
@reference "tailwindcss";
.page-header {
  @apply flex justify-between items-start mb-6 gap-4;
}

.page-header svg,
.page-header [class*="lucide-icon"],
.page-header i {
  @apply mr-2;
}

.page-header h1 {
  @apply text-[20px] font-semibold m-0;
  color: var(--text-primary);
}

.page-subtitle {
  @apply text-sm mt-1 mb-0;
  color: var(--text-muted);
}

.header-actions {
  @apply flex gap-3 shrink-0;
}
</style>
