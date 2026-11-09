import { computed, onMounted, ref } from 'vue'
import { api } from '@/api/client'
import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'

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
  const shares = ref<Share[]>([])
  const accounts = ref<Account[]>([])
  const loading = ref(false)
  const showCreateModal = ref(false)
  const createdShare = ref<Share | null>(null)
  const createForm = ref<CreateShareForm>({ account_id: '', expires_at: '' })

  const loadData = async () => {
    try {
      const [sharesRes, accountsRes] = await Promise.all([
        api.get('/api/v1/shared'),
        api.get('/api/v1/accounts')
      ])
      shares.value = sharesRes.data || []
      accounts.value = accountsRes.data.filter((a: Account) => !a.archived) || []
    } catch (e) {
      uiStore.error(t('shared.access.failLoad'))
    }
  }

  const createShare = async () => {
    loading.value = true
    try {
      const payload: { account_id: string; expires_at?: string } = {
        account_id: createForm.value.account_id
      }
      if (createForm.value.expires_at) {
        payload.expires_at = new Date(createForm.value.expires_at).toISOString()
      }
      const res = await api.post('/api/v1/shared', payload)
      createdShare.value = res.data
      showCreateModal.value = false
      await loadData()
    } catch (e: any) {
      uiStore.error(e.response?.data?.error || t('shared.access.failCreate'))
    } finally {
      loading.value = false
    }
  }

  const revokeShare = async (id: string) => {
    if (!confirm(t('shared.access.revokeConfirm'))) return
    try {
      await api.delete(`/api/v1/shared/${id}`)
      uiStore.success(t('shared.access.revokedSuccess'))
      await loadData()
    } catch (e) {
      uiStore.error(t('shared.access.failRevoke'))
    }
  }

  const copyShareUrl = async (token: string) => {
    const url = `${window.location.origin}/shared/public/${token}`
    await navigator.clipboard.writeText(url)
    uiStore.success(t('shared.access.copied'))
  }

  const fullShareUrl = computed(() => {
    if (!createdShare.value) return ''
    return `${window.location.origin}/shared/public/${createdShare.value.token}`
  })

  const copyFullShareUrl = async () => {
    await navigator.clipboard.writeText(fullShareUrl.value)
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

  return {
    accounts, closeCreateModal, closeCreatedModal, copyFullShareUrl, copyShareUrl,
    createForm, createShare, createdShare, formatDate, fullShareUrl, loading,
    revokeShare, shares, showCreateModal, t
  }
}
