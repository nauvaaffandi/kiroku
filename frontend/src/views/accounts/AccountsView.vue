<template>
  <div class="accounts-page">
    <header class="page-header">
      <div>
        <h1>{{ t('accounts.title') }}</h1>
        <p class="page-subtitle">{{ t('accounts.subtitle') }}</p>
      </div>
      <BaseButton @click="openCreateModal">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
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
  if (confirm('Are you sure you want to delete this account?')) {
    accountsStore.deleteAccount(id)
  }
}

const openCreateModal = () => {
  editingAccount.value = null
  showFormModal.value = true
}

const handleFormSubmit = async (data: Partial<Account>) => {
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

.page-header h1 {
  @apply text-2xl font-bold m-0;
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