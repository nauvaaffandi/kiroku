import { computed, onMounted, ref } from 'vue'
import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { useAccountsStore } from '@/stores/accounts'

export interface Share {
  id: string
  account_id: string
  account_name: string
  token: string
  expires_at?: string | null
  created_at?: string
  last_accessed_at?: string | null
}

export interface Account {
  id: string
  name: string
  currency: string
  archived?: boolean
}

export interface CreateShareForm {
  account_id: string
  expires_at: string
}

export function useShare() {
  const { t, locale: i18nLocale } = useI18n()
  const uiStore = useUIStore()
  const accountsStore = useAccountsStore()
  const shares = ref<Share[]>([])
  const loading = ref(false)
  const showCreateModal = ref(false)
  const createdShare = ref<Share | null>(null)
  const createForm = ref<CreateShareForm>({ account_id: '', expires_at: '' })

  const loadData = async () => {
    try {
      await accountsStore.fetchAccounts()
      shares.value = []
    } catch (e) {
      uiStore.error(t('shared.access.failLoad'))
    }
  }

  const createShare = async () => {
    uiStore.error('Share creation not available offline')
  }

  const revokeShare = async (_id: string) => {
    uiStore.error('Share revoke not available offline')
  }

  const copyShareUrl = async (_token: string) => {
    await navigator.clipboard.writeText('')
    uiStore.success(t('shared.access.copied'))
  }

  const fullShareUrl = computed(() => '')

  const copyFullShareUrl = async () => {
    await navigator.clipboard.writeText('')
    uiStore.success(t('shared.access.copied'))
  }

  const closeCreateModal = () => {
    showCreateModal.value = false
    createForm.value = { account_id: '', expires_at: '' }
  }

  const closeCreatedModal = () => {
    createdShare.value = null
  }

  const formatDate = (dateStr?: string | null) => {
    if (!dateStr) return ''
    return new Date(dateStr).toLocaleDateString(i18nLocale, { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
  }

  onMounted(loadData)

  const accounts = computed(() => accountsStore.accounts.filter((a: Account) => !a.archived))

  return {
    accounts, closeCreateModal, closeCreatedModal, copyFullShareUrl, copyShareUrl,
    createForm, createShare, createdShare, formatDate, fullShareUrl, loading,
    revokeShare, shares, showCreateModal, t
  }
}